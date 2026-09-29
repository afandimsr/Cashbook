package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServiceAuthMiddleware_ValidKey(t *testing.T) {
	r := gin.New()
	r.Use(middleware.ServiceAuthMiddleware("secret-key"))
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Internal-Api-Key", "secret-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestServiceAuthMiddleware_WrongKey(t *testing.T) {
	r := gin.New()
	r.Use(middleware.ErrorHandler(), middleware.ServiceAuthMiddleware("secret-key"))
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Internal-Api-Key", "wrong-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestServiceAuthMiddleware_MissingHeader(t *testing.T) {
	r := gin.New()
	r.Use(middleware.ErrorHandler(), middleware.ServiceAuthMiddleware("secret-key"))
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestServiceAuthMiddleware_EmptyConfiguredKeyDeniesAll(t *testing.T) {
	r := gin.New()
	r.Use(middleware.ErrorHandler(), middleware.ServiceAuthMiddleware(""))
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Internal-Api-Key", "")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
