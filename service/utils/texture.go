package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/texture"
	"github.com/xmdhs/authlib-skin/db/ent/usertexture"
)

func DelTexture(ctx context.Context, userProfileID int, textureType string, client *ent.Client, texturePath string) error {
	// 查找此用户该类型下是否已经存在皮肤
	tl, err := client.UserTexture.Query().Where(usertexture.And(
		usertexture.UserProfileID(userProfileID),
		usertexture.Type(textureType),
	)).All(ctx)
	if err != nil {
		return fmt.Errorf("DelTexture: %w", err)
	}
	if len(tl) == 0 {
		return nil
	}
	// 若存在，查找是否被引用
	for _, v := range tl {
		c, err := client.UserTexture.Query().Where(usertexture.TextureID(v.TextureID)).Count(ctx)
		if err != nil {
			return fmt.Errorf("DelTexture: %w", err)
		}
		if c == 1 {
			// 若没有其他用户使用该皮肤，删除文件和记录
			t, err := client.Texture.Query().Where(texture.ID(v.TextureID)).Only(ctx)
			if err != nil {
				var nf *ent.NotFoundError
				if errors.As(err, &nf) {
					continue
				}
				return fmt.Errorf("DelTexture: %w", err)
			}
			path := filepath.Join(texturePath, t.TextureHash[:2], t.TextureHash[2:4], t.TextureHash)
			err = os.Remove(path)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("DelTexture: %w", err)
			}
			// Texture 表中删除记录
			err = client.Texture.DeleteOneID(v.TextureID).Exec(ctx)
			if err != nil {
				return fmt.Errorf("DelTexture: %w", err)
			}
		}
	}
	ids := lo.Map[*ent.UserTexture, int](tl, func(item *ent.UserTexture, index int) int {
		return item.ID
	})
	// 中间表删除记录
	// UserProfile 上没有于此相关的字段，所以无需操作
	_, err = client.UserTexture.Delete().Where(usertexture.IDIn(ids...)).Exec(ctx)
	// 小概率皮肤上传后，高并发时被此处清理。问题不大重新上传一遍就行。
	// 条件为使用一个独一无二的皮肤的用户，更换皮肤时，另一个用户同时更换自己的皮肤到这个独一无二的皮肤上。
	if err != nil {
		return fmt.Errorf("DelTexture: %w", err)
	}
	return nil
}
