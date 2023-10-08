package yggdrasil

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/texture"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/model"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

var (
	ErrUUIDNotEq = errors.New("uuid 不相同")
)

func (y *Yggdrasil) delTexture(ctx context.Context, userProfileID int, textureType string) error {
	return utilsService.DelTexture(ctx, userProfileID, textureType, y.client, y.config)
}

func (y *Yggdrasil) DelTexture(ctx context.Context, uuid string, t *model.TokenClaims, textureType string) error {
	if uuid != t.Subject {
		return fmt.Errorf("PutTexture: %w", ErrUUIDNotEq)
	}
	up, err := y.client.UserProfile.Query().Where(userprofile.HasUserWith(user.ID(t.UID))).First(ctx)
	if err != nil {
		return fmt.Errorf("DelTexture: %w", err)
	}
	err = y.delTexture(ctx, up.ID, textureType)
	if err != nil {
		return fmt.Errorf("DelTexture: %w", err)
	}
	err = y.cache.Del([]byte("Profile" + uuid))
	if err != nil {
		return fmt.Errorf("DelTexture: %w", err)
	}
	return nil
}

func (y *Yggdrasil) PutTexture(ctx context.Context, t *model.TokenClaims, texturebyte []byte, model string, uuid string, textureType string) error {
	if uuid != t.Subject {
		return fmt.Errorf("PutTexture: %w", ErrUUIDNotEq)
	}

	up, err := y.client.UserProfile.Query().Where(userprofile.HasUserWith(user.ID(t.UID))).First(ctx)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}

	err = y.delTexture(ctx, up.ID, textureType)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}

	hashstr := getHash(texturebyte)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}
	u, err := y.client.User.Query().Where(user.HasProfileWith(userprofile.ID(up.ID))).Only(ctx)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}

	err = utils.WithTx(ctx, y.client, func(tx *ent.Tx) error {
		t, err := tx.Texture.Query().Where(texture.TextureHash(hashstr)).Only(ctx)
		if err != nil {
			var ne *ent.NotFoundError
			if !errors.As(err, &ne) {
				return err
			}
		}
		if t == nil {
			t, err = tx.Texture.Create().SetCreatedUser(u).SetTextureHash(hashstr).Save(ctx)
			if err != nil {
				return err
			}
		}
		err = tx.UserTexture.Create().SetTexture(t).SetType(textureType).SetUserProfile(up).SetVariant(model).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}
	err = createTextureFile(y.config.TexturePath, texturebyte, hashstr)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}
	err = y.cache.Del([]byte("Profile" + uuid))
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}
	return nil
}

func getHash(b []byte) string {
	hashed := sha256.Sum256(b)
	return hex.EncodeToString(hashed[:])
}

func createTextureFile(path string, b []byte, hashstr string) error {
	p := filepath.Join(path, hashstr[:2], hashstr[2:4], hashstr)
	err := os.MkdirAll(filepath.Dir(p), 0755)
	if err != nil {
		return fmt.Errorf("createTextureFile: %w", err)
	}
	f, err := os.Stat(p)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("createTextureFile: %w", err)
	}
	if f == nil {
		err := os.WriteFile(p, b, 0644)
		if err != nil {
			return fmt.Errorf("createTextureFile: %w", err)
		}
	}
	return nil
}
