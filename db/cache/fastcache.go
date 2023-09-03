package cache

import (
	"fmt"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/alecthomas/binary"
)

var _ Cache = (*FastCache)(nil)

type FastCache struct {
	c *fastcache.Cache
}

type ttlCache struct {
	TimeOut int64
	V       []byte
}

func NewFastCache(maxBytes int) *FastCache {
	c := fastcache.New(maxBytes)
	return &FastCache{c: c}
}

func (f *FastCache) Put(k []byte, v []byte, timeOut time.Time) error {
	b, err := binary.Marshal(ttlCache{V: v, TimeOut: timeOut.Unix()})
	if err != nil {
		return fmt.Errorf("FastCache.Put: %w", err)
	}
	f.c.SetBig(k, b)
	return nil
}

func (f *FastCache) Del(k []byte) error {
	f.c.Del(k)
	return nil
}

func (f *FastCache) Get(k []byte) ([]byte, error) {
	b := f.c.GetBig(nil, k)
	me := ttlCache{}
	err := binary.Unmarshal(b, &me)
	if err != nil {
		return nil, fmt.Errorf("FastCache.Get: %w", err)
	}
	return me.V, nil
}
