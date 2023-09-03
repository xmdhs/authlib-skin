package yggdrasil

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/xmdhs/authlib-skin/db/cache"
	"github.com/xmdhs/authlib-skin/db/ent"
	"github.com/xmdhs/authlib-skin/model/yggdrasil"
)

var (
	ErrRate = errors.New("频率限制")
)

func Authenticate(cxt context.Context, client *ent.Client, auth yggdrasil.Authenticate, cache cache.Cache) error {
	key := []byte("Authenticate" + auth.Username)

	v, err := cache.Get(key)
	if err != nil {
		return fmt.Errorf("Authenticate: %w", err)
	}
	if v != nil {
		u := binary.BigEndian.Uint64(v)
		t := time.Unix(int64(u), 0)
		if time.Now().Before(t) {
			return fmt.Errorf("Authenticate: %w", ErrRate)
		}
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(time.Now().Add(10*time.Second).Unix()))
	err = cache.Put(key, b, time.Now().Add(20*time.Second))
	if err != nil {
		return fmt.Errorf("Authenticate: %w", err)
	}

	return nil
}
