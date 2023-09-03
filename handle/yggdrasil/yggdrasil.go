package yggdrasil

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	yggdrasilS "github.com/xmdhs/authlib-skin/service/yggdrasil"
)

type Yggdrasil struct {
	logger           *slog.Logger
	validate         *validator.Validate
	yggdrasilService *yggdrasilS.Yggdrasil
}

func NewYggdrasil(logger *slog.Logger, validate *validator.Validate, yggdrasilService *yggdrasilS.Yggdrasil) *Yggdrasil {
	return &Yggdrasil{
		logger:           logger,
		validate:         validate,
		yggdrasilService: yggdrasilService,
	}
}
