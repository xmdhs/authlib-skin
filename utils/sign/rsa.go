package sign

import (
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
	p, err := x509.ParsePKCS8PrivateKey(b.Bytes)
	if err != nil {
		return nil, fmt.Errorf("NewAuthlibSign: %w", err)
	}
	priv, ok := p.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("NewAuthlibSign: %w", ErrPem)
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

func (a *AuthlibSign) getPKIXPubKey(typeStr string) (string, error) {
	derBytes, err := x509.MarshalPKIXPublicKey(&a.key.PublicKey)
	if err != nil {
		return "", fmt.Errorf("getPKIXPubKey: %w", err)
	}
	pemKey := &pem.Block{
		Type:  typeStr,
		Bytes: derBytes,
	}
	return string(pem.EncodeToMemory(pemKey)), nil
}

// PKIX PUBLIC KEY
func (a *AuthlibSign) GetPKIXPubKeyWithOutRsa() (string, error) {
	return a.getPKIXPubKey("PUBLIC KEY")
}

// PKIX RSA PUBLIC KEY
func (a *AuthlibSign) GetPKIXPubKey() (string, error) {
	return a.getPKIXPubKey("RSA PUBLIC KEY")
}

// PKCS #8
func (a *AuthlibSign) GetPriKey() (string, error) {
	derBytes, err := x509.MarshalPKCS8PrivateKey(a.key)
	if err != nil {
		return "", fmt.Errorf("GetPriKey: %w", err)
	}
	pemKey := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derBytes,
	}
	return string(pem.EncodeToMemory(pemKey)), nil
}

func (a *AuthlibSign) Sign(data []byte) (string, error) {
	hashed := sha1.Sum(data)
	signature, err := rsa.SignPKCS1v15(nil, a.key, crypto.SHA1, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}
