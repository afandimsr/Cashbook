package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

// newRouter returns a gin engine wired with the error handler and a stub that
// injects a user_id (as the AuthMiddleware would).
func newRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	r.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Set("roles", []string{"ADMIN", "USER"})
		c.Next()
	})
	return r
}

func doJSON(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
