package yggdrasil

import (
	"context"
	"fmt"
	"time"

	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	"github.com/xmdhs/authlib-skin/service/auth"
)

type sessionWithIP struct {
	User model.TokenClaims
	IP   string
}

func (y *Yggdrasil) SessionJoin(ctx context.Context, s yggdrasil.Session, t *model.TokenClaims, ip string) error {
	if s.SelectedProfile != t.Subject {
		return fmt.Errorf("SessionJoin: %w", auth.ErrTokenInvalid)
	}
	err := cache.CacheHelp[sessionWithIP]{Cache: y.cache}.Put([]byte("session"+s.ServerID), sessionWithIP{
		User: *t,
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

	if up.UUID != sIP.User.Subject {
		return yggdrasil.UserInfo{}, fmt.Errorf("uuid 不相同")
	}

	u, err := y.GetProfile(ctx, up.UUID, false, host)
	if err != nil {
		return yggdrasil.UserInfo{}, fmt.Errorf("HasJoined: %w", err)
	}
	return u, nil
}
