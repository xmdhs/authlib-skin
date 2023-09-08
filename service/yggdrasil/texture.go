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
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
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

	up, err := y.client.UserProfile.Query().Where(userprofile.UUIDEQ(uuid)).WithUser().First(ctx)

	err = utils.WithTx(ctx, y.client, func(tx *ent.Tx) error {
		// 查找此用户该类型下是否已经存在皮肤
		tl, err := tx.UserProfile.QueryTexture(up).Where(texture.TypeEQ(textureType)).ForUpdate().All(ctx)
		if err != nil {
			return err
		}
		if len(tl) == 0 {
			return nil
		}
		// 若存在，查找是否被引用
		for _, v := range tl {
			c, err := tx.UserProfile.Query().Where(userprofile.HasTextureWith(texture.IDEQ(v.ID))).Count(ctx)
			if err != nil {
				return err
			}
			if c == 1 {
				// 若没有其他用户使用该皮肤，删除文件和记录
				path := filepath.Join(y.TexturePath, v.TextureHash[:2], v.TextureHash[2:4], v.TextureHash)
				err = os.Remove(path)
				if err != nil {
					return err
				}
				err = tx.Texture.DeleteOneID(v.ID).Exec(ctx)
				if err != nil {
					return err
				}
			}
		}
		ids := lo.Map[*ent.Texture, int](tl, func(item *ent.Texture, index int) int {
			return item.ID
		})
		return tx.UserProfile.UpdateOne(up).RemoveTextureIDs(ids...).Exec(ctx)
	})

	hashstr, err := createTextureFile(y.config.TexturePath, texturebyte)
	if err != nil {
		return fmt.Errorf("PutTexture: %w", err)
	}

	var textureEnt *ent.Texture
	err = utils.WithTx(ctx, y.client, func(tx *ent.Tx) error {
		textureEnt, err = tx.Texture.Query().Where(texture.TextureHashEQ(hashstr)).ForUpdate().First(ctx)
		var nr *ent.NotFoundError
		if err != nil && !errors.As(err, &nr) {
			return err
		}
		if textureEnt == nil {
			textureEnt, err = tx.Texture.Create().SetCreatedUser(up.Edges.User).
				SetTextureHash(hashstr).
				SetType(textureType).
				SetVariant(model).
				Save(ctx)
			if err != nil {
				return err
			}
		}
		return tx.UserProfile.UpdateOne(up).AddTexture(textureEnt).Exec(ctx)
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
