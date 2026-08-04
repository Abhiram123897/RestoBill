package router

import (
	"github.com/gin-gonic/gin"
	"restobill/internal/db"
	"restobill/internal/handler"
)

// Setup creates the Gin engine with all routes
func Setup(dbClient *db.Client) *gin.Engine {
	r := gin.Default()

	// Handlers
	healthHandler := handler.NewHealthHandler(dbClient)

	// Routes
	api := r.Group("/api/v1")
	{
		api.GET("/health", healthHandler.GetHealth)
	}

	return r
}
