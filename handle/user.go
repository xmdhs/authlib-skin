package handle

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/utils"
)

func (h *Handel) Reg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ip, err := utils.GetIP(r)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		u, err := utils.DeCodeBody[model.UserReg](r.Body, h.validate)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		rip, err := getPrefix(ip)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrUnknown, 500, slog.LevelWarn)
			return
		}
		lr, err := h.webService.Reg(ctx, u, rip, ip)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson(w, model.API[model.LoginRep]{
			Code: 0,
			Data: lr,
		})
	}
}

func (h *Handel) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ip, err := utils.GetIP(r)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		l, err := utils.DeCodeBody[model.Login](r.Body, h.validate)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		lr, err := h.webService.Login(ctx, l, ip)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson(w, model.API[model.LoginRep]{
			Code: 0,
			Data: lr,
		})
	}
}

func (h *Handel) UserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		u, err := h.webService.Info(ctx, t)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson(w, model.API[model.UserInfo]{
			Code: 0,
			Data: u,
		})
	}
}

func (h *Handel) ChangePasswd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)

		c, err := utils.DeCodeBody[model.ChangePasswd](r.Body, h.validate)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		err = h.webService.ChangePasswd(ctx, c, t)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})

	}
}

func (h *Handel) ChangeName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		c, err := utils.DeCodeBody[model.ChangeName](r.Body, h.validate)
		if err != nil {
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		err = h.webService.ChangeName(ctx, c.Name, t)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}

func (h *Handel) PutTexture() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		models := r.FormValue("model")

		textureType := chi.URLParamFromCtx(ctx, "textureType")
		if textureType != "skin" && textureType != "cape" {
			h.logger.DebugContext(ctx, "上传类型错误")
			h.handleError(ctx, w, "上传类型错误", model.ErrInput, 400, slog.LevelDebug)
		}

		skin, err := func() ([]byte, error) {
			f, _, err := r.FormFile("file")
			if err != nil {
				return nil, err
			}
			b, err := io.ReadAll(io.LimitReader(f, 50*1000))
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
			h.handleError(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		switch models {
		case "slim":
		case "":
		default:
			h.logger.DebugContext(ctx, "错误的皮肤的材质模型")
			h.handleError(ctx, w, "错误的皮肤的材质模型", model.ErrInput, 400, slog.LevelDebug)
			return
		}

		err = h.webService.PutTexture(ctx, t, skin, models, textureType)
		if err != nil {
			h.handleErrorService(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}
