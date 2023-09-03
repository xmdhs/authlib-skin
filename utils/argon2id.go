package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func Argon2ID(pass string) (password string, salt string) {
	s := make([]byte, 16)
	_, err := rand.Read(s)
	if err != nil {
		panic(err)
	}
	b := argon2.IDKey([]byte(pass), s, 1, 64*1024, 1, 32)
	return base64.StdEncoding.EncodeToString(b), base64.StdEncoding.EncodeToString(s)
}

func Argon2Compare(pass, hashPass string, salt string) bool {
	s, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false
	}
	b := argon2.IDKey([]byte(pass), s, 1, 64*1024, 1, 32)
	hb, err := base64.StdEncoding.DecodeString(hashPass)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(b, hb) == 1
}
