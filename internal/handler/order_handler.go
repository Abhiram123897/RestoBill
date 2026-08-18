package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"restobill/internal/dto"
	"restobill/internal/logger"
	"restobill/internal/service"
)

// OrderHandler handles order HTTP requests
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder handles POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	if req.TableID != nil {
		if _, err := uuid.Parse(*req.TableID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid table_id format", "type": "Bad request"})
			return
		}
	}

	if _, err := uuid.Parse(req.CreatedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid created_by format", "type": "Bad request"})
		return
	}

	response, appErr := h.orderService.CreateOrder(c.Request.Context(), &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetAllOrders handles GET /api/v1/orders
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	status := c.Query("status")
	orderType := c.Query("order_type")

	// Validate status filter if provided
	if status != "" {
		validStatuses := map[string]bool{"open": true, "served": true, "billed": true, "paid": true}
		if !validStatuses[status] {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid status filter: must be 'open', 'served', 'billed', or 'paid'", "type": "Bad request"})
			return
		}
	}

	// Validate order_type filter if provided
	if orderType != "" {
		validTypes := map[string]bool{"dine_in": true, "takeaway": true, "delivery": true}
		if !validTypes[orderType] {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid order_type filter: must be 'dine_in', 'takeaway', or 'delivery'", "type": "Bad request"})
			return
		}
	}

	response, appErr := h.orderService.GetAllOrders(c.Request.Context(), status, orderType)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetOrderByID handles GET /api/v1/orders/:id
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	response, appErr := h.orderService.GetOrderByID(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateOrderStatus handles PATCH /api/v1/orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	response, appErr := h.orderService.UpdateOrderStatus(c.Request.Context(), id, &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// AddItemToOrder handles POST /api/v1/orders/:id/items
func (h *OrderHandler) AddItemToOrder(c *gin.Context) {
	orderID := c.Param("id")
	if _, err := uuid.Parse(orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid order id format", "type": "Bad request"})
		return
	}

	var req dto.AddOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	if _, err := uuid.Parse(req.MenuItemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid menu_item_id format", "type": "Bad request"})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "quantity must be greater than 0", "type": "Bad request"})
		return
	}

	response, appErr := h.orderService.AddItemToOrder(c.Request.Context(), orderID, &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// RemoveItemFromOrder handles DELETE /api/v1/orders/:id/items/:itemId
func (h *OrderHandler) RemoveItemFromOrder(c *gin.Context) {
	orderID := c.Param("id")
	if _, err := uuid.Parse(orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid order id format", "type": "Bad request"})
		return
	}

	itemID := c.Param("itemId")
	if _, err := uuid.Parse(itemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid item id format", "type": "Bad request"})
		return
	}

	appErr := h.orderService.RemoveItemFromOrder(c.Request.Context(), orderID, itemID)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from order successfully"})
}
