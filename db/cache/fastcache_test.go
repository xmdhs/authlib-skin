package cache

import (
	"testing"
	"time"

	"github.com/samber/lo"
)

func TestFastCache(t *testing.T) {
	f := NewFastCache(100000)
	c := CacheHelp[string]{
		Cache: f,
	}
	c.Put([]byte("123"), "123", time.Now().Add(10*time.Second))

	if lo.Must(c.Get([]byte("123"))) != "123" {
		t.FailNow()
	}
}
