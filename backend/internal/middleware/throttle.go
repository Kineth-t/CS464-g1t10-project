package middleware

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// SlidingWindowThrottle implements a sliding window rate limiter using Redis.
// It supports both Global (total app) and Per-User (IP-based) limits.
func SlidingWindowThrottle(rdb *redis.Client, limit int, window time.Duration, contextName string, isGlobal bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx context.Context = r.Context()
			var key string

			if isGlobal {
				// Global: All users share the same bucket (Protects DB/Infrastructure)
				key = "throttle:global:" + contextName
			} else {
				// Per-User: Each IP gets its own bucket (Protects against Bots/Scrapers)
				ip, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					// Fallback if RemoteAddr doesn't have a port (e.g., local testing)
					ip = r.RemoteAddr
				}
				key = "throttle:user:" + contextName + ":" + ip
			}

			// If Redis is not configured, skip rate limiting entirely
			if rdb == nil {
				next.ServeHTTP(w, r)
				return
			}

			now := time.Now().UnixNano()
			// Calculate the start of the window (e.g., 60 seconds ago)
			threshold := now - window.Nanoseconds()

			// Use a Pipeline to reduce network round-trips to Redis
			pipe := rdb.Pipeline()
			
			// Remove old requests outside the current window
			pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(threshold, 10))
			
			// Add current request timestamp
			pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
			
			// Set expiry so Redis cleans up unused IP keys automatically
			pipe.Expire(ctx, key, window)
			
			// Count remaining items in the window
			countCmd := pipe.ZCard(ctx, key)

			_, err := pipe.Exec(ctx)
			if err != nil {
				// Fail-safe: If Redis is down, we allow the request to pass
				// rather than blocking all users (Availability > Strict Limiting)
				next.ServeHTTP(w, r)
				return
			}

			// Check if limit exceeded
			if countCmd.Val() > int64(limit) {
				msg := "Too many requests. Please wait."
				if isGlobal {
					msg = "Our servers are currently busy due to high demand. Please try again in a few seconds."
				}

				// Standard HTTP header to tell clients when they can retry
				w.Header().Set("X-Retry-After", strconv.FormatInt(int64(window.Seconds()), 10))
				http.Error(w, msg, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}