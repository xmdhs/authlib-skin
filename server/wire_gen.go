// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"context"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/handle"
	yggdrasil2 "github.com/xmdhs/authlib-skin/handle/yggdrasil"
	"github.com/xmdhs/authlib-skin/server/route"
	"github.com/xmdhs/authlib-skin/service"
	"github.com/xmdhs/authlib-skin/service/yggdrasil"
	"net/http"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

// Injectors from wire.go:

func InitializeRoute(ctx context.Context, c config.Config) (*http.Server, func(), error) {
	handler := ProvideSlog(c)
	logger := NewSlog(handler)
	validate := ProvideValidate()
	db, cleanup, err := ProvideDB(c)
	if err != nil {
		return nil, nil, err
	}
	client, cleanup2, err := ProvideEnt(ctx, db, c, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	cache := ProvideCache(c)
	privateKey, err := ProvidePriKey(c)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	yggdrasilYggdrasil := yggdrasil.NewYggdrasil(client, cache, c, privateKey)
	pubRsaKey, err := ProvidePubKey(privateKey)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	yggdrasil3 := yggdrasil2.NewYggdrasil(logger, validate, yggdrasilYggdrasil, c, pubRsaKey)
	httpClient := ProvideHttpClient()
	webService := service.NewWebService(c, client, httpClient, cache, privateKey)
	handel := handle.NewHandel(webService, validate, c, logger)
	router, err := route.NewRoute(yggdrasil3, handel)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	server, cleanup3 := NewServer(c, logger, router)
	return server, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
