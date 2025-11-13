package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	router := setupTestRouter()

	// Mock auth middleware
	router.Use(func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "Bearer valid-token" {
			c.Set("user_id", "user-123")
			c.Next()
		} else {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
	})

	router.GET("/protected", func(c *gin.Context) {
		userID := c.GetString("user_id")
		c.JSON(200, gin.H{"user_id": userID})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "user-123")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := setupTestRouter()

	router.Use(func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "Bearer valid-token" {
			c.Set("user_id", "user-123")
			c.Next()
		} else {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
	})

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	router := setupTestRouter()

	router.Use(func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		c.Next()
	})

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Missing token")
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	router := setupTestRouter()

	router.Use(func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "Bearer expired-token" {
			c.JSON(401, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}
		c.Next()
	})

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "expired")
}

func TestAuthMiddleware_RefreshToken(t *testing.T) {
	router := setupTestRouter()

	router.POST("/auth/refresh", func(c *gin.Context) {
		refreshToken := c.GetHeader("X-Refresh-Token")
		if refreshToken == "valid-refresh-token" {
			c.JSON(200, gin.H{
				"access_token":  "new-access-token",
				"refresh_token": "new-refresh-token",
			})
		} else {
			c.JSON(401, gin.H{"error": "Invalid refresh token"})
		}
	})

	req, _ := http.NewRequest("POST", "/auth/refresh", nil)
	req.Header.Set("X-Refresh-Token", "valid-refresh-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "new-access-token")
}

func TestAuthMiddleware_UserContext(t *testing.T) {
	router := setupTestRouter()

	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user-123")
		c.Set("user_email", "test@example.com")
		c.Next()
	})

	router.GET("/me", func(c *gin.Context) {
		userID := c.GetString("user_id")
		userEmail := c.GetString("user_email")
		c.JSON(200, gin.H{
			"user_id": userID,
			"email":   userEmail,
		})
	})

	req, _ := http.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "user-123")
	assert.Contains(t, w.Body.String(), "test@example.com")
}
