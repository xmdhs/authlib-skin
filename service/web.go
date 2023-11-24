package service

import (
	"context"
	"fmt"

	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/utils"
)

type WebService struct {
	config config.Config
}

func NewWebService(c config.Config) *WebService {
	return &WebService{
		config: c,
	}
}

func validatePass(ctx context.Context, u *ent.User, password string) error {
	if !utils.Argon2Compare(password, u.Password, u.Salt) {
		return fmt.Errorf("validatePass: %w", ErrPassWord)
	}
	return nil
}

func (w *WebService) GetConfig(ctx context.Context) model.Config {
	return model.Config{
		Captcha: model.Captcha{
			Type:    w.config.Captcha.Type,
			SiteKey: w.config.Captcha.SiteKey,
		},
		ServerName:      w.config.ServerName,
		AllowChangeName: !w.config.OfflineUUID,
	}
}
