package yggdrasil

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	sutils "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

func (y *Yggdrasil) Validate() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cxt := r.Context()
		a, err := utils.DeCodeBody[yggdrasil.ValidateToken](r.Body, y.validate)
		if err != nil {
			y.logger.DebugContext(cxt, err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 400)
			return
		}
		err = y.yggdrasilService.ValidateToken(cxt, a)
		if err != nil {
			if errors.Is(err, sutils.ErrTokenInvalid) {
				y.logger.DebugContext(cxt, err.Error())
				handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: "Invalid token.", Error: "ForbiddenOperationException"}, 403)
				return
			}
			y.logger.WarnContext(cxt, err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 500)
			return
		}
		w.WriteHeader(204)
	}
}
