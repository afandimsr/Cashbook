package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/middleware"
	"github.com/afandimsr/cashbook-backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() { gin.SetMode(gin.TestMode) }

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwt.SetSecret("test-secret")
	token, _ := jwt.GenerateToken(42, "u@e.com", "User", []string{"USER"})

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.GET("/x", func(c *gin.Context) {
		uid, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{"uid": uid})
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "42")
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	r := gin.New()
	r.Use(middleware.ErrorHandler(), middleware.AuthMiddleware())
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_BadFormat(t *testing.T) {
	r := gin.New()
	r.Use(middleware.ErrorHandler(), middleware.AuthMiddleware())
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Token abc")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwt.SetSecret("test-secret")
	r := gin.New()
	r.Use(middleware.ErrorHandler(), middleware.AuthMiddleware())
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer not-a-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminOnly_Allows(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("roles", []string{"USER", "ADMIN"}); c.Next() }, middleware.AdminOnly())
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminOnly_Forbids(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("roles", []string{"USER"}); c.Next() }, middleware.AdminOnly())
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
