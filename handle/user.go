package handle

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

func (h *Handel) Reg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ip, err := utils.GetIP(r)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		u, err := utils.DeCodeBody[model.User](r.Body, h.validate)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		rip, err := getPrefix(ip)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrUnknown, 500, slog.LevelWarn)
			return
		}
		err = h.webService.Reg(ctx, u, rip, ip)
		if err != nil {
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
			h.handleError(ctx, w, err.Error(), model.ErrService, 500, slog.LevelWarn)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}

func (h *Handel) UserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		u, err := h.webService.Info(ctx, t)
		if err != nil {
			if errors.Is(err, utilsService.ErrTokenInvalid) {
				h.handleError(ctx, w, "token 无效", model.ErrAuth, 401, slog.LevelDebug)
				return
			}
			h.handleError(ctx, w, err.Error(), model.ErrService, 500, slog.LevelWarn)
			return
		}
		encodeJson(w, model.API[model.UserInfo]{
			Code: 0,
			Data: u,
		})
	}
}

func (h *Handel) ChangePasswd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)

		c, err := utils.DeCodeBody[model.ChangePasswd](r.Body, h.validate)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		err = h.webService.ChangePasswd(ctx, c, t)
		if err != nil {
			if errors.Is(err, service.ErrPassWord) {
				h.handleError(ctx, w, err.Error(), model.ErrPassWord, 401, slog.LevelDebug)
				return
			}
			h.handleError(ctx, w, err.Error(), model.ErrService, 500, slog.LevelWarn)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})

	}
}

func (h *Handel) ChangeName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		c, err := utils.DeCodeBody[model.ChangeName](r.Body, h.validate)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		err = h.webService.ChangeName(ctx, c.Name, t)
		if err != nil {
			if errors.Is(err, service.ErrExitsName) {
				h.handleError(ctx, w, err.Error(), model.ErrExitsName, 400, slog.LevelDebug)
				return
			}
			h.handleError(ctx, w, err.Error(), model.ErrService, 500, slog.LevelWarn)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}
