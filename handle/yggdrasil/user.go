package yggdrasil

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
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
			y.logger.WarnContext(cxt, err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 500)
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
			y.logger.WarnContext(cxt, err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 500)
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
			y.logger.WarnContext(cxt, err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 500)
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
			y.logger.WarnContext(cxt, err.Error())
			handleYgError(cxt, w, yggdrasil.Error{ErrorMessage: err.Error()}, 500)
			return
		}
		b, _ := json.Marshal(t)
		w.Write(b)
	}
}
