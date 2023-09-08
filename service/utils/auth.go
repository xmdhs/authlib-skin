package utils

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/usertoken"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
)

var (
	ErrTokenInvalid = errors.New("token 无效")
)

func Auth(ctx context.Context, t yggdrasil.ValidateToken, client *ent.Client, pubkey *rsa.PublicKey, tmpInvalid bool) (*model.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(t.AccessToken, &model.TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return pubkey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Auth: %w", errors.Join(err, ErrTokenInvalid))
	}

	claims, ok := token.Claims.(*model.TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Auth: %w", ErrTokenInvalid)
	}
	if t.ClientToken != "" && t.ClientToken != claims.CID {
		return nil, fmt.Errorf("Auth: %w", ErrTokenInvalid)
	}

	if tmpInvalid {
		it, err := claims.GetIssuedAt()
		if err != nil {
			return nil, fmt.Errorf("Auth: %w", errors.Join(err, ErrTokenInvalid))
		}
		et, err := claims.GetExpirationTime()
		if err != nil {
			return nil, fmt.Errorf("Auth: %w", errors.Join(err, ErrTokenInvalid))
		}
		invalidTime := it.Add(et.Time.Sub(it.Time) / 2)
		if time.Now().After(invalidTime) {
			return nil, fmt.Errorf("Auth: %w", ErrTokenInvalid)
		}
	}

	ut, err := client.UserToken.Query().Where(usertoken.HasUserWith(user.ID(claims.UID))).First(ctx)
	if err != nil {
		var ne *ent.NotFoundError
		if !errors.As(err, &ne) {
			return nil, fmt.Errorf("Auth: %w", errors.Join(err, ErrTokenInvalid))
		}
		return nil, fmt.Errorf("Auth: %w", err)
	}
	if strconv.FormatUint(ut.TokenID, 10) != claims.Tid {
		return nil, fmt.Errorf("Auth: %w", ErrTokenInvalid)
	}
	return claims, nil
}
