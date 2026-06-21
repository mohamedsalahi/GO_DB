package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	rdb        *redis.Client
	maxReqs    int
	windowSize time.Duration
}

func NewRateLimiter(rdb *redis.Client, maxReqs int, windowSize time.Duration) *RateLimiter {
	return &RateLimiter{
		rdb:        rdb,
		maxReqs:    maxReqs,
		windowSize: windowSize,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := "ratelimit:" + r.RemoteAddr

		count, err := rl.rdb.Incr(r.Context(), key).Result()
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if count == 1 {
			rl.rdb.Expire(r.Context(), key, rl.windowSize)
		}

		remaining := rl.maxReqs - int(count)
		if remaining < 0 {
			remaining = 0
		}

		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.maxReqs))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rl.windowSize).Unix(), 10))

		if count > int64(rl.maxReqs) {
			w.Header().Set("Retry-After", strconv.Itoa(int(rl.windowSize.Seconds())))
			http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
