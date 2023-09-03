package route

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/snowflake"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/handle"
)

func NewRoute(l *slog.Logger, client *ent.Client, v *validator.Validate, snow *snowflake.Node, c config.Config) (*httprouter.Router, error) {
	r := httprouter.New()
	err := newYggdrasil(r)
	if err != nil {
		return nil, fmt.Errorf("NewRoute: %w", err)
	}
	err = newSkinApi(r, l, client, v, snow, c)
	if err != nil {
		return nil, fmt.Errorf("NewRoute: %w", err)
	}
	return r, nil
}

func newYggdrasil(r *httprouter.Router) error {
	r.POST("/api/authserver/authenticate", nil)
	return nil
}

func newSkinApi(r *httprouter.Router, l *slog.Logger, client *ent.Client, v *validator.Validate, snow *snowflake.Node, c config.Config) error {
	r.PUT("/api/v1/user/reg", handle.Reg(l, client, v, snow, c))
	return nil
}
