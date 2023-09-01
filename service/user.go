package service

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/mysql"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/utils"
)

var ErrExistUser = errors.New("用户已存在")

func Reg(ctx context.Context, u model.User, q mysql.Querier, db *sql.DB, snow *snowflake.Node,
	c config.Config,
) error {
	ou, err := q.GetUserByEmail(ctx, u.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("Reg: %w", err)
	}
	if ou.Email != "" {
		return fmt.Errorf("Reg: %w", ErrExistUser)
	}
	err = utils.WithTx(ctx, &sql.TxOptions{}, q, db, func(q mysql.Querier) error {
		p, s := utils.Argon2ID(u.Password)
		userID := snow.Generate().Int64()
		_, err := q.CreateUser(ctx, mysql.CreateUserParams{
			ID:       userID,
			Email:    u.Email,
			Password: p,
			Salt:     s,
			State:    0,
			RegTime:  time.Now().Unix(),
		})
		if err != nil {
			return err
		}
		var userUuid string
		if c.OfflineUUID {
			userUuid = uuidGen(u.Name)
		} else {
			userUuid = strings.ReplaceAll(uuid.New().String(), "-", "")
		}

		_, err = q.CreateUserProfile(ctx, mysql.CreateUserProfileParams{
			UserID: userID,
			Name:   u.Name,
			Uuid:   userUuid,
		})
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
