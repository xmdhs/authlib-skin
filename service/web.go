package service

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/service/auth"
	"github.com/xmdhs/authlib-skin/utils"
)

type WebService struct {
	config      config.Config
	client      *ent.Client
	httpClient  *http.Client
	cache       cache.Cache
	prikey      *rsa.PrivateKey
	authService *auth.AuthService
}

func NewWebService(c config.Config, e *ent.Client, hc *http.Client,
	cache cache.Cache, prikey *rsa.PrivateKey, authService *auth.AuthService) *WebService {
	return &WebService{
		config:      c,
		client:      e,
		httpClient:  hc,
		cache:       cache,
		prikey:      prikey,
		authService: authService,
	}
}

func (w *WebService) validatePass(ctx context.Context, u *ent.User, password string) error {
	if !utils.Argon2Compare(password, u.Password, u.Salt) {
		return fmt.Errorf("validatePass: %w", ErrPassWord)
	}
	return nil
}
