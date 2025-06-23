package storage

import "time"

type Storage interface {
	Increment(key string) (int, error)
	Set(key string, value int, expiration time.Duration) error
	Get(key string) (int, error)
}
