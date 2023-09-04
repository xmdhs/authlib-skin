package yggdrasil

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	yggdrasilS "github.com/xmdhs/authlib-skin/service/yggdrasil"
	"github.com/xmdhs/authlib-skin/utils"
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

func getAnyModel[K any](ctx context.Context, w http.ResponseWriter, r io.Reader, validate *validator.Validate, slog *slog.Logger) (K, bool) {
	a, err := utils.DeCodeBody[K](r, validate)
	if err != nil {
		slog.DebugContext(ctx, err.Error())
		handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: err.Error()}, 400)
		return a, false
	}
	return a, true
}
