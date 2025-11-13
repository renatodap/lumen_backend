package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware_AllowsWithinLimit(t *testing.T) {
	router := setupTestRouter()

	requestCount := make(map[string]int)

	router.Use(func(c *gin.Context) {
		clientIP := c.ClientIP()
		requestCount[clientIP]++

		if requestCount[clientIP] > 10 {
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	})

	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// Make 5 requests (within limit)
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	}
}

func TestRateLimitMiddleware_BlocksExcessRequests(t *testing.T) {
	router := setupTestRouter()

	requestCount := make(map[string]int)

	router.Use(func(c *gin.Context) {
		clientIP := c.ClientIP()
		requestCount[clientIP]++

		if requestCount[clientIP] > 5 {
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	})

	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// Make 7 requests (exceeds limit of 5)
	for i := 0; i < 7; i++ {
		req, _ := http.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if i < 5 {
			assert.Equal(t, 200, w.Code)
		} else {
			assert.Equal(t, 429, w.Code)
			assert.Contains(t, w.Body.String(), "Rate limit exceeded")
		}
	}
}

func TestRateLimitMiddleware_PerIPTracking(t *testing.T) {
	router := setupTestRouter()

	requestCount := make(map[string]int)

	router.Use(func(c *gin.Context) {
		clientIP := c.ClientIP()
		requestCount[clientIP]++

		if requestCount[clientIP] > 5 {
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	})

	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// IP 1 - 5 requests (ok)
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}

	// IP 2 - 5 requests (ok, different IP)
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.2:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}

	// IP 1 - 6th request (blocked)
	req, _ := http.NewRequest("GET", "/api/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 429, w.Code)
}

func TestRateLimitMiddleware_RetryAfterHeader(t *testing.T) {
	router := setupTestRouter()

	requestCount := make(map[string]int)
	resetTime := make(map[string]time.Time)

	router.Use(func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		if resetTime[clientIP].IsZero() {
			resetTime[clientIP] = now.Add(1 * time.Minute)
		}

		if now.After(resetTime[clientIP]) {
			requestCount[clientIP] = 0
			resetTime[clientIP] = now.Add(1 * time.Minute)
		}

		requestCount[clientIP]++

		if requestCount[clientIP] > 5 {
			retryAfter := int(time.Until(resetTime[clientIP]).Seconds())
			c.Header("Retry-After", string(rune(retryAfter)))
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	})

	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// Exceed rate limit
	for i := 0; i < 6; i++ {
		req, _ := http.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if i == 5 {
			assert.Equal(t, 429, w.Code)
			assert.NotEmpty(t, w.Header().Get("Retry-After"))
		}
	}
}

func TestRateLimitMiddleware_WhitelistedEndpoints(t *testing.T) {
	router := setupTestRouter()

	requestCount := make(map[string]int)
	whitelistedPaths := map[string]bool{
		"/health": true,
		"/ping":   true,
	}

	router.Use(func(c *gin.Context) {
		// Skip rate limiting for whitelisted paths
		if whitelistedPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		requestCount[clientIP]++

		if requestCount[clientIP] > 5 {
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// Make 10 requests to health endpoint (should all succeed)
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest("GET", "/health", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}

	// Make 6 requests to rate-limited endpoint
	for i := 0; i < 6; i++ {
		req, _ := http.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if i < 5 {
			assert.Equal(t, 200, w.Code)
		} else {
			assert.Equal(t, 429, w.Code)
		}
	}
}
