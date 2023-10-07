package yggdrasil

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	sutils "github.com/xmdhs/authlib-skin/service/utils"
	"github.com/xmdhs/authlib-skin/utils"
)

func (y *Yggdrasil) SessionJoin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		a, has := getAnyModel[yggdrasil.Session](ctx, w, r.Body, y.validate, y.logger)
		if !has {
			return
		}
		ip, err := utils.GetIP(r)
		if err != nil {
			y.handleYgError(ctx, w, err)
			return
		}
		err = y.yggdrasilService.SessionJoin(ctx, a, ip)
		if err != nil {
			if errors.Is(err, sutils.ErrTokenInvalid) {
				y.logger.DebugContext(ctx, err.Error())
				handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: "Invalid token.", Error: "ForbiddenOperationException"}, 403)
				return
			}
			y.handleYgError(ctx, w, err)
			return
		}
		w.WriteHeader(204)
	}
}

func (y *Yggdrasil) HasJoined() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		name := r.FormValue("username")
		serverId := r.FormValue("serverId")
		ip := r.FormValue("ip")
		if name == "" || serverId == "" {
			y.logger.DebugContext(ctx, "name 或 serverID 为空")
			w.WriteHeader(204)
			return
		}
		u, err := y.yggdrasilService.HasJoined(ctx, name, serverId, ip, r.Host)
		if err != nil {
			y.logger.WarnContext(ctx, err.Error())
			w.WriteHeader(204)
		}
		w.Write(lo.Must(json.Marshal(u)))
	}
}
