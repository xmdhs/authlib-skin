package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service/auth"
	"github.com/xmdhs/authlib-skin/service/captcha"
	"github.com/xmdhs/authlib-skin/service/email"
	"github.com/xmdhs/authlib-skin/utils"
)

var (
	ErrExistUser  = errors.New("邮箱已存在")
	ErrExitsName  = errors.New("用户名已存在")
	ErrRegLimit   = errors.New("超过注册 ip 限制")
	ErrPassWord   = errors.New("错误的密码或用户名")
	ErrChangeName = errors.New("离线模式 uuid 不允许修改用户名")
	ErrUsername   = errors.New("邮箱不存在")
)

type UserService struct {
	config         config.Config
	client         *ent.Client
	captchaService *captcha.CaptchaService
	authService    *auth.AuthService
	cache          cache.Cache
	emailService   *email.EmailService
}

func NewUserSerice(config config.Config, client *ent.Client, captchaService *captcha.CaptchaService,
	authService *auth.AuthService, cache cache.Cache, emailService *email.EmailService) *UserService {
	return &UserService{
		config:         config,
		client:         client,
		captchaService: captchaService,
		authService:    authService,
		cache:          cache,
		emailService:   emailService,
	}
}

func (w *UserService) Reg(ctx context.Context, u model.UserReg, ipPrefix, ip string) (model.LoginRep, error) {
	var userUuid string
	if w.config.OfflineUUID {
		userUuid = utils.UUIDGen(u.Name)
	} else {
		userUuid = strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	if w.config.Email.Enable {
		err := w.emailService.VerifyJwt(u.Email, u.EmailJwt, "/register")
		if err != nil {
			return model.LoginRep{}, fmt.Errorf("Reg: %w", err)
		}
	}

	err := w.captchaService.VerifyCaptcha(ctx, u.CaptchaToken, ip)
	if err != nil {
		return model.LoginRep{}, fmt.Errorf("Reg: %w", err)
	}

	if w.config.MaxIpUser != 0 {
		c, err := w.client.User.Query().Where(user.RegIPEQ(ipPrefix)).Count(ctx)
		if err != nil {
			return model.LoginRep{}, fmt.Errorf("Reg: %w", err)
		}
		if c >= w.config.MaxIpUser {
			return model.LoginRep{}, fmt.Errorf("Reg: %w", ErrRegLimit)
		}
	}

	p, s := utils.Argon2ID(u.Password)

	var du *ent.User

	err = utils.WithTx(ctx, w.client, func(tx *ent.Tx) error {
		count, err := tx.User.Query().Where(user.EmailEQ(u.Email)).ForUpdateA().Count(ctx)
		if err != nil {
			return err
		}
		if count != 0 {
			return ErrExistUser
		}
		nameCount, err := tx.UserProfile.Query().Where(userprofile.NameEQ(u.Name)).ForUpdateA().Count(ctx)
		if err != nil {
			return err
		}
		if nameCount != 0 {
			return ErrExitsName
		}
		du, err = tx.User.Create().
			SetEmail(u.Email).
			SetPassword(p).
			SetSalt(s).
			SetRegTime(time.Now().Unix()).
			SetRegIP(ipPrefix).
			SetState(0).Save(ctx)
		if err != nil {
			return err
		}
		_, err = tx.UserProfile.Create().
			SetUser(du).
			SetName(u.Name).
			SetUUID(userUuid).
			Save(ctx)
		if err != nil {
			return err
		}
		if du.ID == 1 {
			err := tx.User.UpdateOne(du).SetState(auth.SetAdmin(0, true)).Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return model.LoginRep{}, fmt.Errorf("Reg: %w", err)
	}
	jwt, err := w.authService.CreateToken(ctx, du, "web", userUuid)
	if err != nil {
		return model.LoginRep{}, fmt.Errorf("Login: %w", err)
	}

	return model.LoginRep{
		Token: jwt,
		Name:  u.Name,
		UUID:  userUuid,
	}, nil
}

func (w *UserService) Login(ctx context.Context, l model.Login, ip string) (model.LoginRep, error) {
	err := w.captchaService.VerifyCaptcha(ctx, l.CaptchaToken, ip)
	if err != nil {
		return model.LoginRep{}, fmt.Errorf("Login: %w", err)
	}
	u, err := w.client.User.Query().Where(user.Email(l.Email)).WithProfile().Only(ctx)
	if err != nil {
		var ne *ent.NotFoundError
		if errors.As(err, &ne) {
			return model.LoginRep{}, fmt.Errorf("Login: %w", ErrPassWord)
		}
		return model.LoginRep{}, fmt.Errorf("Login: %w", err)
	}
	err = validatePass(ctx, u, l.Password)
	if err != nil {
		return model.LoginRep{}, fmt.Errorf("Login: %w", err)
	}
	jwt, err := w.authService.CreateToken(ctx, u, "web", u.Edges.Profile.UUID)
	if err != nil {
		return model.LoginRep{}, fmt.Errorf("Login: %w", err)
	}
	return model.LoginRep{
		Token: jwt,
		Name:  u.Edges.Profile.Name,
		UUID:  u.Edges.Profile.UUID,
	}, nil
}

func (w *UserService) Info(ctx context.Context, t *model.TokenClaims) (model.UserInfo, error) {
	u, err := w.client.User.Query().Where(user.ID(t.UID)).First(ctx)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("Info: %w", err)
	}
	isAdmin := auth.IsAdmin(u.State)
	return model.UserInfo{
		UID:     t.UID,
		UUID:    t.Subject,
		IsAdmin: isAdmin,
	}, nil
}

func (w *UserService) ChangePasswd(ctx context.Context, p model.ChangePasswd, uid int, validOldPass bool) error {
	u, err := w.client.User.Query().Where(user.IDEQ(uid)).WithToken().First(ctx)
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	if validOldPass {
		err := validatePass(ctx, u, p.Old)
		if err != nil {
			return fmt.Errorf("ChangePasswd: %w", err)
		}
	}
	pass, salt := utils.Argon2ID(p.New)
	if u.Edges.Token != nil {
		err := w.client.UserToken.UpdateOne(u.Edges.Token).AddTokenID(1).Exec(ctx)
		if err != nil {
			return fmt.Errorf("ChangePasswd: %w", err)
		}
	}
	err = w.cache.Del([]byte("auth" + strconv.Itoa(uid)))
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	err = w.client.User.UpdateOne(u).SetPassword(pass).SetSalt(salt).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	return nil
}

func (w *UserService) changeName(ctx context.Context, newName string, uid int, uuid string) error {
	if w.config.OfflineUUID {
		return fmt.Errorf("changeName: %w", ErrChangeName)
	}
	c, err := w.client.UserProfile.Query().Where(userprofile.Name(newName)).Count(ctx)
	if err != nil {
		return fmt.Errorf("changeName: %w", err)
	}
	if c != 0 {
		return fmt.Errorf("changeName: %w", ErrExitsName)
	}
	err = w.client.UserProfile.Update().Where(userprofile.HasUserWith(user.ID(uid))).SetName(newName).Exec(ctx)
	if err != nil {
		return fmt.Errorf("changeName: %w", err)
	}
	w.cache.Del([]byte("Profile" + uuid))
	return err
}

func (w *UserService) ChangeName(ctx context.Context, newName string, t *model.TokenClaims) error {
	err := w.changeName(ctx, newName, t.UID, t.Subject)
	if err != nil {
		return fmt.Errorf("ChangeName: %w", err)
	}
	return nil
}

func (w *UserService) SendRegEmail(ctx context.Context, email, CaptchaToken, host, ip string) error {
	err := w.captchaService.VerifyCaptcha(ctx, CaptchaToken, ip)
	if err != nil {
		return fmt.Errorf("SendRegEmail: %w", err)
	}

	err = w.emailService.SendVerifyUrl(ctx, email, 60, host, "验证你的邮箱以完成注册", "点击下方链接完成注册，1 天内有效", "/register")
	if err != nil {
		return fmt.Errorf("SendRegEmail: %w", err)
	}
	return nil
}

func (w *UserService) SendChangePasswordEmail(ctx context.Context, email, CaptchaToken, host, ip string) error {
	err := w.captchaService.VerifyCaptcha(ctx, CaptchaToken, ip)
	if err != nil {
		return fmt.Errorf("SendChangePasswordEmail: %w", err)
	}
	c, err := w.client.User.Query().Where(user.Email(email)).Count(ctx)
	if err != nil {
		return fmt.Errorf("SendChangePasswordEmail: %w", err)
	}
	if c == 0 {
		return fmt.Errorf("SendChangePasswordEmail: %w", ErrUsername)
	}
	err = w.emailService.SendVerifyUrl(ctx, email, 60, host, "重设密码", "点击下方链接更改你的密码，1 天内有效", "/forgot")
	if err != nil {
		return fmt.Errorf("SendChangePasswordEmail: %w", err)
	}
	return nil
}

func (w *UserService) ForgotPassword(ctx context.Context, email, passWord, emailJwt string) error {
	err := w.emailService.VerifyJwt(email, emailJwt, "/forgot")
	if err != nil {
		return fmt.Errorf("ForgotPassword: %w", err)
	}
	u, err := w.client.User.Query().Where(user.Email(email)).First(ctx)
	if err != nil {
		return fmt.Errorf("ForgotPassword: %w", err)
	}

	err = w.ChangePasswd(ctx, model.ChangePasswd{New: passWord}, u.ID, false)
	if err != nil {
		return fmt.Errorf("ForgotPassword: %w", err)
	}
	return nil
}
