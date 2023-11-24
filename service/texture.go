package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/texture"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/model"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

type TextureService struct {
	client *ent.Client
	config config.Config
	cache  cache.Cache
}

func NewTextureService(client *ent.Client, config config.Config, cache cache.Cache) *TextureService {
	return &TextureService{
		client: client,
		config: config,
		cache:  cache,
	}
}

func (w *TextureService) PutTexture(ctx context.Context, t *model.TokenClaims, texturebyte []byte, model string, textureType string) error {
	up, err := w.client.UserProfile.Query().Where(userprofile.HasUserWith(user.ID(t.UID))).First(ctx)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}
	err = utilsService.DelTexture(ctx, up.ID, textureType, w.client, w.config.TexturePath)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}

	hashstr := getHash(texturebyte)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}
	u, err := w.client.User.Query().Where(user.HasProfileWith(userprofile.ID(up.ID))).Only(ctx)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}

	err = utils.WithTx(ctx, w.client, func(tx *ent.Tx) error {
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
	err = createTextureFile(w.config.TexturePath, texturebyte, hashstr)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}
	err = w.cache.Del([]byte("Profile" + t.Subject))
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
