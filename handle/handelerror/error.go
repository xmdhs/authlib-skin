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
	"github.com/xmdhs/authlib-skin/service/email"
)

type HandleError struct {
	logger *slog.Logger
}

func NewHandleError(logger *slog.Logger) *HandleError {
	return &HandleError{
		logger: logger,
	}
}

type errorHandler struct {
	ErrorType  error
	ModelError model.APIStatus
	StatusCode int
	LogLevel   slog.Level
}

var errorHandlers = []errorHandler{
	{service.ErrExistUser, model.ErrExistUser, 400, slog.LevelDebug},
	{service.ErrExitsName, model.ErrExitsName, 400, slog.LevelDebug},
	{service.ErrRegLimit, model.ErrRegLimit, 400, slog.LevelInfo},
	{captcha.ErrCaptcha, model.ErrCaptcha, 400, slog.LevelDebug},
	{service.ErrPassWord, model.ErrPassWord, 401, slog.LevelInfo},
	{auth.ErrUserDisable, model.ErrUserDisable, 401, slog.LevelDebug},
	{service.ErrNotAdmin, model.ErrNotAdmin, 401, slog.LevelDebug},
	{auth.ErrTokenInvalid, model.ErrAuth, 401, slog.LevelDebug},
	{email.ErrTokenInvalid, model.ErrAuth, 401, slog.LevelDebug},
	{email.ErrSendLimit, model.ErrEmailSend, 403, slog.LevelDebug},
	{service.ErrUsername, model.ErrPassWord, 401, slog.LevelInfo},
}

func (h *HandleError) Service(ctx context.Context, w http.ResponseWriter, err error) {
	for _, errorHandler := range errorHandlers {
		if errors.Is(err, errorHandler.ErrorType) {
			h.Error(ctx, w, err.Error(), errorHandler.ModelError, errorHandler.StatusCode, errorHandler.LogLevel)
			return
		}
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
