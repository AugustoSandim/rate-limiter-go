package main

import (
	"log"
	"net/http"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/middleware"
	"rate-limiter/internal/storage"
	"rate-limiter/pkg/config"
)

func main() {
	cfg := config.Load()

	redisStorage, err := storage.NewRedisStorage(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	rateLimiter := limiter.NewLimiter(redisStorage, cfg)
	limiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiter)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	log.Println("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", limiterMiddleware.Handler(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
