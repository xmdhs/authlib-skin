package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/predicate"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/db/ent/usertoken"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	"github.com/xmdhs/authlib-skin/service/auth"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

type AdminService struct {
	authService *auth.AuthService
	client      *ent.Client
	config      config.Config
	cache       cache.Cache
}

func NewAdminService(authService *auth.AuthService, client *ent.Client,
	config config.Config, cache cache.Cache) *AdminService {
	return &AdminService{
		authService: authService,
		client:      client,
		config:      config,
		cache:       cache,
	}
}

var ErrNotAdmin = errors.New("无权限")

func (w *AdminService) Auth(ctx context.Context, token string) (*model.TokenClaims, error) {
	t, err := w.authService.Auth(ctx, yggdrasil.ValidateToken{AccessToken: token}, false)
	if err != nil {
		return nil, fmt.Errorf("WebService.Auth: %w", err)
	}
	return t, nil
}

func (w *AdminService) IsAdmin(ctx context.Context, t *model.TokenClaims) error {
	u, err := w.client.User.Query().Where(user.ID(t.UID)).First(ctx)
	if err != nil {
		return fmt.Errorf("IsAdmin: %w", err)
	}
	if !auth.IsAdmin(u.State) {
		return fmt.Errorf("IsAdmin: %w", ErrNotAdmin)
	}
	return nil
}

func (w *AdminService) ListUser(ctx context.Context, page int, email, name string) ([]model.UserList, int, error) {
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
				IsAdmin: auth.IsAdmin(v.State),
			},
			Email:     v.Email,
			RegIp:     v.RegIP,
			Name:      v.Edges.Profile.Name,
			IsDisable: auth.IsDisable(v.State),
		})
	}

	uc, err := w.client.User.Query().Where(user.And(whereL...)).Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("ListUser: %w", err)
	}
	return ul, uc, nil
}

func (w *AdminService) EditUser(ctx context.Context, u model.EditUser, uid int) error {
	uuid := ""
	changePasswd := false
	err := utils.WithTx(ctx, w.client, func(tx *ent.Tx) error {
		upUser := tx.User.UpdateOneID(uid)

		if u.Email != "" {
			c, err := tx.User.Query().Where(user.Email(u.Email)).Count(ctx)
			if err != nil {
				return err
			}
			if c != 0 {
				return ErrExistUser
			}
			upUser = upUser.SetEmail(u.Email)
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
				err := utilsService.DelTexture(ctx, userProfile.ID, v, tx.Client(), w.config.TexturePath)
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
			state = auth.SetAdmin(state, *u.IsAdmin)
		}
		if u.IsDisable != nil {
			if *u.IsDisable {
				changePasswd = true
			}
			state = auth.SetDisable(state, *u.IsDisable)
		}
		if state != aUser.State {
			upUser = upUser.SetState(state)
		}
		if u.Password != "" {
			pass, salt := utils.Argon2ID(u.Password)
			upUser = upUser.SetPassword(pass).SetSalt(salt)
			changePasswd = true
		}

		err = upUser.Exec(ctx)
		if err != nil {
			return err
		}

		if changePasswd {
			err = tx.UserToken.Update().Where(usertoken.HasUserWith(user.ID(uid))).AddTokenID(1).Exec(ctx)
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
