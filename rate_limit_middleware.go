// login-plugin/rate_limit_middleware.go
package login

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	rate     rate.Limit
	burst    int
	duration time.Duration
}

func NewRateLimiter(r rate.Limit, b int, d time.Duration) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
		duration: d,
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
		go func() {
			time.Sleep(rl.duration)
			delete(rl.limiters, key)
		}()
	}
	return limiter
}

func (rl *RateLimiter) LoginLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		email := ""

		if c.Request.Method == http.MethodPost && strings.Contains(c.FullPath(), "/login") {
			var body struct {
				Email string `json:"email"`
			}
			if err := c.ShouldBindJSON(&body); err == nil {
				email = strings.ToLower(body.Email)
			}
		}

		keys := []string{fmt.Sprintf("ip:%s", ip)}
		if email != "" {
			keys = append(keys, fmt.Sprintf("email:%s", email))
		}

		for _, key := range keys {
			limiter := rl.getLimiter(key)
			if !limiter.Allow() {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Tentativas excedidas, tente novamente em instantes"})
				return
			}
		}

		c.Next()
	}
}
