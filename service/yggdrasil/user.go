package yggdrasil

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/db/ent/usertoken"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	sutils "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

var (
	ErrRate     = errors.New("频率限制")
	ErrPassWord = errors.New("错误的密码或邮箱")
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
			return nil, fmt.Errorf("validatePass: %w", ErrPassWord)
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

	var utoken *ent.UserToken
	err = utils.WithTx(cxt, y.client, func(tx *ent.Tx) error {
		utoken, err = tx.User.QueryToken(u).ForUpdate().First(cxt)
		if err != nil {
			var nf *ent.NotFoundError
			if !errors.As(err, &nf) {
				return err
			}
		}
		if utoken == nil {
			ut, err := tx.UserToken.Create().SetTokenID(1).SetUser(u).Save(cxt)
			if err != nil {
				return err
			}
			utoken = ut
		}
		return nil
	})
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", err)
	}

	jwts, err := newJwtToken(y.prikey, strconv.FormatUint(utoken.TokenID, 10), clientToken, u.Edges.Profile.UUID, u.ID)
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", err)
	}

	p := yggdrasil.TokenProfile{
		ID:   u.Edges.Profile.UUID,
		Name: u.Edges.Profile.Name,
	}
	return yggdrasil.Token{
		AccessToken:       jwts,
		AvailableProfiles: []yggdrasil.TokenProfile{p},
		ClientToken:       clientToken,
		SelectedProfile:   p,
		User: yggdrasil.TokenUser{
			ID:         u.Edges.Profile.UUID,
			Properties: []any{},
		},
	}, nil
}

func (y *Yggdrasil) ValidateToken(ctx context.Context, t yggdrasil.ValidateToken) error {
	_, err := sutils.Auth(ctx, t, y.client, &y.prikey.PublicKey, true)
	if err != nil {
		return fmt.Errorf("ValidateToken: %w", err)
	}
	return nil
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
	return nil
}

func (y *Yggdrasil) Invalidate(ctx context.Context, accessToken string) error {
	t, err := sutils.Auth(ctx, yggdrasil.ValidateToken{AccessToken: accessToken}, y.client, &y.prikey.PublicKey, true)
	if err != nil {
		return fmt.Errorf("Invalidate: %w", err)
	}
	err = y.client.UserToken.Update().Where(usertoken.HasUserWith(user.ID(t.UID))).AddTokenID(1).Exec(ctx)
	if err != nil {
		return fmt.Errorf("Invalidate: %w", err)
	}
	return nil
}

func (y *Yggdrasil) Refresh(ctx context.Context, token yggdrasil.RefreshToken) (yggdrasil.Token, error) {
	t, err := sutils.Auth(ctx, yggdrasil.ValidateToken{AccessToken: token.AccessToken, ClientToken: token.ClientToken}, y.client, &y.prikey.PublicKey, false)
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Refresh: %w", err)
	}
	jwts, err := newJwtToken(y.prikey, t.Tid, t.CID, t.Subject, t.UID)
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", err)
	}

	up, err := y.client.UserProfile.Query().Where(userprofile.HasUserWith(user.ID(t.UID))).First(ctx)
	if err != nil {
		return yggdrasil.Token{}, fmt.Errorf("Authenticate: %w", err)
	}

	return yggdrasil.Token{
		AccessToken: jwts,
		ClientToken: t.CID,
		SelectedProfile: yggdrasil.TokenProfile{
			ID:   up.UUID,
			Name: up.Name,
		},
		User: yggdrasil.TokenUser{
			ID:         t.Subject,
			Properties: []any{},
		},
	}, nil
}
