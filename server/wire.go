//go:build wireinject

package server

import (
	"context"
	"net/http"

	"github.com/google/wire"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/handle"
	"github.com/xmdhs/authlib-skin/handle/handelerror"
	"github.com/xmdhs/authlib-skin/handle/yggdrasil"
	"github.com/xmdhs/authlib-skin/server/route"
	"github.com/xmdhs/authlib-skin/service"
	"github.com/xmdhs/authlib-skin/service/auth"
	"github.com/xmdhs/authlib-skin/service/captcha"
	"github.com/xmdhs/authlib-skin/service/email"
	yggdrasilS "github.com/xmdhs/authlib-skin/service/yggdrasil"
)

var serviceSet = wire.NewSet(service.Service, yggdrasilS.NewYggdrasil, email.NewEmail, auth.NewAuthService,
	captcha.NewCaptchaService,
)

var handleSet = wire.NewSet(handelerror.NewHandleError, handle.HandelSet, yggdrasil.NewYggdrasil)

func InitializeRoute(ctx context.Context, c config.Config) (*http.Server, func(), error) {
	panic(wire.Build(Set, route.NewRoute, NewSlog,
		NewServer, handleSet, serviceSet,
	))
}
