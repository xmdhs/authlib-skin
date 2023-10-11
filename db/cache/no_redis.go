//go:build !redis

package cache

func NewRedis(addr, pass string) Cache {
	panic("not tag redis")
}
