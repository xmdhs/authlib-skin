package handle

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

func (h *Handel) handleErrorService(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrExistUser) {
		h.handleError(ctx, w, err.Error(), model.ErrExistUser, 400, slog.LevelDebug)
		return
	}
	if errors.Is(err, service.ErrExitsName) {
		h.handleError(ctx, w, err.Error(), model.ErrExitsName, 400, slog.LevelDebug)
		return
	}
	if errors.Is(err, service.ErrRegLimit) {
		h.handleError(ctx, w, err.Error(), model.ErrRegLimit, 400, slog.LevelDebug)
		return
	}
	if errors.Is(err, captcha.ErrCaptcha) {
		h.handleError(ctx, w, err.Error(), model.ErrCaptcha, 400, slog.LevelDebug)
		return
	}
	if errors.Is(err, service.ErrPassWord) {
		h.handleError(ctx, w, err.Error(), model.ErrPassWord, 401, slog.LevelDebug)
		return
	}
	if errors.Is(err, auth.ErrUserDisable) {
		h.handleError(ctx, w, err.Error(), model.ErrUserDisable, 401, slog.LevelDebug)
		return
	}
	if errors.Is(err, service.ErrNotAdmin) {
		h.handleError(ctx, w, err.Error(), model.ErrNotAdmin, 401, slog.LevelDebug)
		return
	}
	if errors.Is(err, auth.ErrTokenInvalid) {
		h.handleError(ctx, w, err.Error(), model.ErrAuth, 401, slog.LevelDebug)
		return
	}

	h.handleError(ctx, w, err.Error(), model.ErrService, 500, slog.LevelWarn)
}

func (h *Handel) handleError(ctx context.Context, w http.ResponseWriter, msg string, code model.APIStatus, httpcode int, level slog.Level) {
	h.logger.Log(ctx, level, msg)
	w.WriteHeader(httpcode)
	b, err := json.Marshal(model.API[any]{Code: code, Msg: msg, Data: nil})
	if err != nil {
		panic(err)
	}
	w.Write(b)
}
