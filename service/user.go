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
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

var (
	ErrExistUser = errors.New("邮箱已存在")
	ErrExitsName = errors.New("用户名已存在")
	ErrRegLimit  = errors.New("超过注册 ip 限制")
	ErrPassWord  = errors.New("错误的密码")
)

func (w *WebService) Reg(ctx context.Context, u model.User, ipPrefix, ip string) error {
	var userUuid string
	if w.config.OfflineUUID {
		userUuid = utils.UUIDGen(u.Name)
	} else {
		userUuid = strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	if w.config.Captcha.Type == "turnstile" {
		err := w.verifyTurnstile(ctx, u.CaptchaToken, ip)
		if err != nil {
			return fmt.Errorf("Reg: %w", err)
		}
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

	err := utils.WithTx(ctx, w.client, func(tx *ent.Tx) error {
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

func (w *WebService) Info(ctx context.Context, token string) (model.UserInfo, error) {
	t, err := utilsService.Auth(ctx, yggdrasil.ValidateToken{AccessToken: token}, w.client, w.cache, &w.prikey.PublicKey, false)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("Info: %w", err)
	}
	u, err := w.client.User.Query().Where(user.ID(t.UID)).First(ctx)
	if err != nil {
		return model.UserInfo{}, fmt.Errorf("Info: %w", err)
	}
	isAdmin := false
	if u.State&1 == 1 {
		isAdmin = true
	}
	return model.UserInfo{
		UID:     t.UID,
		UUID:    t.Subject,
		IsAdmin: isAdmin,
	}, nil
}

func (w *WebService) ChangePasswd(ctx context.Context, p model.ChangePasswd, token string) error {
	t, err := utilsService.Auth(ctx, yggdrasil.ValidateToken{AccessToken: token}, w.client, w.cache, &w.prikey.PublicKey, false)
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	u, err := w.client.User.Query().Where(user.IDEQ(t.UID)).WithToken().First(ctx)
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	if !utils.Argon2Compare(p.Old, u.Password, u.Salt) {
		return fmt.Errorf("ChangePasswd: %w", ErrPassWord)
	}
	pass, salt := utils.Argon2ID(p.New)

	err = w.client.UserToken.UpdateOne(u.Edges.Token).AddTokenID(1).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	w.cache.Del([]byte("auth" + strconv.Itoa(t.UID)))

	err = w.client.User.UpdateOne(u).SetPassword(pass).SetSalt(salt).Exec(ctx)
	if err != nil {
		return fmt.Errorf("ChangePasswd: %w", err)
	}
	return nil
}
