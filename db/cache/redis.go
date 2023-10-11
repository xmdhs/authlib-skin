//go:build redis

package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ Cache = (*RedisCache)(nil)

type RedisCache struct {
	c *redis.Client
}

func NewRedis(addr, pass string) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})
	return &RedisCache{c: rdb}
}

func (r *RedisCache) Del(k []byte) error {
	_, err := r.c.Del(context.Background(), string(k)).Result()
	if err != nil {
		return fmt.Errorf("RedisCache.Del: %w", err)
	}
	return nil
}

func (r *RedisCache) Get(k []byte) ([]byte, error) {
	value, err := r.c.Get(context.Background(), string(k)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("RedisCache.Get: %w", err)
	}
	return value, nil
}

func (r *RedisCache) Put(k []byte, v []byte, timeOut time.Time) error {
	err := r.c.Set(context.Background(), string(k), v, timeOut.Sub(time.Now())).Err()
	if err != nil {
		return fmt.Errorf("RedisCache.Put: %w", err)
	}
	return nil
}
