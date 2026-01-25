package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/config"

	"github.com/gin-contrib/cors"
)

// Cors returns a CORS configuration based on the environment variables
func Cors(co *config.Config) cors.Config {
	origins := strings.Split(co.CorsAllowedOrigins, ",")

	// show origins for debugging
	log.Println("CORS Allowed Origins:", origins)

	return cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}
