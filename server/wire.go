//go:build wireinject

package server

import (
	"context"
	"net/http"

	"github.com/google/wire"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/server/route"
)

func InitializeRoute(ctx context.Context, c config.Config) (*http.Server, func(), error) {
	panic(wire.Build(Set, route.NewRoute, NewSlog, NewServer))
}
