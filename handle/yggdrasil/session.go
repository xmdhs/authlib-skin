package yggdrasil

import (
	"encoding/json"
	"net/http"

	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	"github.com/xmdhs/authlib-skin/utils"
)

func (y *Yggdrasil) SessionJoin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		a, has := getAnyModel[yggdrasil.Session](ctx, w, r.Body, y.validate, y.logger)
		if !has {
			return
		}
		t := ctx.Value(tokenKey).(*model.TokenClaims)

		ip, err := utils.GetIP(r)
		if err != nil {
			y.handleYgError(ctx, w, err)
			return
		}
		err = y.yggdrasilService.SessionJoin(ctx, a, t, ip)
		if err != nil {
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
