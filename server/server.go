package server

import (
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/utils"
)

func NewServer(c config.Config, sl *slog.Logger, route *httprouter.Router) (*http.Server, func()) {
	trackid := atomic.Uint64{}
	s := &http.Server{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      20 * time.Second,
		Addr:              c.Port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if sl.Enabled(ctx, slog.LevelInfo) {
				ip, _ := utils.GetIP(r, c.RaelIP)
				trackid.Add(1)
				ctx = setCtx(ctx, &reqInfo{
					URL:     r.URL.String(),
					IP:      ip,
					TrackId: trackid.Load(),
				})
				r = r.WithContext(ctx)
			}
			if c.Debug && sl.Enabled(ctx, slog.LevelDebug) {
				sl.DebugContext(ctx, r.Method)
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			route.ServeHTTP(w, r)
		}),
	}
	return s, func() { s.Close() }
}
