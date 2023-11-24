package service

import "github.com/google/wire"

var Service = wire.NewSet(NewUserSerice, NewTextureService, NewAdminService, NewWebService)
