package yggdrasil

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/texture"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/db/ent/usertoken"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	"github.com/xmdhs/authlib-skin/utils"
	"github.com/xmdhs/authlib-skin/utils/sign"
)

var (
	ErrRate        = errors.New("频率限制")
	ErrPassWord    = errors.New("错误的密码或邮箱")
	ErrNotUser     = errors.New("没有这个用户")
	ErrUserDisable = errors.New("用户被禁用")
)

func (y *Yggdrasil) validatePass(cxt context.Context, email, pass string) (*ent.User, error) {
	err := rate("validatePass"+email, y.cache, 10*time.Second, 3)
	if err != nil {
		return nil, fmt.Errorf("validatePass: %w", err)
	}
	u, err := y.client.User.Query().Where(user.EmailEQ(email)).WithProfile().First(cxt)
	if err != nil {
		var nf *ent.NotFoundError
		if errors.As(err, &nf) {
			return nil, fmt.Errorf("validatePass: %w", errors.Join(ErrPassWord, err))
		}
		return nil, fmt.Errorf("validatePass: %w", err)
	}
	if !utils.Argon2Compare(pass, u.Password, u.Salt) {
		return nil, fmt.Errorf("validatePass: %w", ErrPassWord)
	}
	return u, nil
}

func (y *Yggdrasil) Authenticate(cxt context.Context, auth yggdrasil.Authenticate) (yggdrasil.Token, error) {
	u, err := y.validatePass(cxt, auth.Username, auth.Password)
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", err)
	}

	clientToken := auth.ClientToken
	if clientToken == "" {
		clientToken = strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	jwts, err := y.authService.CreateToken(cxt, u, clientToken, u.Edges.Profile.UUID)
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", err)
	}

	if u.Edges.Profile == nil {
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", ErrUserDisable)
	}
	p := yggdrasil.UserInfo{
		ID:   u.Edges.Profile.UUID,
		Name: u.Edges.Profile.Name,
	}
	return yggdrasil.Token{
		AccessToken:       jwts,
		AvailableProfiles: []yggdrasil.UserInfo{p},
		ClientToken:       clientToken,
		SelectedProfile:   p,
		User: yggdrasil.TokenUserID{
			ID:         utils.UUIDGen(strconv.Itoa(u.ID)),
			Properties: []any{},
		},
	}, nil
}

func (y *Yggdrasil) SignOut(ctx context.Context, t yggdrasil.Pass) error {
	u, err := y.validatePass(ctx, t.Username, t.Password)
	if err != nil {
		return fmt.Errorf("SignOut: %w", err)
	}
	ut, err := y.client.UserToken.Query().Where(usertoken.HasUserWith(user.IDEQ(u.ID))).First(ctx)
	if err != nil {
		var nf *ent.NotFoundError
		if !errors.As(err, &nf) {
			return fmt.Errorf("SignOut: %w", err)
		}
		return nil
	}
	err = y.client.UserToken.UpdateOne(ut).AddTokenID(1).Exec(ctx)
	if err != nil {
		return fmt.Errorf("SignOut: %w", err)
	}
	err = y.cache.Del([]byte("auth" + strconv.Itoa(u.ID)))
	if err != nil {
		return fmt.Errorf("SignOut: %w", err)
	}
	return nil
}

func (y *Yggdrasil) Invalidate(ctx context.Context, t *model.TokenClaims) error {
	err := y.client.UserToken.Update().Where(usertoken.HasUserWith(user.ID(t.UID))).AddTokenID(1).Exec(ctx)
	if err != nil {
		return fmt.Errorf("Invalidate: %w", err)
	}
	err = y.cache.Del([]byte("auth" + strconv.Itoa(t.UID)))
	if err != nil {
		return fmt.Errorf("Invalidate: %w", err)
	}
	return nil
}

func (y *Yggdrasil) Refresh(ctx context.Context, t *model.TokenClaims) (yggdrasil.Token, error) {
	jwts, err := newJwtToken(y.prikey, t.Tid, t.CID, t.Subject, t.UID)
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Refresh: %w", err)
	}

	up, err := y.client.UserProfile.Query().Where(userprofile.HasUserWith(user.ID(t.UID))).First(ctx)
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Refresh: %w", err)
	}
	u := yggdrasil.UserInfo{ID: up.UUID, Name: up.Name}

	return yggdrasil.Token{
		AccessToken:       jwts,
		AvailableProfiles: []yggdrasil.UserInfo{u},
		ClientToken:       t.CID,
		SelectedProfile:   u,
		User: yggdrasil.TokenUserID{
			ID:         utils.UUIDGen(strconv.Itoa(t.UID)),
			Properties: []any{},
		},
	}, nil
}

