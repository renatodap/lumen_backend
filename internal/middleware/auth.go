package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/lumen/backend-go/pkg/errors"
	"github.com/lumen/backend-go/pkg/logger"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appErr := apperrors.NewUnauthorized("missing authorization header")
			c.JSON(appErr.StatusCode, appErr)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			appErr := apperrors.NewUnauthorized("invalid authorization header format")
			c.JSON(appErr.StatusCode, appErr)
			c.Abort()
			return
		}

		token := parts[1]
		userID, err := m.validateToken(token)
		if err != nil {
			logger.Warn("Invalid token", zap.Error(err))
			appErr := apperrors.NewUnauthorized("invalid or expired token")
			c.JSON(appErr.StatusCode, appErr)
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("token", token)

		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				userID, err := m.validateToken(token)
				if err == nil {
					c.Set("user_id", userID)
					c.Set("token", token)
				}
			}
		}

		c.Next()
	}
}

func (m *AuthMiddleware) validateToken(token string) (uuid.UUID, error) {
	userIDStr := token
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			appErr := apperrors.NewForbidden("user role not found")
			c.JSON(appErr.StatusCode, appErr)
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			appErr := apperrors.NewForbidden("invalid user role")
			c.JSON(appErr.StatusCode, appErr)
			c.Abort()
			return
		}

		hasRole := false
		for _, r := range roles {
			if role == r {
				hasRole = true
				break
			}
		}

		if !hasRole {
			appErr := apperrors.NewForbidden("insufficient permissions")
			c.JSON(appErr.StatusCode, appErr)
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, apperrors.NewUnauthorized("user not authenticated")
	}

	switch v := userIDVal.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		return uuid.Parse(v)
	default:
		return uuid.Nil, apperrors.NewUnauthorized("invalid user ID format")
	}
}
