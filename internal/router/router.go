package router

import (
	"github.com/gin-gonic/gin"
	"restobill/internal/config"
	"restobill/internal/db"
	"restobill/internal/handler"
	"restobill/internal/middleware"
	"restobill/internal/repository"
	"restobill/internal/service"
)

// Setup creates the Gin engine with all routes
func Setup(dbClient *db.Client, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// ===== Initialize Repositories =====
	categoryRepo := repository.NewCategoryRepository(dbClient.Pool)
	menuItemRepo := repository.NewMenuItemRepository(dbClient.Pool)
	tableRepo := repository.NewRestaurantTableRepository(dbClient.Pool)
	orderRepo := repository.NewOrderRepository(dbClient.Pool)
	billRepo := repository.NewBillRepository(dbClient.Pool)
	paymentRepo := repository.NewPaymentRepository(dbClient.Pool)
	userRepo := repository.NewUserRepository(dbClient.Pool)

	// ===== Initialize Services =====
	categoryService := service.NewCategoryService(categoryRepo)
	menuItemService := service.NewMenuItemService(menuItemRepo)
	tableService := service.NewRestaurantTableService(tableRepo)
	orderService := service.NewOrderService(orderRepo, menuItemRepo, tableRepo)
	billService := service.NewBillService(billRepo, orderRepo)
	paymentService := service.NewPaymentService(paymentRepo, billRepo, orderRepo, tableRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	// ===== Initialize Handlers =====
	healthHandler := handler.NewHealthHandler(dbClient)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	menuItemHandler := handler.NewMenuItemHandler(menuItemService)
	tableHandler := handler.NewRestaurantTableHandler(tableService)
	orderHandler := handler.NewOrderHandler(orderService)
	billHandler := handler.NewBillHandler(billService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	authHandler := handler.NewAuthHandler(authService)

	// ===== Routes =====
	api := r.Group("/api/v1")
	{
		// Public routes (no auth required)
		api.GET("/health", healthHandler.GetHealth)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", middleware.OptionalAuthMiddleware(cfg.JWTSecret), authHandler.Register)

		// Protected routes (auth required)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Auth
			protected.GET("/auth/me", authHandler.GetMe)

			// Categories (admin can create/update/delete, all can read)
			protected.GET("/categories", categoryHandler.GetAllCategories)
			protected.GET("/categories/:id", categoryHandler.GetCategoryByID)
			protected.POST("/categories", middleware.RoleMiddleware("admin"), categoryHandler.CreateCategory)
			protected.PUT("/categories/:id", middleware.RoleMiddleware("admin"), categoryHandler.UpdateCategory)
			protected.DELETE("/categories/:id", middleware.RoleMiddleware("admin"), categoryHandler.DeleteCategory)

			// Menu Items (admin can create/update/delete, all can read)
			protected.GET("/menu-items", menuItemHandler.GetAllMenuItems)
			protected.GET("/menu-items/:id", menuItemHandler.GetMenuItemByID)
			protected.POST("/menu-items", middleware.RoleMiddleware("admin"), menuItemHandler.CreateMenuItem)
			protected.PUT("/menu-items/:id", middleware.RoleMiddleware("admin"), menuItemHandler.UpdateMenuItem)
			protected.PATCH("/menu-items/:id/availability", middleware.RoleMiddleware("admin"), menuItemHandler.UpdateAvailability)
			protected.DELETE("/menu-items/:id", middleware.RoleMiddleware("admin"), menuItemHandler.DeleteMenuItem)

			// Restaurant Tables (admin can create/update/delete, all can read/update status)
			protected.GET("/tables", tableHandler.GetAllTables)
			protected.GET("/tables/:id", tableHandler.GetTableByID)
			protected.PATCH("/tables/:id/status", tableHandler.UpdateTableStatus)
			protected.POST("/tables", middleware.RoleMiddleware("admin"), tableHandler.CreateTable)
			protected.PUT("/tables/:id", middleware.RoleMiddleware("admin"), tableHandler.UpdateTable)
			protected.DELETE("/tables/:id", middleware.RoleMiddleware("admin"), tableHandler.DeleteTable)

			// Orders (all authenticated users can manage orders)
			protected.POST("/orders", orderHandler.CreateOrder)
			protected.GET("/orders", orderHandler.GetAllOrders)
			protected.GET("/orders/:id", orderHandler.GetOrderByID)
			protected.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)
			protected.POST("/orders/:id/items", orderHandler.AddItemToOrder)
			protected.DELETE("/orders/:id/items/:itemId", orderHandler.RemoveItemFromOrder)

			// Bills (cashier and admin can manage bills)
			protected.POST("/bills", middleware.RoleMiddleware("admin", "cashier"), billHandler.GenerateBill)
			protected.GET("/bills/:id", billHandler.GetBillByID)
			protected.GET("/orders/:id/bill", billHandler.GetBillByOrderID)
			protected.PATCH("/bills/:id/discount", middleware.RoleMiddleware("admin", "cashier"), billHandler.ApplyDiscount)
			protected.PATCH("/bills/:id/void", middleware.RoleMiddleware("admin"), billHandler.VoidBill)

			// Payments (cashier and admin can record payments)
			protected.POST("/payments", middleware.RoleMiddleware("admin", "cashier"), paymentHandler.RecordPayment)
			protected.GET("/bills/:id/payments", paymentHandler.GetPaymentsByBillID)
		}
	}

	return r
}
