package yggdrasil

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	sutils "github.com/xmdhs/authlib-skin/service/utils"
	yggdrasilS "github.com/xmdhs/authlib-skin/service/yggdrasil"
)

func (y *Yggdrasil) Authenticate() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func (y *Yggdrasil) Validate() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cxt := r.Context()
		a, has := getAnyModel[yggdrasil.ValidateToken](cxt, w, r.Body, y.validate, y.logger)
		if !has {
			return
		}
		err := y.yggdrasilService.ValidateToken(cxt, a)
		if err != nil {
			if errors.Is(err, sutils.ErrTokenInvalid) {
				y.logger.DebugContext(cxt, err.Error())
				handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: "Invalid token.", Error: "ForbiddenOperationException"}, 403)
				return
			}
			y.handleYgError(cxt, w, err)
			return
		}
		w.WriteHeader(204)
	}
}

func (y *Yggdrasil) Signout() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func (y *Yggdrasil) Invalidate() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(204)
		cxt := r.Context()
		a, has := getAnyModel[yggdrasil.ValidateToken](cxt, w, r.Body, y.validate, y.logger)
		if !has {
			return
		}
		err := y.yggdrasilService.Invalidate(cxt, a.AccessToken)
		if err != nil {
			if errors.Is(err, sutils.ErrTokenInvalid) {
				y.logger.DebugContext(cxt, err.Error())
				return
			}
			y.logger.WarnContext(cxt, err.Error())
		}
	}
}

func (y *Yggdrasil) Refresh() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cxt := r.Context()
		a, has := getAnyModel[yggdrasil.RefreshToken](cxt, w, r.Body, y.validate, y.logger)
		if !has {
			return
		}
		t, err := y.yggdrasilService.Refresh(cxt, a)
		if err != nil {
			if errors.Is(err, sutils.ErrTokenInvalid) {
				y.logger.DebugContext(cxt, err.Error())
				handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: "Invalid token.", Error: "ForbiddenOperationException"}, 403)
				return
			}
			y.handleYgError(cxt, w, err)
			return
		}
		b, _ := json.Marshal(t)
		w.Write(b)
	}
}

func (y *Yggdrasil) GetProfile() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()
		uuid := p.ByName("uuid")
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

func (y *Yggdrasil) BatchProfile() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func (y *Yggdrasil) PlayerCertificates() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()
		token := y.getTokenbyAuthorization(ctx, w, r)
		if token == "" {
			return
		}
		c, err := y.yggdrasilService.PlayerCertificates(ctx, token)
		if err != nil {
			if errors.Is(err, sutils.ErrTokenInvalid) {
				y.logger.DebugContext(ctx, err.Error())
				handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: "Invalid token.", Error: "ForbiddenOperationException"}, 403)
				return
			}
			y.handleYgError(ctx, w, err)
			return
		}
		w.Write(lo.Must(json.Marshal(c)))
	}
}
