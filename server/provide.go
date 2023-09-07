package server

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/db/ent/migrate"
	"github.com/xmdhs/authlib-skin/utils/sign"
)

func ProvideSlog(c config.Config) slog.Handler {
	var level slog.Level
	switch c.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	o := &slog.HandlerOptions{Level: level}

	var h slog.Handler
	if c.Log.Json {
		h = slog.NewJSONHandler(os.Stderr, o)
	} else {
		h = slog.NewTextHandler(os.Stderr, o)
	}

	return h
}

func ProvideDB(c config.Config) (*sql.DB, func(), error) {
	db, err := sql.Open("mysql", c.Sql.MysqlDsn)
	if err != nil {
		return nil, nil, fmt.Errorf("ProvideDB: %w", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return db, func() { db.Close() }, nil
}

func ProvideEnt(ctx context.Context, db *sql.DB, c config.Config, sl *slog.Logger) (*ent.Client, func(), error) {
	drv := entsql.OpenDB("mysql", db)
	opts := []ent.Option{ent.Driver(drv), ent.Log(
		func(a ...any) {
			sl.Debug(fmt.Sprint(a))
		},
	)}
	if c.Debug {
		opts = append(opts, ent.Debug())
	}
	e := ent.NewClient(opts...)
	err := e.Schema.Create(ctx, migrate.WithForeignKeys(false), migrate.WithDropIndex(true), migrate.WithDropColumn(true))
	if err != nil {
		return nil, nil, fmt.Errorf("ProvideEnt: %w", err)
	}
	return e, func() { e.Close() }, nil
}

func ProvideValidate() *validator.Validate {
	return validator.New()
}

func ProvideCache(c config.Config) cache.Cache {
	return cache.NewFastCache(c.Cache.Ram)
}

func ProvidePriKey(c config.Config) (*rsa.PrivateKey, error) {
	a, err := sign.NewAuthlibSign([]byte(c.RsaPriKey))
	if err != nil {
		return nil, fmt.Errorf("ProvidePriKey: %w", err)
	}
	return a.GetKey(), nil
}

var Set = wire.NewSet(ProvideSlog, ProvideDB, ProvideEnt, ProvideValidate, ProvideCache, ProvidePriKey)
