package yggdrasil

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	yggdrasilS "github.com/xmdhs/authlib-skin/service/yggdrasil"
)

func (y *Yggdrasil) Authenticate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cxt := r.Context()
		a, has := getAnyModel[yggdrasil.Authenticate](cxt, w, r.Body, y.validate, y.logger)
		if !has {
			return
		}
		t, err := y.yggdrasilService.Authenticate(cxt, a)
		if err != nil {
			if errors.Is(err, yggdrasilS.ErrPassWord) || errors.Is(err, yggdrasilS.ErrRate) {
				y.logger.DebugContext(cxt, err.Error())
				handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: "Invalid credentials. Invalid username or password.", Error: "ForbiddenOperationException"}, 403)
				return
			}
			y.handleYgError(cxt, w, err)
			return
		}
		b, _ := json.Marshal(t)
		w.Write(b)
	}
}

func (y *Yggdrasil) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}
}

func (y *Yggdrasil) Signout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cxt := r.Context()
		a, has := getAnyModel[yggdrasil.Pass](cxt, w, r.Body, y.validate, y.logger)
		if !has {
			return
		}
		err := y.yggdrasilService.SignOut(cxt, a)
		if err != nil {
			if errors.Is(err, yggdrasilS.ErrPassWord) || errors.Is(err, yggdrasilS.ErrRate) {
				y.logger.DebugContext(cxt, err.Error())
				handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: "Invalid credentials. Invalid username or password.", Error: "ForbiddenOperationException"}, 403)
				return
			}
			y.handleYgError(cxt, w, err)
			return
		}
		w.WriteHeader(204)
	}
}

func (y *Yggdrasil) Invalidate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
		cxt := r.Context()

		t := cxt.Value(tokenKey).(*model.TokenClaims)

		err := y.yggdrasilService.Invalidate(cxt, t)
		if err != nil {
			y.logger.WarnContext(cxt, err.Error())
		}
	}
}

func (y *Yggdrasil) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cxt := r.Context()
		token := cxt.Value(tokenKey).(*model.TokenClaims)
		t, err := y.yggdrasilService.Refresh(cxt, token)
		if err != nil {
			y.handleYgError(cxt, w, err)
			return
		}
		b, _ := json.Marshal(t)
		w.Write(b)
	}
}

func (y *Yggdrasil) GetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		uuid := chi.URLParamFromCtx(ctx, "uuid")

		unsigned := r.FormValue("unsigned")

		unsignedBool := true

		switch unsigned {
		case "true":
		case "false":
			unsignedBool = false
		case "":
		default:
			y.logger.DebugContext(ctx, "unsigned 参数类型错误")
			handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: "unsigned 参数类型错误"}, 400)
			return

		}
		u, err := y.yggdrasilService.GetProfile(ctx, uuid, unsignedBool, r.Host)
		if err != nil {
			if errors.Is(err, yggdrasilS.ErrNotUser) {
				y.logger.DebugContext(ctx, err.Error())
				w.WriteHeader(204)
				return
			}
			y.handleYgError(ctx, w, err)
			return
		}
		b, _ := json.Marshal(u)
		w.Write(b)
	}
}

func (y *Yggdrasil) BatchProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		a, has := getAnyModel[[]string](ctx, w, r.Body, nil, y.logger)
		if !has {
			return
		}
		if len(a) > 5 {
			y.logger.DebugContext(ctx, "最多同时查询五个")
			handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: "最多同时查询五个"}, 400)
			return
		}
		ul, err := y.yggdrasilService.BatchProfile(ctx, a)
		if err != nil {
			y.handleYgError(ctx, w, err)
			return
		}
		w.Write(lo.Must1(json.Marshal(ul)))
	}
}

func (y *Yggdrasil) PlayerCertificates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		c, err := y.yggdrasilService.PlayerCertificates(ctx, t)
		if err != nil {
			y.handleYgError(ctx, w, err)
			return
		}
		w.Write(lo.Must(json.Marshal(c)))
	}
}

func (y *Yggdrasil) PlayerAttributes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"privileges":{"onlineChat":{"enabled":true},"multiplayerServer":{"enabled":true},"multiplayerRealms":{"enabled":true},"telemetry":{"enabled":true},"optionalTelemetry":{"enabled":true}},"profanityFilterPreferences":{"profanityFilterOn":true},"banStatus":{"bannedScopes":{}}}`))
	}
}
