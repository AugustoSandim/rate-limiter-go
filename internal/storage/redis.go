package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr string) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &RedisStorage{client: client}, nil
}

func (r *RedisStorage) Increment(key string) (int, error) {
	return int(r.client.Incr(context.Background(), key).Val()), nil
}

func (r *RedisStorage) Set(key string, value int, expiration time.Duration) error {
	return r.client.Set(context.Background(), key, value, expiration).Err()
}

func (r *RedisStorage) Get(key string) (int, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}
