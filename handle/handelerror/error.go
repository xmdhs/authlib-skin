package handelerror

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service"
	"github.com/xmdhs/authlib-skin/service/auth"
	"github.com/xmdhs/authlib-skin/service/captcha"
)

type HandleError struct {
	logger *slog.Logger
}

func NewHandleError(logger *slog.Logger) *HandleError {
	return &HandleError{
		logger: logger,
	}
}

func (h *HandleError) Service(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrExistUser) {
		h.Error(ctx, w, err.Error(), model.ErrExistUser, 400, slog.LevelDebug)
		return
	}
	if errors.Is(err, service.ErrExitsName) {
		h.Error(ctx, w, err.Error(), model.ErrExitsName, 400, slog.LevelDebug)
		return
	}
	if errors.Is(err, service.ErrRegLimit) {
		h.Error(ctx, w, err.Error(), model.ErrRegLimit, 400, slog.LevelDebug)
		return
	}
	if errors.Is(err, captcha.ErrCaptcha) {
		h.Error(ctx, w, err.Error(), model.ErrCaptcha, 400, slog.LevelDebug)
		return
	}
	if errors.Is(err, service.ErrPassWord) {
		h.Error(ctx, w, err.Error(), model.ErrPassWord, 401, slog.LevelDebug)
		return
	}
	if errors.Is(err, auth.ErrUserDisable) {
		h.Error(ctx, w, err.Error(), model.ErrUserDisable, 401, slog.LevelDebug)
		return
	}
	if errors.Is(err, service.ErrNotAdmin) {
		h.Error(ctx, w, err.Error(), model.ErrNotAdmin, 401, slog.LevelDebug)
		return
	}
	if errors.Is(err, auth.ErrTokenInvalid) {
		h.Error(ctx, w, err.Error(), model.ErrAuth, 401, slog.LevelDebug)
		return
	}

	h.Error(ctx, w, err.Error(), model.ErrService, 500, slog.LevelWarn)
}

func (h *HandleError) Error(ctx context.Context, w http.ResponseWriter, msg string, code model.APIStatus, httpcode int, level slog.Level) {
	h.logger.Log(ctx, level, msg)
	w.WriteHeader(httpcode)
	b, err := json.Marshal(model.API[any]{Code: code, Msg: msg, Data: nil})
	if err != nil {
		panic(err)
	}
	w.Write(b)
}
