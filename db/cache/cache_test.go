package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	cl := []Cache{NewFastCache(100000), NewRedis("127.0.0.1:6379", "")}

	for i, c := range cl {
		key := []byte("key")
		value := []byte("value")

		require.Nil(t, c.Put(key, value, time.Now().Add(1*time.Hour)), i)

		v, err := c.Get(key)
		require.Nil(t, err, i)

		assert.Equal(t, v, value, i)

		require.Nil(t, c.Del(key), i)

		v, err = c.Get(key)
		require.Nil(t, err, i)
		require.Nil(t, v, i)

		require.Nil(t, c.Put(key, value, time.Now().Add(2*time.Second)), i)
		time.Sleep(3 * time.Second)

		v, err = c.Get(key)
		require.Nil(t, err, i)
		require.Nil(t, v, i)
	}
}
