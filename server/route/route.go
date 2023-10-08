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
)

func NewRoute(handelY *yggdrasil.Yggdrasil, handel *handle.Handel, c config.Config, sl slog.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	if sl.Enabled(context.Background(), slog.LevelDebug) {
		r.Use(NewStructuredLogger(sl))
	}
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)
	if c.RaelIP {
		r.Use(middleware.RealIP)
	}

	r.Mount("/api/v1", newSkinApi(handel))
	r.Mount("/api/yggdrasil", newYggdrasil(handelY))

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

		r.Put("/api/user/profile/{uuid}/{textureType}", handelY.PutTexture())
		r.Delete("/api/user/profile/{uuid}/{textureType}", handelY.DelTexture())

		r.Post("/sessionserver/session/minecraft/join", handelY.SessionJoin())
		r.Post("/minecraftservices/player/certificates", handelY.PlayerCertificates())

	})

	r.Post("/authserver/authenticate", handelY.Authenticate())
	r.Post("/authserver/signout", handelY.Signout())

	r.Get("/sessionserver/session/minecraft/profile/{uuid}", handelY.GetProfile())
	r.Post("/api/profiles/minecraft", handelY.BatchProfile())

	r.Get("/sessionserver/session/minecraft/hasJoined", handelY.HasJoined())

	r.Get("/", handelY.YggdrasilRoot())
	return r
}

func newSkinApi(handel *handle.Handel) http.Handler {
	r := chi.NewRouter()

	r.Put("/user/reg", handel.Reg())
	r.Get("/config", handel.GetConfig())

	r.Group(func(r chi.Router) {
		r.Use(handel.NeedAuth)
		r.Get("/user", handel.UserInfo())
		r.Post("/user/password", handel.ChangePasswd())
		r.Post("/user/name", handel.ChangeName())

		r.Group(func(r chi.Router) {
			r.Use(handel.NeedAdmin)
			r.Get("/admin/users", handel.ListUser())
		})
	})

	return r
}
