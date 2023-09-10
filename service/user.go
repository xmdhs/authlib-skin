package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/utils"
)

var (
	ErrExistUser = errors.New("邮箱已存在")
	ErrExitsName = errors.New("用户名已存在")
	ErrRegLimit  = errors.New("超过注册 ip 限制")
)

func (w *WebService) Reg(ctx context.Context, u model.User, ip string) error {
	var userUuid string
	if w.config.OfflineUUID {
		userUuid = utils.UUIDGen(u.Name)
	} else {
		userUuid = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	p, s := utils.Argon2ID(u.Password)

	if w.config.MaxIpUser != 0 {
		c, err := w.client.User.Query().Where(user.RegIPEQ(ip)).Count(ctx)
		if err != nil {
			return fmt.Errorf("Reg: %w", err)
		}
		if c >= w.config.MaxIpUser {
			return fmt.Errorf("Reg: %w", ErrRegLimit)
		}
	}

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
			SetRegIP(ip).
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
