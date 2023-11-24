package service

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service/auth"
	"github.com/xmdhs/authlib-skin/service/captcha"
	"github.com/xmdhs/authlib-skin/utils"
)

type WebService struct {
	config         config.Config
	client         *ent.Client
	httpClient     *http.Client
	cache          cache.Cache
	prikey         *rsa.PrivateKey
	authService    *auth.AuthService
	captchaService *captcha.CaptchaService
}

func NewWebService(c config.Config, e *ent.Client, hc *http.Client,
	cache cache.Cache, prikey *rsa.PrivateKey, authService *auth.AuthService, captchaService *captcha.CaptchaService) *WebService {
	return &WebService{
		config:         c,
		client:         e,
		httpClient:     hc,
		cache:          cache,
		prikey:         prikey,
		authService:    authService,
		captchaService: captchaService,
	}
}

func (w *WebService) validatePass(ctx context.Context, u *ent.User, password string) error {
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
