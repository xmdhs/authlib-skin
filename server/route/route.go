package route

import (
	"fmt"

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
	r.POST("/api/authserver/authenticate", warpHJSON(handelY.Authenticate()))
	r.POST("/api/authserver/validate", warpHJSON(handelY.Validate()))
	return nil
}

func newSkinApi(r *httprouter.Router, handel *handle.Handel) error {
	r.PUT("/api/v1/user/reg", handel.Reg())
	return nil
}
