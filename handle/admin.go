package handle

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
)

func (h *Handel) NeedAdmin(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()
		token := h.getTokenbyAuthorization(ctx, w, r)
		if token == "" {
			return
		}
		err := h.webService.IsAdmin(ctx, token)
		if err != nil {
			if errors.Is(err, utilsService.ErrTokenInvalid) {
				h.handleError(ctx, w, "token 无效", model.ErrAuth, 401, slog.LevelDebug)
				return
			}
			h.handleError(ctx, w, err.Error(), model.ErrService, 500, slog.LevelWarn)
			return
		}
		handle(w, r, p)
	}
}

func (h *Handel) ListUser() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
