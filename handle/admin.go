package handle

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service"
	"github.com/xmdhs/authlib-skin/service/utils"
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
			if errors.Is(err, utils.ErrTokenInvalid) {
				h.handleError(ctx, w, err.Error(), model.ErrAuth, 401, slog.LevelDebug)
				return
			}
			h.handleError(ctx, w, err.Error(), model.ErrService, 500, slog.LevelWarn)
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
			if errors.Is(err, service.ErrNotAdmin) {
				h.handleError(ctx, w, err.Error(), model.ErrNotAdmin, 401, slog.LevelDebug)
				return
			}
			h.handleError(ctx, w, err.Error(), model.ErrService, 500, slog.LevelWarn)
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

		ul, uc, err := h.webService.ListUser(ctx, pagei)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrService, 500, slog.LevelWarn)
			return
		}
		encodeJson(w, model.API[model.List[model.UserList]]{Data: model.List[model.UserList]{List: ul, Total: uc}})
	}
}
