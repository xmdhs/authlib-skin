package server

import (
	"context"
	"log/slog"
)

type reqInfo struct {
	URL     string
	IP      string
	TrackId uint64
}

type reqInfoKeyType string

var reqinfoKey reqInfoKeyType = "reqinfoKey"

func setCtx(ctx context.Context, r *reqInfo) context.Context {
	return context.WithValue(ctx, reqinfoKey, r)
}

func getFromCtx(ctx context.Context) *reqInfo {
	v := ctx.Value(reqinfoKey)
	if v == nil {
		return nil
	}
	return v.(*reqInfo)
}

type warpSlogHandle struct {
	slog.Handler
}

func (w *warpSlogHandle) Handle(ctx context.Context, r slog.Record) error {
	if w.Enabled(ctx, slog.LevelInfo) {
		ri := getFromCtx(ctx)
		if ri != nil {
			r.AddAttrs(slog.String("ip", ri.IP), slog.String("url", ri.URL), slog.Uint64("trackID", ri.TrackId))
		}
	}
	return w.Handler.Handle(ctx, r)
}

func NewSlog(h slog.Handler) *slog.Logger {
	l := slog.New(&warpSlogHandle{
		Handler: h,
	})
	return l
}
