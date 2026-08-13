package router

import (
	"github.com/gin-gonic/gin"
	"restobill/internal/db"
	"restobill/internal/handler"
	"restobill/internal/repository"
	"restobill/internal/service"
)

// Setup creates the Gin engine with all routes
func Setup(dbClient *db.Client) *gin.Engine {
	r := gin.Default()

	// ===== Initialize Repositories =====
	categoryRepo := repository.NewCategoryRepository(dbClient.Pool)
	menuItemRepo := repository.NewMenuItemRepository(dbClient.Pool)

	// ===== Initialize Services =====
	categoryService := service.NewCategoryService(categoryRepo)
	menuItemService := service.NewMenuItemService(menuItemRepo)

	// ===== Initialize Handlers =====
	healthHandler := handler.NewHealthHandler(dbClient)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	menuItemHandler := handler.NewMenuItemHandler(menuItemService)

	// ===== Routes =====
	api := r.Group("/api/v1")
	{
		// Health
		api.GET("/health", healthHandler.GetHealth)

		// Categories
		api.POST("/categories", categoryHandler.CreateCategory)
		api.GET("/categories", categoryHandler.GetAllCategories)
		api.GET("/categories/:id", categoryHandler.GetCategoryByID)
		api.PUT("/categories/:id", categoryHandler.UpdateCategory)
		api.DELETE("/categories/:id", categoryHandler.DeleteCategory)

		// Menu Items
		api.POST("/menu-items", menuItemHandler.CreateMenuItem)
		api.GET("/menu-items", menuItemHandler.GetAllMenuItems)
		api.GET("/menu-items/:id", menuItemHandler.GetMenuItemByID)
		api.PUT("/menu-items/:id", menuItemHandler.UpdateMenuItem)
		api.PATCH("/menu-items/:id/availability", menuItemHandler.UpdateAvailability)
		api.DELETE("/menu-items/:id", menuItemHandler.DeleteMenuItem)
	}

	return r
}
