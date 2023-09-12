package yggdrasil

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xmdhs/authlib-skin/model/yggdrasil"
)

func handleYgError(ctx context.Context, w http.ResponseWriter, e yggdrasil.Error, httpcode int) {
	w.WriteHeader(httpcode)
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func (y *Yggdrasil) handleYgError(ctx context.Context, w http.ResponseWriter, err error) {
	y.logger.WarnContext(ctx, err.Error())
	handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: err.Error()}, 500)
	return
}
