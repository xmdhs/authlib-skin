package route

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/handle"
	"github.com/xmdhs/authlib-skin/handle/yggdrasil"
	"github.com/xmdhs/authlib-skin/server/static"
)

func NewRoute(handelY *yggdrasil.Yggdrasil, handel *handle.Handel, c config.Config, sl slog.Handler,
	userHandel *handle.UserHandel, adminHandel *handle.AdminHandel) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	if c.RaelIP {
		r.Use(middleware.RealIP)
	}
	if sl.Enabled(context.Background(), slog.LevelDebug) {
		r.Use(NewStructuredLogger(sl))
	}
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)
	r.Use(APILocationIndication)

	r.Mount("/", static.StaticServer())
	r.Mount("/api/v1", newSkinApi(handel, userHandel, adminHandel))
	r.Mount("/api/yggdrasil", newYggdrasil(handelY))

	if c.Debug {
		r.Mount("/debug", middleware.Profiler())
	}

	r.Get("/texture/*", handelY.TextureAssets())

	return r
}

func newYggdrasil(handelY *yggdrasil.Yggdrasil) http.Handler {
	r := chi.NewRouter()
	r.Use(warpHJSON)

	r.Group(func(r chi.Router) {
		r.Use(handelY.Auth)
		r.Post("/authserver/validate", handelY.Validate())
		r.Post("/authserver/invalidate", handelY.Invalidate())
		r.Post("/authserver/refresh", handelY.Refresh())

		r.Delete("/api/user/profile/{uuid}/{textureType}", handelY.DelTexture())

		r.Post("/sessionserver/session/minecraft/join", handelY.SessionJoin())
		r.Post("/minecraftservices/player/certificates", handelY.PlayerCertificates())

	})

	r.Post("/authserver/authenticate", handelY.Authenticate())
	r.Post("/authserver/signout", handelY.Signout())

	r.Get("/sessionserver/session/minecraft/profile/{uuid}", handelY.GetProfile())
	r.Post("/api/profiles/minecraft", handelY.BatchProfile())

	r.Get("/sessionserver/session/minecraft/hasJoined", handelY.HasJoined())

	r.Get("/minecraftservices/player/attributes", handelY.PlayerAttributes())
	r.Post("/minecraftservices/player/attributes", handelY.PlayerAttributes())
	r.Post("/minecraftservices/player/report", handelY.PlayerReport())
	r.Get("/minecraftservices/publickeys", handelY.PublicKeys())

	r.Get("/", handelY.YggdrasilRoot())
	return r
}

func newSkinApi(handel *handle.Handel, userHandel *handle.UserHandel, adminHandel *handle.AdminHandel) http.Handler {
	r := chi.NewRouter()

	r.Post("/user/reg", userHandel.Reg())
	r.Post("/user/login", userHandel.Login())
	r.Get("/config", handel.GetConfig())

	r.Group(func(r chi.Router) {
		r.Use(userHandel.NeedEnableEmail)
		r.Post("/user/reg_email", userHandel.SendRegEmail())
		r.Post("/user/forgot_email", userHandel.SendForgotEmail())
		r.Post("/user/forgot", userHandel.ForgotPassword())
	})

	r.Group(func(r chi.Router) {
		r.Use(adminHandel.NeedAuth)
		r.Get("/user", userHandel.UserInfo())
		r.Post("/user/password", userHandel.ChangePasswd())
		r.Post("/user/name", userHandel.ChangeName())
		r.Put("/user/skin/{textureType}", userHandel.PutTexture())
	})

	r.Group(func(r chi.Router) {
		r.Use(adminHandel.NeedAuth)
		r.Use(adminHandel.NeedAdmin)
		r.Get("/admin/users", adminHandel.ListUser())
		r.Patch("/admin/user/{uid}", adminHandel.EditUser())
	})

	return r
}
