package yggdrasil

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	"github.com/xmdhs/authlib-skin/utils"
)

func Authenticate(l *slog.Logger, v *validator.Validate) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cxt := r.Context()
		a, err := utils.DeCodeBody[yggdrasil.Authenticate](r.Body, v)
		if err != nil {
			l.InfoContext(cxt, err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 400)
			return
		}

		_ = a
	}
}
