package handle

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/xmdhs/authlib-skin/model"
)

func (h *Handel) handleError(ctx context.Context, w http.ResponseWriter, msg string, code model.APIStatus, httpcode int, level slog.Level) {
	h.logger.Log(ctx, level, msg)
	w.WriteHeader(httpcode)
	b, err := json.Marshal(model.API[any]{Code: code, Msg: msg, Data: nil})
	if err != nil {
		panic(err)
	}
	w.Write(b)
}
