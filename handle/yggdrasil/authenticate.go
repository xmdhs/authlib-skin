package yggdrasil

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
)

func Authenticate(l *slog.Logger, v *validator.Validate) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cxt := r.Context()
		jr := json.NewDecoder(r.Body)
		var a yggdrasil.Authenticate
		err := jr.Decode(&a)
		if err != nil {
			l.Info(err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 400)
			return
		}

		err = v.Struct(a)
		if err != nil {
			l.Info(err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 400)
			return
		}

	}
}
