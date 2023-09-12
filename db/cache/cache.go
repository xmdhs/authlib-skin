package cache

import (
	"encoding/json"
	"time"
)

type Cache interface {
	Del(k []byte) error
	Get(k []byte) ([]byte, error)
	Put(k []byte, v []byte, timeOut time.Time) error
}

type CacheHelp[T any] struct {
	Cache
}

func (c CacheHelp[T]) Get(k []byte) (T, error) {
	var t T
	b, err := c.Cache.Get(k)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (c CacheHelp[T]) Put(k []byte, v T, timeOut time.Time) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.Cache.Put(k, b, timeOut)
}
