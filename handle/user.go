package handle

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service"
	"github.com/xmdhs/authlib-skin/utils"
)

func (h *Handel) Reg() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()

		u, err := utils.DeCodeBody[model.User](r.Body, h.validate)
		if err != nil {
			h.logger.InfoContext(ctx, err.Error())
			handleError(ctx, w, err.Error(), model.ErrInput, 400)
			return
		}
		err = h.webService.Reg(ctx, u)
		if err != nil {
			if errors.Is(err, service.ErrExistUser) {
				h.logger.DebugContext(ctx, err.Error())
				handleError(ctx, w, err.Error(), model.ErrExistUser, 400)
				return
			}
			h.logger.WarnContext(ctx, err.Error())
			handleError(ctx, w, err.Error(), model.ErrService, 500)
			return
		}
	}
}
