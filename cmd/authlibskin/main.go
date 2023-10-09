package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	_ "embed"

	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/server"
)

var configPath string

//go:embed config.yaml.template
var configTempLate []byte

func init() {
	flag.StringVar(&configPath, "c", "config.yaml", "")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	b, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			lo.Must0(os.WriteFile("config.yaml", configTempLate, 0600))
			fmt.Println("未找到配置文件，已写入模板配置文件")
			return
		}
		panic(err)
	}
	config := lo.Must(config.YamlDeCode(b))

	s, cancel := lo.Must2(server.InitializeRoute(ctx, config))
	defer cancel()
	panic(s.ListenAndServe())
}
