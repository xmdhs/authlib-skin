package route

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/handle"
	"github.com/xmdhs/authlib-skin/handle/yggdrasil"
)

func NewRoute(yggService *yggdrasil.Yggdrasil, handel *handle.Handel) (*httprouter.Router, error) {
	r := httprouter.New()
	err := newYggdrasil(r, *yggService)
	if err != nil {
		return nil, fmt.Errorf("NewRoute: %w", err)
	}
	err = newSkinApi(r, handel)
	if err != nil {
		return nil, fmt.Errorf("NewRoute: %w", err)
	}
	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
	return r, nil
}

func newYggdrasil(r *httprouter.Router, handelY yggdrasil.Yggdrasil) error {
	r.POST("/api/yggdrasil/authserver/authenticate", warpHJSON(handelY.Authenticate()))
	r.POST("/api/yggdrasil/authserver/validate", warpHJSON(handelY.Validate()))
	r.POST("/api/yggdrasil/authserver/signout", warpHJSON(handelY.Signout()))
	r.POST("/api/yggdrasil/authserver/invalidate", handelY.Invalidate())
	r.POST("/api/yggdrasil/authserver/refresh", warpHJSON(handelY.Refresh()))

	r.PUT("/api/yggdrasil/api/user/profile/:uuid/:textureType", handelY.PutTexture())
	r.DELETE("/api/yggdrasil/api/user/profile/:uuid/:textureType", warpHJSON(handelY.DelTexture()))

	r.GET("/api/yggdrasil/sessionserver/session/minecraft/profile/:uuid", warpHJSON(handelY.GetProfile()))
	r.POST("/api/yggdrasil/api/profiles/minecraft", warpHJSON(handelY.BatchProfile()))

	r.POST("/api/yggdrasil/sessionserver/session/minecraft/join", warpHJSON(handelY.SessionJoin()))
	r.GET("/api/yggdrasil/sessionserver/session/minecraft/hasJoined", warpHJSON(handelY.HasJoined()))

	r.POST("/api/yggdrasil/minecraftservices/player/certificates", warpHJSON(handelY.PlayerCertificates()))

	r.GET("/api/yggdrasil", warpHJSON(handelY.YggdrasilRoot()))
	r.GET("/api/yggdrasil/", warpHJSON(handelY.YggdrasilRoot()))

	r.GET("/texture/*filepath", handelY.TextureAssets())
	return nil
}

func newSkinApi(r *httprouter.Router, handel *handle.Handel) error {
	r.PUT("/api/v1/user/reg", handel.Reg())
	return nil
}
