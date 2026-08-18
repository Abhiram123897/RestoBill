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
	tableRepo := repository.NewRestaurantTableRepository(dbClient.Pool)
	orderRepo := repository.NewOrderRepository(dbClient.Pool)

	// ===== Initialize Services =====
	categoryService := service.NewCategoryService(categoryRepo)
	menuItemService := service.NewMenuItemService(menuItemRepo)
	tableService := service.NewRestaurantTableService(tableRepo)
	orderService := service.NewOrderService(orderRepo, menuItemRepo, tableRepo)

	// ===== Initialize Handlers =====
	healthHandler := handler.NewHealthHandler(dbClient)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	menuItemHandler := handler.NewMenuItemHandler(menuItemService)
	tableHandler := handler.NewRestaurantTableHandler(tableService)
	orderHandler := handler.NewOrderHandler(orderService)

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

		// Restaurant Tables
		api.POST("/tables", tableHandler.CreateTable)
		api.GET("/tables", tableHandler.GetAllTables)
		api.GET("/tables/:id", tableHandler.GetTableByID)
		api.PUT("/tables/:id", tableHandler.UpdateTable)
		api.PATCH("/tables/:id/status", tableHandler.UpdateTableStatus)
		api.DELETE("/tables/:id", tableHandler.DeleteTable)

		// Orders
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders", orderHandler.GetAllOrders)
		api.GET("/orders/:id", orderHandler.GetOrderByID)
		api.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)
		api.POST("/orders/:id/items", orderHandler.AddItemToOrder)
		api.DELETE("/orders/:id/items/:itemId", orderHandler.RemoveItemFromOrder)
	}

	return r
}
