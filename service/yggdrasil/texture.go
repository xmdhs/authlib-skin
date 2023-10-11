package yggdrasil

import (
	"context"
	"errors"
	"fmt"

	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/model"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
)

var (
	ErrUUIDNotEq = errors.New("uuid 不相同")
)

func (y *Yggdrasil) delTexture(ctx context.Context, userProfileID int, textureType string) error {
	return utilsService.DelTexture(ctx, userProfileID, textureType, y.client, y.config.TexturePath)
}

func (y *Yggdrasil) DelTexture(ctx context.Context, t *model.TokenClaims, textureType string) error {
	up, err := y.client.UserProfile.Query().Where(userprofile.HasUserWith(user.ID(t.UID))).First(ctx)
	if err != nil {
		return fmt.Errorf("DelTexture: %w", err)
	}
	err = y.delTexture(ctx, up.ID, textureType)
	if err != nil {
		return fmt.Errorf("DelTexture: %w", err)
	}
	err = y.cache.Del([]byte("Profile" + t.Subject))
	if err != nil {
		return fmt.Errorf("DelTexture: %w", err)
	}
	return nil
}
