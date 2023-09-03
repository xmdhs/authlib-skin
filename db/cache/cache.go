package cache

import "time"

type Cache interface {
	Del(k []byte) error
	Get(k []byte) ([]byte, error)
	Put(k []byte, v []byte, timeOut time.Time) error
}
