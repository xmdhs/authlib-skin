package yggdrasil

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (y *Yggdrasil) Validate() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	}
}
