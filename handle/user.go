package handle

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/handle/handelerror"
	"github.com/xmdhs/authlib-skin/model"
	"github.com/xmdhs/authlib-skin/service"
	"github.com/xmdhs/authlib-skin/utils"
)

type UserHandel struct {
	handleError    *handelerror.HandleError
	validate       *validator.Validate
	userService    *service.UserService
	logger         *slog.Logger
	textureService *service.TextureService
	config         config.Config

	emailReg func() (*regexp.Regexp, error)
}

func NewUserHandel(handleError *handelerror.HandleError, validate *validator.Validate,
	userService *service.UserService, logger *slog.Logger, textureService *service.TextureService, config config.Config) *UserHandel {
	emailReg := sync.OnceValues[*regexp.Regexp, error](func() (*regexp.Regexp, error) {
		return regexp.Compile(config.Email.EmailReg)
	})

	return &UserHandel{
		handleError:    handleError,
		validate:       validate,
		userService:    userService,
		logger:         logger,
		textureService: textureService,
		config:         config,

		emailReg: emailReg,
	}
}

func (h *UserHandel) Reg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ip, err := utils.GetIP(r)
		if err != nil {
			h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		u, err := utils.DeCodeBody[model.UserReg](r.Body, h.validate)
		if err != nil {
			h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		rip, err := getPrefix(ip)
		if err != nil {
			h.handleError.Error(ctx, w, err.Error(), model.ErrUnknown, 500, slog.LevelWarn)
			return
		}
		lr, err := h.userService.Reg(ctx, u, rip, ip)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson(w, model.API[model.LoginRep]{
			Code: 0,
			Data: lr,
		})
	}
}

func (h *UserHandel) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ip, err := utils.GetIP(r)
		if err != nil {
			h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		l, err := utils.DeCodeBody[model.Login](r.Body, h.validate)
		if err != nil {
			h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		lr, err := h.userService.Login(ctx, l, ip)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson(w, model.API[model.LoginRep]{
			Code: 0,
			Data: lr,
		})
	}
}

func (h *UserHandel) UserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		u, err := h.userService.Info(ctx, t)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson(w, model.API[model.UserInfo]{
			Code: 0,
			Data: u,
		})
	}
}

func (h *UserHandel) ChangePasswd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)

		c, err := utils.DeCodeBody[model.ChangePasswd](r.Body, h.validate)
		if err != nil {
			h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		err = h.userService.ChangePasswd(ctx, c, t.UID, true)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})

	}
}

func (h *UserHandel) ChangeName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		c, err := utils.DeCodeBody[model.ChangeName](r.Body, h.validate)
		if err != nil {
			h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		err = h.userService.ChangeName(ctx, c.Name, t)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}

func (h *UserHandel) PutTexture() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := ctx.Value(tokenKey).(*model.TokenClaims)
		models := r.FormValue("model")

		textureType := chi.URLParamFromCtx(ctx, "textureType")
		if textureType != "skin" && textureType != "cape" {
			h.logger.DebugContext(ctx, "上传类型错误")
			h.handleError.Error(ctx, w, "上传类型错误", model.ErrInput, 400, slog.LevelDebug)
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
			h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}

		switch models {
		case "slim":
		case "":
		default:
			h.logger.DebugContext(ctx, "错误的皮肤的材质模型")
			h.handleError.Error(ctx, w, "错误的皮肤的材质模型", model.ErrInput, 400, slog.LevelDebug)
			return
		}

		err = h.textureService.PutTexture(ctx, t, skin, models, textureType)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}

func (h *UserHandel) NeedEnableEmail(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !h.config.Email.Enable {
			h.handleError.Error(ctx, w, "未开启邮件功能", model.ErrUnknown, 403, slog.LevelInfo)
		}
		handle.ServeHTTP(w, r)
	})
}

var ErrNotAllowDomain = errors.New("不在允许域名列表内")

func (h *UserHandel) SendRegEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, ip, shouldReturn := h.sendMailParameter(ctx, r, w)
		if shouldReturn {
			return
		}

		if len(h.config.Email.AllowDomain) != 0 {
			allow := false
			for _, v := range h.config.Email.AllowDomain {
				if strings.HasSuffix(c.Email, v) {
					allow = true
					break
				}
			}
			if !allow {
				h.handleError.Error(ctx, w, "不在允许邮箱域名内", model.ErrInput, 400, slog.LevelDebug)
				return
			}
		}
		if h.config.Email.EmailReg != "" {
			r, err := h.emailReg()
			if err != nil {
				h.handleError.Error(ctx, w, "正则错误", model.ErrUnknown, 500, slog.LevelError)
				return
			}
			if !r.MatchString(c.Email) {
				h.handleError.Error(ctx, w, "邮箱不符合正则要求", model.ErrInput, 400, slog.LevelDebug)
				return
			}
		}

		err := h.userService.SendRegEmail(ctx, c.Email, c.CaptchaToken, r.Host, ip)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}

func (h *UserHandel) SendForgotEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, ip, shouldReturn := h.sendMailParameter(ctx, r, w)
		if shouldReturn {
			return
		}

		err := h.userService.SendChangePasswordEmail(ctx, c.Email, c.CaptchaToken, r.Host, ip)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}

func (h *UserHandel) sendMailParameter(ctx context.Context, r *http.Request, w http.ResponseWriter) (model.SendRegEmail, string, bool) {
	c, err := utils.DeCodeBody[model.SendRegEmail](r.Body, h.validate)
	if err != nil {
		h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
		return model.SendRegEmail{}, "", true
	}
	ip, err := utils.GetIP(r)
	if err != nil {
		h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
		return model.SendRegEmail{}, "", true
	}
	return c, ip, false
}

func (h *UserHandel) ForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, err := utils.DeCodeBody[model.ForgotPassword](r.Body, h.validate)
		if err != nil {
			h.handleError.Error(ctx, w, err.Error(), model.ErrInput, 400, slog.LevelDebug)
			return
		}
		err = h.userService.ForgotPassword(ctx, c.Email, c.PassWord, c.EmailJwt)
		if err != nil {
			h.handleError.Service(ctx, w, err)
			return
		}
		encodeJson(w, model.API[any]{
			Code: 0,
		})
	}
}
