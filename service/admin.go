package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
)

var ErrNotAdmin = errors.New("无权限")

func (w *WebService) IsAdmin(ctx context.Context, token string) error {
	t, err := utilsService.Auth(ctx, yggdrasil.ValidateToken{AccessToken: token}, w.client, w.cache, &w.prikey.PublicKey, false)
	if err != nil {
		return fmt.Errorf("IsAdmin: %w", err)
	}
	u, err := w.client.User.Query().Where(user.ID(t.UID)).First(ctx)
	if err != nil {
		return fmt.Errorf("IsAdmin: %w", err)
	}
	if !utilsService.IsAdmin(u.State) {
		return fmt.Errorf("IsAdmin: %w", ErrNotAdmin)
	}
	return nil
}

func (w *WebService) ListUser(ctx context.Context, page int) ([]model.UserList, int, error) {
	u, err := w.client.User.Query().WithProfile().Limit(20).Offset((page - 1) * 20).All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("ListUser: %w", err)
	}
	ul := make([]model.UserList, 0, len(u))

	for _, v := range u {
		ul = append(ul, model.UserList{
			UserInfo: model.UserInfo{
				UID:     v.ID,
				UUID:    v.Edges.Profile.UUID,
				IsAdmin: utilsService.IsAdmin(v.State),
			},
			Email: v.Email,
			RegIp: v.RegIP,
		})
	}

	uc, err := w.client.User.Query().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("ListUser: %w", err)
	}
	return ul, uc, nil
}
