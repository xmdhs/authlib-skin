package yggdrasil

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	sutils "github.com/xmdhs/authlib-skin/service/utils"
)

type Yggdrasil struct {
	client *ent.Client
	cache  cache.Cache
	config config.Config
	prikey *rsa.PrivateKey

	pubStr func() string
}

func NewYggdrasil(client *ent.Client, cache cache.Cache, c config.Config, prikey *rsa.PrivateKey) *Yggdrasil {
	return &Yggdrasil{
		client: client,
		cache:  cache,
		config: c,
		prikey: prikey,
		pubStr: sync.OnceValue[string](func() string {
			derBytes := lo.Must(x509.MarshalPKIXPublicKey(&prikey.PublicKey))
			return base64.StdEncoding.EncodeToString(derBytes)
		}),
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

func newJwtToken(jwtKey *rsa.PrivateKey, tokenID, clientToken, UUID string, userID int) (string, error) {
	return sutils.NewJwtToken(jwtKey, tokenID, clientToken, UUID, userID)
}

func (y *Yggdrasil) Auth(ctx context.Context, t yggdrasil.ValidateToken) (*model.TokenClaims, error) {
	u, err := sutils.Auth(ctx, t, y.client, y.cache, &y.prikey.PublicKey, true)
	if err != nil {
		return nil, fmt.Errorf("ValidateToken: %w", err)
	}
	return u, nil
}

//go:embed yggdrasil_session_pubkey.der
var mojangPubKey []byte

var mojangPubKeyStr = sync.OnceValue(func() string {
	pub := lo.Must(x509.ParsePKIXPublicKey(mojangPubKey)).(*rsa.PublicKey)
	derBytes := lo.Must(x509.MarshalPKIXPublicKey(pub))
	return base64.StdEncoding.EncodeToString(derBytes)
})

func (y *Yggdrasil) PublicKeys(ctx context.Context) yggdrasil.PublicKeys {
	mojangPub := mojangPubKeyStr()
	myPub := y.pubStr()

	pl := []yggdrasil.PublicKeyList{{PublicKey: mojangPub}, {PublicKey: myPub}}

	return yggdrasil.PublicKeys{
		PlayerCertificateKeys: pl,
		ProfilePropertyKeys:   pl,
	}
}
