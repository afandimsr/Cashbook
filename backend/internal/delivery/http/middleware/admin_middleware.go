package middleware

import (
	"net/http"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/gin-gonic/gin"
)

// AdminOnly ensures the user has ADMIN role
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// roles set by AuthMiddleware
		if rolesI, exists := c.Get("roles"); exists {
			if roles, ok := rolesI.([]string); ok {
				for _, r := range roles {
					if r == "ADMIN" {
						c.Next()
						return
					}
				}
			}
		}
		// fallback to single role
		if roleI, exists := c.Get("role"); exists {
			if roleStr, ok := roleI.(string); ok && roleStr == "ADMIN" {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "FORBIDDEN", "admin role required", "")
		c.Abort()
	}
}
