package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xmdhs/authlib-skin/model"
)

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
		return "", fmt.Errorf("NewJwtToken: %w", err)
	}
	return jwts, nil
}
