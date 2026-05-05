package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"lakoo/backend/pkg/response"
)

// RateLimitMiddleware blocks IPs that exceed a certain amount of requests per window.
// Specifically tailored for Login brute-force protection.
func RateLimitMiddleware(rdb *redis.Client, maxAttempts int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:login:%s", ip)

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			// Fail open but log it if Redis is down
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, window)
		}

		if int(count) > maxAttempts {
			response.Error(c, 429, "TOO_MANY_REQUESTS", "Terlalu banyak percobaan login. Silakan coba lagi dalam 15 menit.")
			c.Abort()
			return
		}

		c.Next()
	}
}
