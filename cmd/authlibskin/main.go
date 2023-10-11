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

	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/server"
	"github.com/xmdhs/authlib-skin/utils/sign"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "config.toml", "")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	b, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			rsa2048 := lo.Must(rsa.GenerateKey(rand.Reader, 4096))
			as := sign.NewAuthlibSignWithKey(rsa2048)

			c := config.Default()
			c.RsaPriKey = lo.Must(as.GetPriKey())

			lo.Must0(os.WriteFile(configPath, lo.Must(toml.Marshal(c)), 0600))
			fmt.Println("未找到配置文件，已写入模板配置文件")
			return
		}
		panic(err)
	}
	var config config.Config
	lo.Must0(toml.Unmarshal(b, &config))
	s, cancel := lo.Must2(server.InitializeRoute(ctx, config))
	defer cancel()
	panic(s.ListenAndServe())
}
