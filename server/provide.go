package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/bwmarrin/snowflake"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/mysql"
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
		return nil, nil, fmt.Errorf("newDB: %w", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return db, func() { db.Close() }, nil
}

func ProvideQuerier(ctx context.Context, db *sql.DB) (mysql.Querier, func(), error) {
	q, err := mysql.Prepare(ctx, db)
	if err != nil {
		return nil, nil, fmt.Errorf("newQuerier: %w", err)
	}
	return q, func() { q.Close() }, nil
}

func ProvideValidate() *validator.Validate {
	return validator.New()
}

func ProvideSnowflake(c config.Config) (*snowflake.Node, error) {
	snowflake.Epoch = c.Epoch
	n, err := snowflake.NewNode(c.Node)
	if err != nil {
		return nil, fmt.Errorf("newSnowflake: %w", err)
	}
	return n, nil
}

var Set = wire.NewSet(ProvideSlog, ProvideDB, ProvideQuerier, ProvideValidate, ProvideSnowflake)
