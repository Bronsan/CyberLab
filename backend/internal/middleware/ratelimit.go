package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter provides per-IP rate limiting using a sliding window.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

type visitor struct {
	count    int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter.
// limit: max requests per window per IP.
// window: time window (e.g., 1 minute).
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
	// Background cleanup every 10 minutes
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window*2 {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if the given IP is allowed to proceed.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[ip]

	if !exists {
		rl.visitors[ip] = &visitor{count: 1, lastSeen: now}
		return true
	}

	// Reset if window has passed
	if now.Sub(v.lastSeen) > rl.window {
		v.count = 1
		v.lastSeen = now
		return true
	}

	v.lastSeen = now
	v.count++
	return v.count <= rl.limit
}

// RateLimitMiddleware returns a Gin middleware for rate limiting.
// Returns 429 Too Many Requests if limit exceeded.
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Default login rate limiter: 10 requests per minute per IP
var LoginLimiter = NewRateLimiter(10, 1*time.Minute)

// Default register rate limiter: 5 requests per minute per IP
var RegisterLimiter = NewRateLimiter(5, 1*time.Minute)
