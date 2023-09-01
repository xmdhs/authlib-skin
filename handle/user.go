package handle

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/db/mysql"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/utils"
)

func Reg(l *slog.Logger, q mysql.Querier, v *validator.Validate) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()

		u, err := utils.DeCodeBody[model.User](r.Body, v)
		if err != nil {
			l.InfoContext(ctx, err.Error())
		}
		_ = u

	}
}
