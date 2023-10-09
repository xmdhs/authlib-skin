package handle

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/xmdhs/authlib-skin/model"

	U "github.com/xmdhs/authlib-skin/utils"
)

type tokenValue string

const tokenKey = tokenValue("token")

func (h *Handel) NeedAuth(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := h.getTokenbyAuthorization(ctx, w, r)
		if token == "" {
			return
		}
		t, err := h.webService.Auth(ctx, token)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		r = r.WithContext(context.WithValue(ctx, tokenKey, t))
		handle.ServeHTTP(w, r)
	})
}

func (h *Handel) NeedAdmin(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		err := h.webService.IsAdmin(ctx, t)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		handle.ServeHTTP(w, r)
	})
}

func (h *Handel) ListUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		page := r.FormValue("page")
		pagei := 1
		if page != "" {
			p, err := strconv.Atoi(page)
			if err != nil {
				h.handleError(ctx, w, "page 必须为数字", model.ErrInput, 400, slog.LevelDebug)
				return
			}
			if p == 0 {
				p = 1
			}
			pagei = p
		}
		email := r.FormValue("email")
		name := r.FormValue("name")

		ul, uc, err := h.webService.ListUser(ctx, pagei, email, name)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson(w, model.API[model.List[model.UserList]]{Data: model.List[model.UserList]{List: ul, Total: uc}})
	}
}

func (h *Handel) EditUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		uid := chi.URLParamFromCtx(ctx, "uid")
		if uid == "" {
			h.handleError(ctx, w, "uid 为空", model.ErrInput, 400, slog.LevelDebug)
			return
		}
		uidi, err := strconv.Atoi(uid)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		a, err := U.DeCodeBody[model.EditUser](r.Body, h.validate)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		err = h.webService.EditUser(ctx, a, uidi)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson[any](w, model.API[any]{
			Code: 0,
		})
	}
}
