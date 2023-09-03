package handle

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/service"
)

type Handel struct {
	webService *service.WebService
	validate   *validator.Validate
	config     config.Config
	logger     *slog.Logger
}

func NewHandel(webService *service.WebService, validate *validator.Validate, config config.Config, logger *slog.Logger) *Handel {
	return &Handel{
		webService: webService,
		validate:   validate,
		config:     config,
		logger:     logger,
	}
}
