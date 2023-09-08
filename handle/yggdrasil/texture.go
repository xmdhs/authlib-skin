package yggdrasil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
	"github.com/xmdhs/authlib-skin/service/utils"
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

func (y *Yggdrasil) PutTexture() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()
		uuid := p.ByName("uuid")
		textureType := p.ByName("textureType")
		if uuid == "" || textureType == "" {
			y.logger.DebugContext(ctx, "路径中缺少参数 uuid / textureType")
			handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: "路径中缺少参数 uuid / textureType"}, 400)
			return
		}
		token := y.getTokenbyAuthorization(ctx, w, r)

		model := r.FormValue("model")

		skin, err := func() ([]byte, error) {
			f, _, err := r.FormFile("file")
			if err != nil {
				return nil, err
			}
			b, err := io.ReadAll(io.LimitReader(f, 5*1000*1000))
			if err != nil {
				return nil, err
			}
			pc, err := png.DecodeConfig(bytes.NewReader(b))
			if err != nil {
				return nil, err
			}
			if pc.Height > 200 || pc.Width > 200 {
				return nil, fmt.Errorf("材质大小超过限制")
			}
			p, err := png.Decode(bytes.NewReader(b))
			if err != nil {
				return nil, err
			}
			bw := bytes.NewBuffer(nil)
			err = png.Encode(bw, p)
			return bw.Bytes(), err
		}()
		if err != nil {
			y.logger.DebugContext(ctx, err.Error())
			handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: err.Error()}, 400)
			return
		}
		if !y.validTextureType(ctx, w, textureType) {
			return
		}

		switch model {
		case "slim":
		case "":
		default:
			y.logger.DebugContext(ctx, "错误的皮肤的材质模型")
			handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: "错误的皮肤的材质模型"}, 400)
			return
		}

		err = y.yggdrasilService.PutTexture(ctx, token, skin, model, uuid, textureType)
		if err != nil {
			if errors.Is(err, utils.ErrTokenInvalid) {
				y.logger.DebugContext(ctx, err.Error())
				w.WriteHeader(401)
				return
			}

			y.logger.WarnContext(ctx, err.Error())
			handleYgError(ctx, w, yggdrasil.Error{ErrorMessage: err.Error()}, 500)
			return
		}
		w.WriteHeader(204)
	}
}
