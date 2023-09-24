package service

import (
	"net/http"

	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/ent"
)

type WebService struct {
	config     config.Config
	client     *ent.Client
	httpClient *http.Client
}

func NewWebService(c config.Config, e *ent.Client, hc *http.Client) *WebService {
	return &WebService{
		config:     c,
		client:     e,
		httpClient: hc,
	}
}