func (y *Yggdrasil) GetProfile(ctx context.Context, uuid string, unsigned bool, host string) (yggdrasil.UserInfo, error) {
	c := cache.CacheHelp[yggdrasil.UserTextures]{Cache: y.cache}
	key := []byte("Profile" + uuid)
	ut, err := c.Get(key)
	if err != nil {
		return yggdrasil.UserInfo{}, fmt.Errorf("GetProfile: %w", err)
	}
	if ut.ProfileName == "" {
		up, err := y.client.UserProfile.Query().Where(userprofile.UUID(uuid)).WithUsertexture().Only(ctx)
		if err != nil {
			var nf *ent.NotFoundError
			if errors.As(err, &nf) {
				return yggdrasil.UserInfo{}, fmt.Errorf("GetProfile: %w", ErrNotUser)
			}
			return yggdrasil.UserInfo{}, fmt.Errorf("GetProfile: %w", err)
		}

		ut = yggdrasil.UserTextures{
			ProfileID:   up.UUID,
			ProfileName: up.Name,
			Textures:    map[string]yggdrasil.Textures{},
			Timestamp:   time.Now().UnixMilli(),
		}

		baseURl := func() string {
			if y.config.TextureBaseUrl == "" {
				u := &url.URL{}
				u.Host = host
				u.Scheme = "http"
				u.Path = "texture"
				return u.String()
			}
			return y.config.TextureBaseUrl
		}()

		for _, v := range up.Edges.Usertexture {
			dt, err := y.client.Texture.Query().Where(texture.ID(v.TextureID)).Only(ctx)
			if err != nil {
				return yggdrasil.UserInfo{}, fmt.Errorf("GetProfile: %w", err)
			}
			hashstr := dt.TextureHash
			t := yggdrasil.Textures{
				Url:      lo.Must1(url.JoinPath(baseURl, hashstr[:2], hashstr[2:4], hashstr)),
				Metadata: map[string]string{},
			}
			if v.Variant == "slim" {
				t.Metadata["model"] = "slim"
			}
			ut.Textures[strings.ToTitle(v.Type)] = t
		}
		err = c.Put(key, ut, time.Now().Add(30*time.Minute))
		if err != nil {
			return yggdrasil.UserInfo{}, fmt.Errorf("GetProfile: %w", err)
		}
	}

	texturesBase64 := ut.Base64()

	pl := []yggdrasil.UserProperties{{
		Name:  "textures",
		Value: texturesBase64,
	}}

	if !unsigned {
		s := sign.NewAuthlibSignWithKey(y.prikey)
		for i, v := range pl {
			sign, err := s.Sign([]byte(v.Value))
			if err != nil {
				return yggdrasil.UserInfo{}, fmt.Errorf("GetProfile: %w", err)
			}
			pl[i].Signature = sign
		}
	}

	uinfo := yggdrasil.UserInfo{
		ID:         ut.ProfileID,
		Name:       ut.ProfileName,
		Properties: pl,
	}

	return uinfo, nil
}

func (y *Yggdrasil) BatchProfile(ctx context.Context, names []string) ([]yggdrasil.UserInfo, error) {
	pl, err := y.client.UserProfile.Query().Where(userprofile.NameIn(names...)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("BatchProfile: %w", err)
	}
	return lo.Map[*ent.UserProfile, yggdrasil.UserInfo](pl, func(item *ent.UserProfile, index int) yggdrasil.UserInfo {
		return yggdrasil.UserInfo{
			ID:   item.UUID,
			Name: item.Name,
		}
	}), nil
}

// publicKey 为 PKIX，但要求 pem type 为 RSA PUBLIC KEY
// privateKey 为 PKCS #8， pem type 为 RSA PUBLIC KEY
// 签名使用 rsaWIthsha1

func (y *Yggdrasil) PlayerCertificates(ctx context.Context, t *model.TokenClaims) (yggdrasil.Certificates, error) {
	rsa2048, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return yggdrasil.Certificates{}, fmt.Errorf("PlayerCertificates: %w", err)
	}

	s := sign.NewAuthlibSignWithKey(rsa2048)
	priKey := lo.Must(s.GetPriKey())
	pubKey := lo.Must(s.GetPKIXPubKey())

	expiresAt := time.Now().Add(24 * time.Hour)
	expiresAtUnix := expiresAt.UnixMilli()

	pubV2, err := publicKeySignatureV2(&rsa2048.PublicKey, t.Subject, expiresAtUnix)
	if err != nil {
		return yggdrasil.Certificates{}, fmt.Errorf("PlayerCertificates: %w", err)
	}
	pub := publicKeySignature(pubKey, expiresAtUnix)

	servicePri := sign.NewAuthlibSignWithKey(y.prikey)

	pubV2Base64, err := servicePri.Sign(pubV2)
	if err != nil {
		return yggdrasil.Certificates{}, fmt.Errorf("PlayerCertificates: %w", err)
	}
	pubBase64, err := servicePri.Sign(pub)
	if err != nil {
		return yggdrasil.Certificates{}, fmt.Errorf("PlayerCertificates: %w", err)
	}

	return yggdrasil.Certificates{
		ExpiresAt: expiresAt.Format(time.RFC3339Nano),
		KeyPair: yggdrasil.CertificatesKeyPair{
			PrivateKey: priKey,
			PublicKey:  pubKey,
		},
		PublicKeySignature:   pubBase64,
		PublicKeySignatureV2: pubV2Base64,
		RefreshedAfter:       time.Now().Format(time.RFC3339Nano),
	}, nil

}

func publicKeySignatureV2(key *rsa.PublicKey, uuid string, expiresAt int64) ([]byte, error) {
	bf := &bytes.Buffer{}
	u := big.Int{}
	u.SetString(uuid, 16)
	bf.Write(u.Bytes())

	eb := make([]byte, 8)
	binary.BigEndian.PutUint64(eb, uint64(expiresAt))
	bf.Write(eb)
	pubKey, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, fmt.Errorf("publicKeySignatureV2: %w", err)
	}
	bf.Write(pubKey)
	return bf.Bytes(), nil
}

func publicKeySignature(key string, expiresAt int64) []byte {
	bf := &bytes.Buffer{}
	bf.WriteString(strconv.FormatInt(expiresAt, 10))
	bf.WriteString(key)
	return bf.Bytes()
}
