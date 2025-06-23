package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	RedisAddr          string
	RateLimitByIP      int
	RateLimitByToken   int
	BlockTimeInSeconds int
	TokensRateLimit    map[string]int
}

var (
	instance *Config
	once     sync.Once
)

func Load() *Config {
	once.Do(func() {
		instance = &Config{
			RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
			RateLimitByIP:      getEnvAsInt("RATE_LIMIT_BY_IP", 10),
			RateLimitByToken:   getEnvAsInt("RATE_LIMIT_BY_TOKEN", 100),
			BlockTimeInSeconds: getEnvAsInt("BLOCK_TIME_IN_SECONDS", 300),
			TokensRateLimit:    getTokensRateLimit(),
		}
	})
	return instance
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}

func getTokensRateLimit() map[string]int {
	tokens := make(map[string]int)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "TOKEN_") && strings.HasSuffix(pair[0], "_RATE_LIMIT") {
			token := strings.TrimSuffix(strings.TrimPrefix(pair[0], "TOKEN_"), "_RATE_LIMIT")
			limit, err := strconv.Atoi(pair[1])
			if err == nil {
				tokens[token] = limit
			}
		}
	}
	return tokens
}
