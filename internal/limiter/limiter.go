package limiter

import (
	"fmt"
	"rate-limiter/internal/storage"
	"rate-limiter/pkg/config"
	"time"
)

type Limiter struct {
	storage storage.Storage
	config  *config.Config
}

func NewLimiter(storage storage.Storage, config *config.Config) *Limiter {
	return &Limiter{
		storage: storage,
		config:  config,
	}
}

func (l *Limiter) Allow(identifier string, isToken bool) (bool, error) {
	limit := l.getLimit(identifier, isToken)
	key := l.getKey(identifier, isToken)

	count, err := l.storage.Get(key)
	if err != nil {
		return false, err
	}

	if count >= limit {
		return false, nil
	}

	count, err = l.storage.Increment(key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		expiration := time.Duration(l.config.BlockTimeInSeconds) * time.Second
		if err := l.storage.Set(key, count, expiration); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (l *Limiter) getLimit(identifier string, isToken bool) int {
	if isToken {
		if limit, ok := l.config.TokensRateLimit[identifier]; ok {
			return limit
		}
		return l.config.RateLimitByToken
	}
	return l.config.RateLimitByIP
}

func (l *Limiter) getKey(identifier string, isToken bool) string {
	prefix := "ip"
	if isToken {
		prefix = "token"
	}
	return fmt.Sprintf("%s:%s", prefix, identifier)
}
