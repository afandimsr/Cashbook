package middleware

import (
	"strings"

	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(apperror.Unauthorized("Access Denied, Login Required", nil))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(apperror.Unauthorized("invalid authorization format", nil))
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(parts[1])
		if err != nil {
			c.Error(apperror.Unauthorized("invalid or expired token", err))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		// set roles (array) and a primary role for backward compatibility
		if len(claims.Roles) > 0 {
			c.Set("roles", claims.Roles)
			c.Set("role", claims.Roles[0])
		}
		c.Next()
	}
}
