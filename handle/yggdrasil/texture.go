package yggdrasil

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"

	yggdrasilS "github.com/xmdhs/authlib-skin/service/yggdrasil"
)

func (y *Yggdrasil) getTokenbyAuthorization(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		y.logger.DebugContext(ctx, "缺少 Authorization")
		w.WriteHeader(401)
		return ""
	}
	al := strings.Split(auth, " ")
	if len(al) != 2 || al[0] != "Bearer" {
		y.logger.DebugContext(ctx, "Authorization 格式错误")
		w.WriteHeader(401)
		return ""
	}
	return al[1]
}

func (y *Yggdrasil) validTextureType(ctx context.Context, w http.ResponseWriter, textureType string) bool {
	switch textureType {
	case "skin":
	case "cape":
	default:
		y.logger.DebugContext(ctx, "错误的材质类型")
		handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: "错误的材质类型"}, 400)
		return false
	}
	return true
}

func getUUIDbyParams(ctx context.Context, l *slog.Logger, w http.ResponseWriter) (string, string, bool) {
	uuid := chi.URLParamFromCtx(ctx, "uuid")
	textureType := chi.URLParamFromCtx(ctx, "textureType")
	if uuid == "" {
		l.DebugContext(ctx, "路径中缺少参数 uuid")
		handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: "路径中缺少参数 uuid"}, 400)
		return "", "", false
	}
	if textureType != "skin" && textureType != "cape" {
		l.DebugContext(ctx, "上传类型错误")
		handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: "上传类型错误"}, 400)
		return "", "", false

	}
	return uuid, textureType, true
}

func (y *Yggdrasil) DelTexture() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		uuid, textureType, ok := getUUIDbyParams(ctx, y.logger, w)
		if !ok {
			return
		}
		t := ctx.Value(tokenKey).(*model.TokenClaims)

		if uuid != t.Subject {
			y.logger.DebugContext(ctx, "uuid 不相同")
			w.WriteHeader(401)
			return
		}

		err := y.yggdrasilService.DelTexture(ctx, t, textureType)
		if err != nil {
			if errors.Is(err, yggdrasilS.ErrUUIDNotEq) {
				y.logger.DebugContext(ctx, err.Error())
				w.WriteHeader(401)
				return
			}
			y.handleYgError(ctx, w, err)
			return
		}
		w.WriteHeader(204)
	}
}
