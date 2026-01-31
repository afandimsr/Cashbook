package middleware

import (
	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/gin-gonic/gin"
)

// RoleGuard is a middleware that checks if the user has at least one of the allowed roles.
func RoleGuard(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, ok := c.Get("roles")
		if !ok {
			response.Error(c, 403, apperror.AuthForbidden, "Forbidden: you don't have access", nil)
			c.Abort()
			return
		}

		// Check if the user has at least one of the allowed roles
		userRoles := roles.([]string)
		for _, r := range userRoles {
			if contains(allowed, r) { // using a helper function to check if the allowed role is in user's roles
				c.Next()
				return
			}
		}

		// If we reach here, it means none of the user's roles were found in the allowed list. So we reject the request.
		response.Error(c, 403, apperror.AuthForbidden, "Forbidden: you don't have access", nil)
	}
}

// Helper function to check if a string is in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// how to use
