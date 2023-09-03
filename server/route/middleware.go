package route

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func warpHJSON(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		handle(w, r, p)
	}
}
