package yggdrasil

import (
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/model"
)

type Yggdrasil struct {
	client *ent.Client
	cache  cache.Cache
	config config.Config
	prikey *rsa.PrivateKey
}

func NewYggdrasil(client *ent.Client, cache cache.Cache, c config.Config, prikey *rsa.PrivateKey) *Yggdrasil {
	return &Yggdrasil{
		client: client,
		cache:  cache,
		config: c,
		prikey: prikey,
	}
}

func rate(k string, c cache.Cache, d time.Duration, count uint) error {
	key := []byte(k)
	v, err := c.Get([]byte(key))
	if err != nil {
		return fmt.Errorf("rate: %w", err)
	}
	if v == nil {
		err := putUint(1, c, key, d)
		if err != nil {
			return fmt.Errorf("rate: %w", err)
		}
		return nil
	}
	n := binary.BigEndian.Uint64(v)
	if n > uint64(count) {
		return fmt.Errorf("rate: %w", ErrRate)
	}
	err = putUint(n+1, c, key, d)
	if err != nil {
		return fmt.Errorf("rate: %w", err)
	}
	return nil
}

func putUint(n uint64, c cache.Cache, key []byte, d time.Duration) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	err := c.Put(key, b, time.Now().Add(d))
	if err != nil {
		return fmt.Errorf("rate: %w", err)
	}
	return nil
}

func newJwtToken(jwtKey *rsa.PrivateKey, tokenID, clientToken, UUID string) (string, error) {
	claims := model.TokenClaims{
		Tid: tokenID,
		CID: clientToken,
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
