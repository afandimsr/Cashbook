package middleware

import (
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/pkg/ratelimit"
	"github.com/gin-gonic/gin"
)

// RateLimit throttles requests per key (typically client IP), independent of
// whether the request ultimately succeeds — a cheap first line of defense
// against broad scripted abuse on public-ish endpoints. For lockouts that
// should only count failed attempts (e.g. login), don't use this middleware —
// call limiter.Blocked/RecordFailure directly in the handler instead.
func RateLimit(limiter *ratelimit.Limiter, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow(keyFunc(c)) {
			c.Error(apperror.TooManyRequests("too many requests, please try again later", nil))
			c.Abort()
			return
		}
		c.Next()
	}
}

// ClientIPKey is the standard keyFunc for per-IP rate limiting.
func ClientIPKey(c *gin.Context) string {
	return c.ClientIP()
}
