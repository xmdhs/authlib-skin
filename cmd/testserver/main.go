package main

import (
	"context"
	"fmt"

	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/server"
)

func main() {
	ctx := context.Background()
	config := config.Config{
		OfflineUUID: true,
		Port:        "127.0.0.1:8080",
		Log: struct {
			Level string
			Json  bool
		}{
			Level: "debug",
		},
		Sql: struct{ MysqlDsn string }{
			MysqlDsn: "sizzle1445:jnjjJQ8^YF&8PN@tcp(192.168.20.1)/skin",
		},
		Node:  0,
		Epoch: 1693645718534,
	}
	s, c, err := server.InitializeRoute(ctx, config)
	if err != nil {
		panic(err)
	}
	defer c()
	fmt.Println(s.ListenAndServe())
}
