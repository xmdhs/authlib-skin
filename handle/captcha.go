package handle

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model"
)

func (h *Handel) GetCaptcha() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()
		c := h.webService.GetCaptcha(ctx)
		m := model.API[model.Captcha]{
			Code: 0,
			Data: c,
		}
		json.NewEncoder(w).Encode(m)
	}
}
