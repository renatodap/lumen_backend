package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lumen/backend-go/internal/repository"
)

type HealthHandler struct {
	db *repository.Database
}

func NewHealthHandler(db *repository.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx := c.Request.Context()

	dbStatus := "healthy"
	if err := h.db.Health(ctx); err != nil {
		dbStatus = "unhealthy"
	}

	response := gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "lumen-api",
		"version":   "1.0.0",
		"database":  dbStatus,
	}

	if dbStatus == "unhealthy" {
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.db.Health(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"message": "database connection unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func (h *HealthHandler) Metrics(c *gin.Context) {
	stats := h.db.Stats()
	c.JSON(http.StatusOK, gin.H{
		"database": stats,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
