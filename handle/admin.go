package handle

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/xmdhs/authlib-skin/handle/handelerror"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service"

	U "github.com/xmdhs/authlib-skin/utils"
)

type AdminHandel struct {
	handleError  *handelerror.HandleError
	adminService *service.AdminService
	validate     *validator.Validate
}

func NewAdminHandel(handleError *handelerror.HandleError, adminService *service.AdminService, validate *validator.Validate) *AdminHandel {
	return &AdminHandel{
		handleError:  handleError,
		adminService: adminService,
		validate:     validate,
	}
}

type tokenValue string

const tokenKey = tokenValue("token")

func (h *AdminHandel) NeedAuth(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := h.getTokenbyAuthorization(ctx, w, r)
		if token == "" {
			return
		}
		t, err := h.adminService.Auth(ctx, token)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		r = r.WithContext(context.WithValue(ctx, tokenKey, t))
		handle.ServeHTTP(w, r)
	})
}

func (h *AdminHandel) NeedAdmin(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		err := h.adminService.IsAdmin(ctx, t)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		handle.ServeHTTP(w, r)
	})
}

func (h *AdminHandel) ListUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		page := r.FormValue("page")
		pagei := 1
		if page != "" {
			p, err := strconv.Atoi(page)
			if err != nil {
				h.handleError.Error(ctx, w, "page 必须为数字", model.ErrInput, 400, slog.LevelDebug)
				return
			}
			if p == 0 {
				p = 1
			}
			pagei = p
		}
		email := r.FormValue("email")
		name := r.FormValue("name")

		ul, uc, err := h.adminService.ListUser(ctx, pagei, email, name)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson(w, model.API[model.List[model.UserList]]{Data: model.List[model.UserList]{List: ul, Total: uc}})
	}
}

func (h *AdminHandel) EditUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		uid := chi.URLParamFromCtx(ctx, "uid")
		if uid == "" {
			h.handleError.Error(ctx, w, "uid 为空", model.ErrInput, 400, slog.LevelDebug)
			return
		}
		uidi, err := strconv.Atoi(uid)
		if err != nil {
			h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		a, err := U.DeCodeBody[model.EditUser](r.Body, h.validate)
		if err != nil {
			h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		err = h.adminService.EditUser(ctx, a, uidi)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson[any](w, model.API[any]{
			Code: 0,
		})
	}
}

func (h *AdminHandel) getTokenbyAuthorization(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		h.handleError.Error(ctx, w, "缺少 Authorization", model.ErrAuth, 401, slog.LevelDebug)
		return ""
	}
	al := strings.Split(auth, " ")
	if len(al) != 2 || al[0] != "Bearer" {
		h.handleError.Error(ctx, w, "Authorization 格式错误", model.ErrAuth, 401, slog.LevelDebug)
		return ""
	}
	return al[1]
}
