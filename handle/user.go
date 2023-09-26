package handle

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service"
	utilsService "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

func (h *Handel) Reg() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()

		ip, err := utils.GetIP(r, h.config.RaelIP)
		if err != nil {
			h.logger.InfoContext(ctx, err.Error())
			handleError(ctx, w, err.Error(), model.ErrInput, 400)
			return
		}

		u, err := utils.DeCodeBody[model.User](r.Body, h.validate)
		if err != nil {
			h.logger.InfoContext(ctx, err.Error())
			handleError(ctx, w, err.Error(), model.ErrInput, 400)
			return
		}
		rip, err := getPrefix(ip)
		if err != nil {
			h.logger.WarnContext(ctx, err.Error())
			handleError(ctx, w, err.Error(), model.ErrUnknown, 500)
			return
		}
		err = h.webService.Reg(ctx, u, rip, ip)
		if err != nil {
			if errors.Is(err, service.ErrExistUser) {
				h.logger.DebugContext(ctx, err.Error())
				handleError(ctx, w, err.Error(), model.ErrExistUser, 400)
				return
			}
			if errors.Is(err, service.ErrRegLimit) {
				h.logger.DebugContext(ctx, err.Error())
				handleError(ctx, w, err.Error(), model.ErrRegLimit, 400)
				return
			}
			h.logger.WarnContext(ctx, err.Error())
			handleError(ctx, w, err.Error(), model.ErrService, 500)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}

func (h *Handel) UserInfo() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()
		token := h.getTokenbyAuthorization(ctx, w, r)
		if token == "" {
			return
		}

		u, err := h.webService.Info(ctx, token)
		if err != nil {
			if errors.Is(err, utilsService.ErrTokenInvalid) {
				h.logger.DebugContext(ctx, "token 无效")
				handleError(ctx, w, "token 无效", model.ErrAuth, 401)
				return
			}
			h.logger.InfoContext(ctx, err.Error())
			handleError(ctx, w, err.Error(), model.ErrUnknown, 500)
			return
		}
		encodeJson(w, model.API[model.UserInfo]{
			Code: 0,
			Data: u,
		})
	}
}
