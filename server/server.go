package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewYggdrasil(r *httprouter.Router) error {

	r.POST("/api/authserver/authenticate", nil)
	return nil
}

func warpCtJSON(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		handle(w, r, p)
	}
}
