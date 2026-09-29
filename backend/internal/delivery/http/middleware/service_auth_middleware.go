package middleware

import (
	"crypto/subtle"

	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/gin-gonic/gin"
)

// ServiceAuthMiddleware protects internal service-to-service routes (e.g. the
// Telegram bot) with a shared secret, instead of the end-user JWT flow used by
// AuthMiddleware. It never sets a user_id in context — every internal/bot
// handler resolves identity itself from the caller-supplied telegram_chat_id,
// so a leaked service key can only act as accounts that were explicitly
// linked, never arbitrary users.
func ServiceAuthMiddleware(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		provided := c.GetHeader("X-Internal-Api-Key")
		if expectedKey == "" || provided == "" ||
			subtle.ConstantTimeCompare([]byte(provided), []byte(expectedKey)) != 1 {
			c.Error(apperror.Unauthorized("invalid or missing internal api key", nil))
			c.Abort()
			return
		}
		c.Next()
	}
}
