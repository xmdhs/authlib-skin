package service

import (
	"crypto/rsa"
	"net/http"

	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
)

type WebService struct {
	config     config.Config
	client     *ent.Client
	httpClient *http.Client
	cache      cache.Cache
	prikey     *rsa.PrivateKey
}

func NewWebService(c config.Config, e *ent.Client, hc *http.Client, cache cache.Cache, prikey *rsa.PrivateKey) *WebService {
	return &WebService{
		config:     c,
		client:     e,
		httpClient: hc,
		cache:      cache,
		prikey:     prikey,
	}
}
