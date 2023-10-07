package server

import (
	"net/http"
	"time"

	"github.com/xmdhs/authlib-skin/config"
)

func NewServer(c config.Config, route http.Handler) (*http.Server, func()) {
	s := &http.Server{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      20 * time.Second,
		Addr:              c.Port,
		Handler:           route,
	}
	return s, func() { s.Close() }
}
