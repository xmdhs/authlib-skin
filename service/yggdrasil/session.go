package yggdrasil

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	sutils "github.com/xmdhs/authlib-skin/service/utils"
)

type sessionWithIP struct {
	user model.TokenClaims
	IP   string
}

func (y *Yggdrasil) SessionJoin(ctx context.Context, s yggdrasil.Session, ip string) error {
	t, err := sutils.Auth(ctx, yggdrasil.ValidateToken{
		AccessToken: s.AccessToken,
	}, y.client, &y.prikey.PublicKey, true)
	if err != nil {
		return fmt.Errorf("SessionJoin: %w", err)
	}
	if s.SelectedProfile != t.Subject {
		return fmt.Errorf("SessionJoin: %w", sutils.ErrTokenInvalid)
	}
	err = y.cache.Put([]byte("session"+s.ServerID), lo.Must1(json.Marshal(sessionWithIP{
		user: *t,
		IP:   ip,
	})), time.Now().Add(30*time.Second))
	if err != nil {
		return fmt.Errorf("SessionJoin: %w", err)
	}
	return nil
}

func (y *Yggdrasil) HasJoined(ctx context.Context, username, serverId string, ip string, host string) (yggdrasil.UserInfo, error) {
	b := lo.Must1(y.cache.Get([]byte("session" + serverId)))
	sIP := sessionWithIP{}
	lo.Must0(json.Unmarshal(b, &sIP))

	if ip != "" && ip != sIP.IP {
		return yggdrasil.UserInfo{}, fmt.Errorf("ip 不相同")
	}

	up, err := y.client.UserProfile.Query().Where(userprofile.Name(username)).Only(ctx)
	if err != nil {
		return yggdrasil.UserInfo{}, fmt.Errorf("HasJoined: %w", err)
	}

	if up.UUID != sIP.user.Subject {
		return yggdrasil.UserInfo{}, fmt.Errorf("uuid 不相同")
	}

	u, err := y.GetProfile(ctx, up.UUID, false, host)
	if err != nil {
		return yggdrasil.UserInfo{}, fmt.Errorf("HasJoined: %w", err)
	}
	return u, nil
}
