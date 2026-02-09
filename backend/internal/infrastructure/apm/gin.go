package apm

import (
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
)

func GinMiddleware() gin.HandlerFunc {
	return apmgin.Middleware(nil)
}
