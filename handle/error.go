package handle

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xmdhs/authlib-skin/model"
)

func handleError(ctx context.Context, w http.ResponseWriter, msg string, code int, httpcode int) {
	w.WriteHeader(httpcode)
	b, err := json.Marshal(model.API[any]{Code: code, Msg: msg, Data: nil})
	if err != nil {
		panic(err)
	}
	w.Write(b)
}
