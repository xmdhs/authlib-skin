package server

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

type warpSlogHandle struct {
	slog.Handler
}

func (w *warpSlogHandle) Handle(ctx context.Context, r slog.Record) error {
	id := middleware.GetReqID(ctx)
	if id != "" {
		r.AddAttrs(slog.String("trackID", id))
	}
	return w.Handler.Handle(ctx, r)
}

func NewSlog(h slog.Handler) *slog.Logger {
	l := slog.New(&warpSlogHandle{
		Handler: h,
	})
	return l
}
