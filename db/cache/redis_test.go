//go:build redis

package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	c := NewRedis("127.0.0.1:6379", "")

	key := []byte("key")
	value := []byte("value")

	require.Nil(t, c.Put(key, value, time.Now().Add(1*time.Hour)))

	v, err := c.Get(key)
	require.Nil(t, err)

	assert.Equal(t, v, value)

	require.Nil(t, c.Del(key))

	v, err = c.Get(key)
	require.Nil(t, err)
	require.Nil(t, v)

	require.Nil(t, c.Put(key, value, time.Now().Add(2*time.Second)))
	time.Sleep(3 * time.Second)

	v, err = c.Get(key)
	require.Nil(t, err)
	require.Nil(t, v)
}
