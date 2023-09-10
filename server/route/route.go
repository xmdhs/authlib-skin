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

	r.GET("/api/yggdrasil/sessionserver/session/minecraft/profile/:uuid", handelY.GetProfile())

	r.GET("/api/yggdrasil", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte(`{
			"meta": {
				"serverName": "test",
				"implementationName": "test",
				"implementationVersion": "999.999.999"
			},
			"skinDomains": [
			],
			"signaturePublickey": "123"
		}`))
	})
	return nil
}

func newSkinApi(r *httprouter.Router, handel *handle.Handel) error {
	r.PUT("/api/v1/user/reg", handel.Reg())
	return nil
}
