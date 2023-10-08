package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/xmdhs/authlib-skin/db/ent/predicate"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
)

var ErrNotAdmin = errors.New("无权限")

func (w *WebService) Auth(ctx context.Context, token string) (*model.TokenClaims, error) {
	t, err := utilsService.Auth(ctx, yggdrasil.ValidateToken{AccessToken: token}, w.client, w.cache, &w.prikey.PublicKey, false)
	if err != nil {
		return nil, fmt.Errorf("Auth: %w", err)
	}
	return t, nil
}

func (w *WebService) IsAdmin(ctx context.Context, t *model.TokenClaims) error {
	u, err := w.client.User.Query().Where(user.ID(t.UID)).First(ctx)
	if err != nil {
		return fmt.Errorf("IsAdmin: %w", err)
	}
	if !utilsService.IsAdmin(u.State) {
		return fmt.Errorf("IsAdmin: %w", ErrNotAdmin)
	}
	return nil
}

func (w *WebService) ListUser(ctx context.Context, page int, email, name string) ([]model.UserList, int, error) {
	whereL := []predicate.User{}
	if email != "" {
		whereL = append(whereL, user.EmailHasPrefix(email))
	}
	if name != "" {
		whereL = append(whereL, user.HasProfileWith(userprofile.NameHasPrefix(name)))
	}
	u, err := w.client.User.Query().WithProfile().
		Where(user.And(whereL...)).
		Limit(20).Offset((page - 1) * 20).All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("ListUser: %w", err)
	}
	ul := make([]model.UserList, 0, len(u))

	for _, v := range u {
		if v.Edges.Profile == nil {
			continue
		}
		ul = append(ul, model.UserList{
			UserInfo: model.UserInfo{
				UID:     v.ID,
				UUID:    v.Edges.Profile.UUID,
				IsAdmin: utilsService.IsAdmin(v.State),
			},
			Email: v.Email,
			RegIp: v.RegIP,
			Name:  v.Edges.Profile.Name,
		})
	}

	uc, err := w.client.User.Query().Where(user.And(whereL...)).Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("ListUser: %w", err)
	}
	return ul, uc, nil
}
