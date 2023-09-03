package yggdrasil

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	"github.com/xmdhs/authlib-skin/utils"
)

var (
	ErrRate     = errors.New("频率限制")
	ErrPassWord = errors.New("错误的密码或邮箱")
)

func (y *Yggdrasil) Authenticate(cxt context.Context, auth yggdrasil.Authenticate) (yggdrasil.Token, error) {
	key := []byte("Authenticate" + auth.Username)

	v, err := y.cache.Get(key)
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", err)
	}
	if v != nil {
		u := binary.BigEndian.Uint64(v)
		t := time.Unix(int64(u), 0)
		if time.Now().Before(t) {
			return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", ErrRate)
		}
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(time.Now().Add(10*time.Second).Unix()))
	err = y.cache.Put(key, b, time.Now().Add(20*time.Second))
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", err)
	}

	u, err := y.client.User.Query().Where(user.EmailEQ(auth.Username)).WithProfile().WithToken().First(cxt)
	if err != nil {
		var nf *ent.NotFoundError
		if errors.As(err, &nf) {
			return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", ErrPassWord)
		}
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", err)
	}
	if !utils.Argon2Compare(auth.Password, u.Password, u.Salt) {
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", ErrPassWord)
	}
	clientToken := auth.ClientToken
	if clientToken == "" {
		clientToken = strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	claims := model.TokenClaims{
		Tid: u.Edges.Profile.UUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * 24 * time.Hour)),
			Issuer:    "authlib-skin",
			Subject:   u.Edges.Profile.UUID,
		},
	}
	_ = claims

	return yggdrasil.Token{}, nil
}
