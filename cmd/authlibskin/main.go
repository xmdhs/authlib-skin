package main

import (
	"context"
	"flag"
	"os"

	"github.com/xmdhs/authlib-skin/config"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "", "")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	os.ReadFile(configPath)

	config.YamlDeCode()
}
