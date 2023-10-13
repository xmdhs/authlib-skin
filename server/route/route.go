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

func NewRoute(handelY *yggdrasil.Yggdrasil, handel *handle.Handel, c config.Config, sl slog.Handler) http.Handler {
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
	r.Mount("/api/v1", newSkinApi(handel))
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

	r.Get("/", handelY.YggdrasilRoot())
	return r
}

func newSkinApi(handel *handle.Handel) http.Handler {
	r := chi.NewRouter()

	r.Post("/user/reg", handel.Reg())
	r.Post("/user/login", handel.Login())
	r.Get("/config", handel.GetConfig())

	r.Group(func(r chi.Router) {
		r.Use(handel.NeedAuth)
		r.Get("/user", handel.UserInfo())
		r.Post("/user/password", handel.ChangePasswd())
		r.Post("/user/name", handel.ChangeName())
		r.Put("/user/skin/{textureType}", handel.PutTexture())
	})

	r.Group(func(r chi.Router) {
		r.Use(handel.NeedAuth)
		r.Use(handel.NeedAdmin)
		r.Get("/admin/users", handel.ListUser())
		r.Patch("/admin/user/{uid}", handel.EditUser())
	})

	return r
}
