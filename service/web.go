package service

import (
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/ent"
)

type WebService struct {
	config config.Config
	client *ent.Client
}

func NewWebService(c config.Config, e *ent.Client) *WebService {
	return &WebService{
		config: c,
		client: e,
	}
}
