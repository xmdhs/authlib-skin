package sign

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type AuthlibSign struct {
	key *rsa.PrivateKey
}

var ErrPem = errors.New("错误的证书")

func NewAuthlibSign(prikey []byte) (*AuthlibSign, error) {
	b, _ := pem.Decode(prikey)
	if b == nil {
		return nil, fmt.Errorf("NewAuthlibSign: %w", ErrPem)
	}
	priv, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
		return nil, fmt.Errorf("NewAuthlibSign: %w", err)
	}
	return &AuthlibSign{
		key: priv,
	}, nil
}

func NewAuthlibSignWithKey(key *rsa.PrivateKey) *AuthlibSign {
	return &AuthlibSign{
		key: key,
	}
}

func (a *AuthlibSign) GetKey() *rsa.PrivateKey {
	return a.key
}

func (a *AuthlibSign) GetPubKey() (string, error) {
	derBytes := x509.MarshalPKCS1PublicKey(&a.key.PublicKey)
	pemKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}
	bw := &bytes.Buffer{}
	err := pem.Encode(bw, pemKey)
	if err != nil {
		return "", fmt.Errorf("GetPubKey: %w", err)
	}
	return bw.String(), nil
}

func (a *AuthlibSign) GetPriKey() (string, error) {
	derBytes := x509.MarshalPKCS1PrivateKey(a.key)
	pemKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derBytes,
	}
	bw := &bytes.Buffer{}
	err := pem.Encode(bw, pemKey)
	if err != nil {
		return "", fmt.Errorf("GetPubKey: %w", err)
	}
	return bw.String(), nil
}

func (a *AuthlibSign) Sign(data []byte) (string, error) {
	hashed := sha1.Sum(data)
	signature, err := rsa.SignPKCS1v15(nil, a.key, crypto.SHA1, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}
