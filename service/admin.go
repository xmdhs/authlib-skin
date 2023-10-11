package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/predicate"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

var ErrNotAdmin = errors.New("无权限")

func (w *WebService) Auth(ctx context.Context, token string) (*model.TokenClaims, error) {
	t, err := utilsService.Auth(ctx, yggdrasil.ValidateToken{AccessToken: token}, w.client, w.cache, &w.prikey.PublicKey, false)
	if err != nil {
		return nil, fmt.Errorf("WebService.Auth: %w", err)
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
			Email:     v.Email,
			RegIp:     v.RegIP,
			Name:      v.Edges.Profile.Name,
			IsDisable: utilsService.IsDisable(v.State),
		})
	}

	uc, err := w.client.User.Query().Where(user.And(whereL...)).Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("ListUser: %w", err)
	}
	return ul, uc, nil
}

func (w *WebService) EditUser(ctx context.Context, u model.EditUser, uid int) error {
	uuid := ""
	changePasswd := false
	err := utils.WithTx(ctx, w.client, func(tx *ent.Tx) error {
		if u.Email != "" {
			c, err := tx.User.Query().Where(user.Email(u.Email)).Count(ctx)
			if err != nil {
				return err
			}
			if c != 0 {
				return ErrExistUser
			}
			err = tx.User.UpdateOneID(uid).SetEmail(u.Email).Exec(ctx)
			if err != nil {
				return err
			}
		}

		if u.Name != "" {
			c, err := tx.UserProfile.Query().Where(userprofile.Name(u.Name)).Count(ctx)
			if err != nil {
				return err
			}
			if c != 0 {
				return ErrExitsName
			}
			err = tx.UserProfile.Update().Where(userprofile.HasUserWith(user.ID(uid))).SetName(u.Name).Exec(ctx)
			if err != nil {
				return err
			}
		}

		if u.DelTextures {
			userProfile, err := tx.UserProfile.Query().Where(userprofile.ID(uid)).First(ctx)
			if err != nil {
				return err
			}
			uuid = userProfile.UUID
			tl := []string{"skin", "cape"}
			for _, v := range tl {
				err := utilsService.DelTexture(ctx, userProfile.ID, v, w.client, w.config)
				if err != nil {
					return err
				}
			}
		}

		aUser, err := tx.User.Get(ctx, uid)
		if err != nil {
			return err
		}

		state := aUser.State
		if u.IsAdmin != nil {
			state = utilsService.SetAdmin(state, *u.IsAdmin)
		}
		if u.IsDisable != nil {
			state = utilsService.SetDisable(state, *u.IsDisable)
		}
		if state != aUser.State {
			err := tx.User.UpdateOneID(uid).SetState(state).Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("EditUser: %w", err)
	}
	if uuid != "" {
		err = w.cache.Del([]byte("Profile" + uuid))
		if err != nil {
			return fmt.Errorf("EditUser: %w", err)
		}
	}
	if changePasswd {
		err = w.cache.Del([]byte("auth" + strconv.Itoa(uid)))
		if err != nil {
			return fmt.Errorf("EditUser: %w", err)
		}
	}
	return nil
}
