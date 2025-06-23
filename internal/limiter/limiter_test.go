package limiter

import (
	"rate-limiter/pkg/config"
	"testing"
	"time"
)

type mockStorage struct {
	store map[string]int
}

func (m *mockStorage) Increment(key string) (int, error) {
	m.store[key]++
	return m.store[key], nil
}

func (m *mockStorage) Set(key string, value int, expiration time.Duration) error {
	m.store[key] = value
	return nil
}

func (m *mockStorage) Get(key string) (int, error) {
	return m.store[key], nil
}

func newMockStorage() *mockStorage {
	return &mockStorage{store: make(map[string]int)}
}

func TestLimiter(t *testing.T) {
	cfg := &config.Config{
		RateLimitByIP:      5,
		RateLimitByToken:   10,
		BlockTimeInSeconds: 60,
		TokensRateLimit: map[string]int{
			"test_token": 2,
		},
	}

	t.Run("IP based limiting", func(t *testing.T) {
		storage := newMockStorage()
		limiter := NewLimiter(storage, cfg)

		for i := 0; i < 5; i++ {
			allowed, err := limiter.Allow("127.0.0.1", false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !allowed {
				t.Fatalf("request should be allowed")
			}
		}

		allowed, err := limiter.Allow("127.0.0.1", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if allowed {
			t.Fatalf("request should be denied")
		}
	})

	t.Run("Token based limiting", func(t *testing.T) {
		storage := newMockStorage()
		limiter := NewLimiter(storage, cfg)

		for i := 0; i < 10; i++ {
			allowed, err := limiter.Allow("some_other_token", true)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !allowed {
				t.Fatalf("request should be allowed")
			}
		}

		allowed, err := limiter.Allow("some_other_token", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if allowed {
			t.Fatalf("request should be denied")
		}
	})

	t.Run("Specific token based limiting", func(t *testing.T) {
		storage := newMockStorage()
		limiter := NewLimiter(storage, cfg)

		for i := 0; i < 2; i++ {
			allowed, err := limiter.Allow("test_token", true)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !allowed {
				t.Fatalf("request should be allowed")
			}
		}

		allowed, err := limiter.Allow("test_token", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if allowed {
			t.Fatalf("request should be denied")
		}
	})
}
