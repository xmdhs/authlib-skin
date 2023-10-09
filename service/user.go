package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/model"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

var (
	ErrExistUser  = errors.New("邮箱已存在")
	ErrExitsName  = errors.New("用户名已存在")
	ErrRegLimit   = errors.New("超过注册 ip 限制")
	ErrPassWord   = errors.New("错误的密码或用户名")
	ErrChangeName = errors.New("离线模式 uuid 不允许修改用户名")
)

func (w *WebService) Reg(ctx context.Context, u model.UserReg, ipPrefix, ip string) error {
	var userUuid string
	if w.config.OfflineUUID {
		userUuid = utils.UUIDGen(u.Name)
	} else {
		userUuid = strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	err := w.verifyCaptcha(ctx, u.CaptchaToken, ip)
	if err != nil {
		return fmt.Errorf("Reg: %w", err)
	}

	if w.config.MaxIpUser != 0 {
		c, err := w.client.User.Query().Where(user.RegIPEQ(ipPrefix)).Count(ctx)
		if err != nil {
			return fmt.Errorf("Reg: %w", err)
		}
		if c >= w.config.MaxIpUser {
			return fmt.Errorf("Reg: %w", ErrRegLimit)
		}
	}

	p, s := utils.Argon2ID(u.Password)

	err = utils.WithTx(ctx, w.client, func(tx *ent.Tx) error {
		count, err := tx.User.Query().Where(user.EmailEQ(u.Email)).ForUpdate().Count(ctx)
		if err != nil {
			return err
		}
		if count != 0 {
			return ErrExistUser
		}
		nameCount, err := tx.UserProfile.Query().Where(userprofile.NameEQ(u.Name)).ForUpdate().Count(ctx)
		if err != nil {
			return err
		}
		if nameCount != 0 {
			return ErrExitsName
		}
		du, err := tx.User.Create().
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
		return nil
	})
	if err != nil {
		return fmt.Errorf("Reg: %w", err)
	}
	return nil
}

func (w *WebService) Login(ctx context.Context, l model.Login, ip string) (model.LoginRep, error) {
	err := w.verifyCaptcha(ctx, l.CaptchaToken, ip)
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
	err = w.validatePass(ctx, u, l.Password)
	if err != nil {
		return model.LoginRep{}, fmt.Errorf("Login: %w", err)
	}
	jwt, err := utilsService.CreateToken(ctx, u, w.client, w.cache, w.prikey, "web")
	if err != nil {
		return model.LoginRep{}, fmt.Errorf("Login: %w", err)
	}
	return model.LoginRep{
		Token: jwt,
		Name:  u.Edges.Profile.Name,
		UUID:  u.Edges.Profile.UUID,
	}, nil
}

func (w *WebService) Info(ctx context.Context, t *model.TokenClaims) (model.UserInfo, error) {
	u, err := w.client.User.Query().Where(user.ID(t.UID)).First(ctx)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("Info: %w", err)
	}
	isAdmin := utilsService.IsAdmin(u.State)
	return model.UserInfo{
		UID:     t.UID,
		UUID:    t.Subject,
		IsAdmin: isAdmin,
	}, nil
}

func (w *WebService) ChangePasswd(ctx context.Context, p model.ChangePasswd, t *model.TokenClaims) error {
	u, err := w.client.User.Query().Where(user.IDEQ(t.UID)).WithToken().First(ctx)
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	err = w.validatePass(ctx, u, p.Old)
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	pass, salt := utils.Argon2ID(p.New)
	if u.Edges.Token != nil {
		err := w.client.UserToken.UpdateOne(u.Edges.Token).AddTokenID(1).Exec(ctx)
		if err != nil {
			return fmt.Errorf("ChangePasswd: %w", err)
		}
	}
	err = w.cache.Del([]byte("auth" + strconv.Itoa(t.UID)))
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	err = w.client.User.UpdateOne(u).SetPassword(pass).SetSalt(salt).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	return nil
}

func (w *WebService) changeName(ctx context.Context, newName string, uid int, uuid string) error {
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

func (w *WebService) ChangeName(ctx context.Context, newName string, t *model.TokenClaims) error {
	err := w.changeName(ctx, newName, t.UID, t.Subject)
	if err != nil {
		return fmt.Errorf("ChangeName: %w", err)
	}
	return nil
}
