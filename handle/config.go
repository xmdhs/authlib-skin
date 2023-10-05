package handle

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model"
)

func (h *Handel) GetConfig() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()
		c := h.webService.GetConfig(ctx)
		m := model.API[model.Config]{
			Code: 0,
			Data: c,
		}
		json.NewEncoder(w).Encode(m)
	}
}
