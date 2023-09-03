package yggdrasil

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	"github.com/xmdhs/authlib-skin/utils"
)

func (y *Yggdrasil) Authenticate() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cxt := r.Context()
		a, err := utils.DeCodeBody[yggdrasil.Authenticate](r.Body, y.validate)
		if err != nil {
			y.logger.InfoContext(cxt, err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 400)
			return
		}

		_ = a
	}
}
