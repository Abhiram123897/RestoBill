package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"restobill/internal/db"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	dbClient *db.Client
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(dbClient *db.Client) *HealthHandler {
	return &HealthHandler{dbClient: dbClient}
}

// GetHealth checks if the server and database are healthy
func (h *HealthHandler) GetHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := h.dbClient.Ping(ctx)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "unhealthy",
			"service":  "restobill",
			"database": "disconnected",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"service":  "restobill",
		"database": "connected",
	})
}
