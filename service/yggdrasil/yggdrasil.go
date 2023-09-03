package yggdrasil

import (
	"github.com/xmdhs/authlib-skin/config"
	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
)

type Yggdrasil struct {
	client *ent.Client
	cache  cache.Cache
	config config.Config
}

func NewYggdrasil(client *ent.Client, cache cache.Cache, c config.Config) *Yggdrasil {
	return &Yggdrasil{
		client: client,
		cache:  cache,
		config: c,
	}
}
