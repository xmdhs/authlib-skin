package test

import (
	"context"
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/server"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	b := lo.Must(os.ReadFile("config.toml"))
	var config config.Config
	lo.Must0(toml.Unmarshal(b, &config))
	s, cancel := lo.Must2(server.InitializeRoute(ctx, config))
	defer cancel()
	go func() {
		s.ListenAndServe()
	}()

	os.Exit(m.Run())
}
