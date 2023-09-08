package yggdrasil

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/texture"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
	"github.com/xmdhs/authlib-skin/db/ent/usertexture"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
	"lukechampine.com/blake3"
)

var (
	ErrUUIDNotEq = errors.New("uuid 不相同")
)

func (y *Yggdrasil) PutTexture(ctx context.Context, token string, texturebyte []byte, model string, uuid string, textureType string) error {
	t, err := utilsService.Auth(ctx, yggdrasil.ValidateToken{AccessToken: token}, y.client, &y.prikey.PublicKey, true)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}
	if uuid != t.Subject {
		return fmt.Errorf("PutTexture: %w", ErrUUIDNotEq)
	}

	up, err := y.client.UserProfile.Query().Where(userprofile.HasUserWith(user.ID(t.UID))).First(ctx)

	err = utils.WithTx(ctx, y.client, func(tx *ent.Tx) error {
		// 查找此用户该类型下是否已经存在皮肤
		tl, err := tx.UserTexture.Query().Where(usertexture.And(
			usertexture.UserProfileID(up.ID),
			usertexture.Type(textureType),
		)).ForUpdate().All(ctx)
		if err != nil {
			return err
		}
		if len(tl) == 0 {
			return nil
		}
		// 若存在，查找是否被引用
		for _, v := range tl {
			c, err := tx.UserTexture.Query().Where(usertexture.TextureID(v.TextureID)).Count(ctx)
			if err != nil {
				return err
			}
			if c == 1 {
				// 若没有其他用户使用该皮肤，删除文件和记录
				t, err := tx.Texture.Query().Where(texture.ID(v.TextureID)).Only(ctx)
				if err != nil {
					return err
				}
				path := filepath.Join(y.TexturePath, t.TextureHash[:2], t.TextureHash[2:4], t.TextureHash)
				err = os.Remove(path)
				if err != nil {
					return err
				}
				// Texture 表中删除记录
				err = tx.Texture.DeleteOneID(v.TextureID).Exec(ctx)
				if err != nil {
					return err
				}
			}
		}
		ids := lo.Map[*ent.UserTexture, int](tl, func(item *ent.UserTexture, index int) int {
			return item.UserProfileID
		})
		// 中间表删除记录
		_, err = tx.UserTexture.Delete().Where(usertexture.UserProfileIDIn(ids...)).Exec(ctx)
		return err
		// 小概率皮肤上传后，高并发时被此处清理。问题不大重新上传一遍就行。
		// 条件为使用一个独一无二的皮肤的用户，更换皮肤时，另一个用户同时更换自己的皮肤到这个独一无二的皮肤上。
	})

	hashstr, err := createTextureFile(y.config.TexturePath, texturebyte)
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
			err = tx.Texture.Create().SetCreatedUser(u).SetTextureHash(hashstr).Exec(ctx)
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
	return nil
}

func createTextureFile(path string, b []byte) (string, error) {
	hashed := blake3.Sum256(b)
	hashstr := hex.EncodeToString(hashed[:])
	p := filepath.Join(path, hashstr[:2], hashstr[2:4], hashstr)
	err := os.MkdirAll(filepath.Dir(p), 0755)
	if err != nil {
		return "", fmt.Errorf("createTextureFile: %w", err)
	}
	f, err := os.Stat(p)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("createTextureFile: %w", err)
	}
	if f == nil {
		err := os.WriteFile(p, b, 0644)
		if err != nil {
			return "", fmt.Errorf("createTextureFile: %w", err)
		}
	}
	return hashstr, nil
}
