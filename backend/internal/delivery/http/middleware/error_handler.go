package middleware

import (
	"errors"
	"net/http"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			response.Error(
				c,
				appErr.Code,
				appErr.ErrorCode,
				appErr.Message,
				appErr.Error(),
			)
			return
		}

		// fallback unknown error
		response.Error(
			c,
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"internal server error",
			err.Error(),
		)
	}
}
