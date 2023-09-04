package yggdrasil

import (
	"context"
	"fmt"

	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	"github.com/xmdhs/authlib-skin/service/utils"
)

func (y *Yggdrasil) ValidateToken(ctx context.Context, t yggdrasil.ValidateToken) error {
	err := utils.Auth(ctx, t, y.client, y.config.JwtKey)
	if err != nil {
		return fmt.Errorf("ValidateToken: %w", err)
	}
	return nil
}
