package handle

import (
	"log/slog"
	"net/http"

	"github.com/xmdhs/authlib-skin/model"
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

		u, err := utils.DeCodeBody[model.UserReg](r.Body, h.validate)
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
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}

func (h *Handel) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ip, err := utils.GetIP(r)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		l, err := utils.DeCodeBody[model.Login](r.Body, h.validate)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		lr, err := h.webService.Login(ctx, l, ip)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson(w, model.API[model.LoginRep]{
			Code: 0,
			Data: lr,
		})
	}
}

func (h *Handel) UserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		u, err := h.webService.Info(ctx, t)
		if err != nil {
			h.handleErrorService(ctx, w, err)
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
			h.handleErrorService(ctx, w, err)
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
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}
