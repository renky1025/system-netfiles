package middleware

import (
	"netfilessys/internal/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RateLimiterMiddleware creates a rate limiter middleware
// Different endpoints have different rate limits:
// - Login: 5 req/min (prevent brute force)
// - Upload: 20 req/min (prevent abuse)
// - Default: 60 req/min (general API)
func RateLimiterMiddleware(rateLimit string) gin.HandlerFunc {
	var store limiter.Store
	var err error

	// Try to use Redis if available, fallback to memory
	if config.AppConfig.Redis.Addr != "" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     config.AppConfig.Redis.Addr,
			Password: config.AppConfig.Redis.Password,
			DB:       config.AppConfig.Redis.DB,
		})
		store, err = sredis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
			Prefix:   "rate_limit",
			MaxRetry: 3,
		})
		if err != nil {
			// Fallback to memory store
			store = memory.NewStore()
		}
	} else {
		// Use memory store if Redis not configured
		store = memory.NewStore()
	}

	// Parse rate limit string (e.g., "60-M" = 60 per minute)
	rate, err := limiter.NewRateFromFormatted(rateLimit)
	if err != nil {
		// Default to 60 per minute if parsing fails
		rate = limiter.Rate{
			Period: 1 * time.Minute,
			Limit:  60,
		}
	}

	// Create limiter instance
	instance := limiter.New(store, rate)

	// Return gin middleware
	return mgin.NewMiddleware(instance)
}

// LoginRateLimiter limits login attempts to prevent brute force attacks
func LoginRateLimiter() gin.HandlerFunc {
	return RateLimiterMiddleware("5-M") // 5 requests per minute
}

// UploadRateLimiter limits upload requests to prevent abuse
func UploadRateLimiter() gin.HandlerFunc {
	return RateLimiterMiddleware("20-M") // 20 requests per minute
}

// DefaultRateLimiter is the default rate limiter for general API endpoints
func DefaultRateLimiter() gin.HandlerFunc {
	return RateLimiterMiddleware("60-M") // 60 requests per minute
}

// GlobalRateLimiter limits all requests globally
func GlobalRateLimiter() gin.HandlerFunc {
	return RateLimiterMiddleware("100-M") // 100 requests per minute
}
