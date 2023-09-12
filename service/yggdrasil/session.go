package yggdrasil

import (
	"context"
	"fmt"
	"time"

	"github.com/xmdhs/authlib-skin/db/cache"
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
	}, y.client, y.cache, &y.prikey.PublicKey, true)
	if err != nil {
		return fmt.Errorf("SessionJoin: %w", err)
	}
	if s.SelectedProfile != t.Subject {
		return fmt.Errorf("SessionJoin: %w", sutils.ErrTokenInvalid)
	}
	err = cache.CacheHelp[sessionWithIP]{Cache: y.cache}.Put([]byte("session"+s.ServerID), sessionWithIP{
		user: *t,
		IP:   ip,
	}, time.Now().Add(30*time.Second))
	if err != nil {
		return fmt.Errorf("SessionJoin: %w", err)
	}
	return nil
}

func (y *Yggdrasil) HasJoined(ctx context.Context, username, serverId string, ip string, host string) (yggdrasil.UserInfo, error) {
	sIP, err := cache.CacheHelp[sessionWithIP]{Cache: y.cache}.Get([]byte("session" + serverId))
	if err != nil {
		return yggdrasil.UserInfo{}, fmt.Errorf("HasJoined: %w", err)
	}
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