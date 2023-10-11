package utils

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/usertoken"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	"github.com/xmdhs/authlib-skin/utils"
)

var (
	ErrTokenInvalid = errors.New("token 无效")
	ErrUserDisable  = errors.New("用户被禁用")
)

func Auth(ctx context.Context, t yggdrasil.ValidateToken, client *ent.Client, c cache.Cache, pubkey *rsa.PublicKey, tmpInvalid bool) (*model.TokenClaims, error) {
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
	tokenID, err := func() (uint64, error) {
		c := cache.CacheHelp[uint64]{Cache: c}
		key := []byte("auth" + strconv.Itoa(claims.UID))
		t, err := c.Get(key)
		if err != nil {
			return 0, err
		}
		if t != 0 {
			return t, nil
		}
		ut, err := client.UserToken.Query().Where(usertoken.HasUserWith(user.ID(claims.UID))).First(ctx)
		if err != nil {
			var ne *ent.NotFoundError
			if errors.As(err, &ne) {
				return 0, ErrTokenInvalid
			}
			return 0, err
		}
		return ut.TokenID, c.Put(key, ut.TokenID, time.Now().Add(20*time.Minute))
	}()
	if err != nil {
		return nil, fmt.Errorf("Auth: %w", err)
	}
	if strconv.FormatUint(tokenID, 10) != claims.Tid {
		return nil, fmt.Errorf("Auth: %w", ErrTokenInvalid)
	}
	return claims, nil
}

func CreateToken(ctx context.Context, u *ent.User, client *ent.Client, cache cache.Cache, jwtKey *rsa.PrivateKey, clientToken string) (string, error) {
	if IsDisable(u.State) {
		return "", fmt.Errorf("CreateToken: %w", ErrUserDisable)
	}
	var utoken *ent.UserToken
	err := utils.WithTx(ctx, client, func(tx *ent.Tx) error {
		var err error
		utoken, err = tx.User.QueryToken(u).ForUpdateA().First(ctx)
		if err != nil {
			var nf *ent.NotFoundError
			if !errors.As(err, &nf) {
				return err
			}
		}
		if utoken == nil {
			ut, err := tx.UserToken.Create().SetTokenID(1).SetUser(u).Save(ctx)
			if err != nil {
				return err
			}
			utoken = ut
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("CreateToken: %w", err)
	}
	err = cache.Del([]byte("auth" + strconv.Itoa(u.ID)))
	if err != nil {
		return "", fmt.Errorf("CreateToken: %w", err)
	}
	t, err := NewJwtToken(jwtKey, strconv.FormatUint(utoken.TokenID, 10), clientToken, u.Edges.Profile.UUID, u.ID)
	if err != nil {
		return "", fmt.Errorf("CreateToken: %w", err)
	}
	return t, nil
}

func NewJwtToken(jwtKey *rsa.PrivateKey, tokenID, clientToken, UUID string, userID int) (string, error) {
	claims := model.TokenClaims{
		Tid: tokenID,
		CID: clientToken,
		UID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * 24 * time.Hour)),
			Issuer:    "authlib-skin",
			Subject:   UUID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwts, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("newJwtToken: %w", err)
	}
	return jwts, nil
}

func IsAdmin(state int) bool {
	return state&1 == 1
}

func IsDisable(state int) bool {
	return state&2 == 2
}

func SetAdmin(state int, is bool) int {
	if is {
		return state | 1
	}
	return state & (state ^ 1)
}

func SetDisable(state int, is bool) int {
	if is {
		return state | 2
	}
	return state & (state ^ 2)
}
