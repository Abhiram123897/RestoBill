package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"restobill/internal/dto"
	"restobill/internal/logger"
	"restobill/internal/service"
)

// BillHandler handles bill HTTP requests
type BillHandler struct {
	billService service.BillService
}

// NewBillHandler creates a new BillHandler
func NewBillHandler(billService service.BillService) *BillHandler {
	return &BillHandler{billService: billService}
}

// GenerateBill handles POST /api/v1/bills
func (h *BillHandler) GenerateBill(c *gin.Context) {
	var req dto.GenerateBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	if _, err := uuid.Parse(req.OrderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid order_id format", "type": "Bad request"})
		return
	}

	if _, err := uuid.Parse(req.GeneratedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid generated_by format", "type": "Bad request"})
		return
	}

	response, appErr := h.billService.GenerateBill(c.Request.Context(), &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetBillByID handles GET /api/v1/bills/:id
func (h *BillHandler) GetBillByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	response, appErr := h.billService.GetBillByID(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetBillByOrderID handles GET /api/v1/orders/:id/bill
func (h *BillHandler) GetBillByOrderID(c *gin.Context) {
	orderID := c.Param("id")
	if _, err := uuid.Parse(orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid order id format", "type": "Bad request"})
		return
	}

	response, appErr := h.billService.GetBillByOrderID(c.Request.Context(), orderID)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ApplyDiscount handles PATCH /api/v1/bills/:id/discount
func (h *BillHandler) ApplyDiscount(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	var req dto.ApplyDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	response, appErr := h.billService.ApplyDiscount(c.Request.Context(), id, &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// VoidBill handles PATCH /api/v1/bills/:id/void
func (h *BillHandler) VoidBill(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	response, appErr := h.billService.VoidBill(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}
