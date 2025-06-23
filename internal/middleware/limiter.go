package middleware

import (
	"net/http"
	"rate-limiter/internal/limiter"
	"strings"
)

type RateLimiterMiddleware struct {
	limiter *limiter.Limiter
}

func NewRateLimiterMiddleware(limiter *limiter.Limiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiter: limiter,
	}
}

func (m *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")
		identifier := ""
		isToken := false

		if token != "" {
			identifier = token
			isToken = true
		} else {
			identifier = strings.Split(r.RemoteAddr, ":")[0]
		}

		allowed, err := m.limiter.Allow(identifier, isToken)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
