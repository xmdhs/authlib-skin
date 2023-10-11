package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"flag"
	"fmt"
	"os"

	_ "embed"

	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/server"
	"github.com/xmdhs/authlib-skin/utils/sign"
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
			lo.Must0(os.WriteFile(configPath, configTempLate, 0600))
			fmt.Println("未找到配置文件，已写入模板配置文件")
			return
		}
		panic(err)
	}
	config := lo.Must(config.YamlDeCode(b))

	if config.RsaPriKey == "" {
		rsa2048 := lo.Must(rsa.GenerateKey(rand.Reader, 4096))
		as := sign.NewAuthlibSignWithKey(rsa2048)
		config.RsaPriKey = lo.Must(as.GetPriKey())
		lo.Must0(os.WriteFile(configPath, []byte(config.RsaPriKey), 0600))
	}

	s, cancel := lo.Must2(server.InitializeRoute(ctx, config))
	defer cancel()
	panic(s.ListenAndServe())
}
