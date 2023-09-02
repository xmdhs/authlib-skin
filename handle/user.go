package handle

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/bwmarrin/snowflake"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/mysql"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service"
	"github.com/xmdhs/authlib-skin/utils"
)

func Reg(l *slog.Logger, q mysql.QuerierWithTx, v *validator.Validate, snow *snowflake.Node, c config.Config) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()

		u, err := utils.DeCodeBody[model.User](r.Body, v)
		if err != nil {
			l.InfoContext(ctx, err.Error())
			handleError(ctx, w, err.Error(), model.ErrInput, 400)
			return
		}
		err = service.Reg(ctx, u, q, snow, c)
		if err != nil {
			if errors.Is(err, service.ErrExistUser) {
				l.DebugContext(ctx, err.Error())
				handleError(ctx, w, err.Error(), model.ErrExistUser, 400)
				return
			}
			l.WarnContext(ctx, err.Error())
			handleError(ctx, w, err.Error(), model.ErrService, 500)
			return
		}
	}
}
