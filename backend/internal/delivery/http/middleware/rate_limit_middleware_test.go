package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/middleware"
	"github.com/afandimsr/cashbook-backend/internal/pkg/ratelimit"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimit_AllowsUpToMax(t *testing.T) {
	limiter := ratelimit.New(2, time.Minute)
	r := gin.New()
	r.Use(middleware.ErrorHandler(), middleware.RateLimit(limiter, middleware.ClientIPKey))
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRateLimit_BlocksOverMax(t *testing.T) {
	limiter := ratelimit.New(2, time.Minute)
	r := gin.New()
	r.Use(middleware.ErrorHandler(), middleware.RateLimit(limiter, middleware.ClientIPKey))
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		r.ServeHTTP(httptest.NewRecorder(), req)
	}

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRateLimit_KeysIndependently(t *testing.T) {
	limiter := ratelimit.New(1, time.Minute)
	r := gin.New()
	r.Use(middleware.ErrorHandler(), middleware.RateLimit(limiter, func(c *gin.Context) string {
		return c.Query("key")
	}))
	r.GET("/x", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	req1 := httptest.NewRequest(http.MethodGet, "/x?key=a", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	req2 := httptest.NewRequest(http.MethodGet, "/x?key=b", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code, "a different key must have its own budget")
}
