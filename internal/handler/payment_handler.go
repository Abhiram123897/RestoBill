package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"restobill/internal/dto"
	"restobill/internal/logger"
	"restobill/internal/service"
)

// PaymentHandler handles payment HTTP requests
type PaymentHandler struct {
	paymentService service.PaymentService
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// RecordPayment handles POST /api/v1/payments
func (h *PaymentHandler) RecordPayment(c *gin.Context) {
	var req dto.RecordPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	if _, err := uuid.Parse(req.BillID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid bill_id format", "type": "Bad request"})
		return
	}

	response, appErr := h.paymentService.RecordPayment(c.Request.Context(), &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetPaymentsByBillID handles GET /api/v1/bills/:id/payments
func (h *PaymentHandler) GetPaymentsByBillID(c *gin.Context) {
	billID := c.Param("id")
	if _, err := uuid.Parse(billID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid bill id format", "type": "Bad request"})
		return
	}

	response, appErr := h.paymentService.GetPaymentsByBillID(c.Request.Context(), billID)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}
