package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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
)

func (w *WebService) Reg(ctx context.Context, u model.User, ip string) error {
	var userUuid string
	if w.config.OfflineUUID {
		userUuid = uuidGen(u.Name)
	} else {
		userUuid = strings.ReplaceAll(uuid.New().String(), "-", "")
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
			SetRegIP("").
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

func uuidGen(t string) string {
	data := []byte("OfflinePlayer:" + t)
	h := md5.New()
	h.Write(data)
	uuid := h.Sum(nil)
	uuid[6] = (uuid[6] & 0x0f) | uint8((3&0xf)<<4)
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return hex.EncodeToString(uuid)
}
