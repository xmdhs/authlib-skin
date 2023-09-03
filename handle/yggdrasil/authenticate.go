package yggdrasil

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	yggdrasilS "github.com/xmdhs/authlib-skin/service/yggdrasil"
	"github.com/xmdhs/authlib-skin/utils"
)

func (y *Yggdrasil) Authenticate() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cxt := r.Context()
		a, err := utils.DeCodeBody[yggdrasil.Authenticate](r.Body, y.validate)
		if err != nil {
			y.logger.DebugContext(cxt, err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 400)
			return
		}
		t, err := y.yggdrasilService.Authenticate(cxt, a)
		if err != nil {
			if errors.Is(err, yggdrasilS.ErrPassWord) || errors.Is(err, yggdrasilS.ErrRate) {
				y.logger.DebugContext(cxt, err.Error())
				handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: "Invalid credentials. Invalid username or password.", Error: "ForbiddenOperationException"}, 403)
				return
			}
			y.logger.WarnContext(cxt, err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 500)
			return
		}
		b, _ := json.Marshal(t)
		w.Write(b)
	}
}
