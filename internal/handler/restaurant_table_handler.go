package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"restobill/internal/dto"
	"restobill/internal/logger"
	"restobill/internal/service"
)

// RestaurantTableHandler handles restaurant table HTTP requests
type RestaurantTableHandler struct {
	tableService service.RestaurantTableService
}

// NewRestaurantTableHandler creates a new RestaurantTableHandler
func NewRestaurantTableHandler(tableService service.RestaurantTableService) *RestaurantTableHandler {
	return &RestaurantTableHandler{tableService: tableService}
}

// CreateTable handles POST /api/v1/tables
func (h *RestaurantTableHandler) CreateTable(c *gin.Context) {
	var req dto.CreateRestaurantTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	response, appErr := h.tableService.CreateTable(c.Request.Context(), &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetAllTables handles GET /api/v1/tables
func (h *RestaurantTableHandler) GetAllTables(c *gin.Context) {
	response, appErr := h.tableService.GetAllTables(c.Request.Context())
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTableByID handles GET /api/v1/tables/:id
func (h *RestaurantTableHandler) GetTableByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	response, appErr := h.tableService.GetTableByID(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTable handles PUT /api/v1/tables/:id
func (h *RestaurantTableHandler) UpdateTable(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	var req dto.UpdateRestaurantTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	response, appErr := h.tableService.UpdateTable(c.Request.Context(), id, &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTableStatus handles PATCH /api/v1/tables/:id/status
func (h *RestaurantTableHandler) UpdateTableStatus(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	var req dto.UpdateTableStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	response, appErr := h.tableService.UpdateTableStatus(c.Request.Context(), id, &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTable handles DELETE /api/v1/tables/:id
func (h *RestaurantTableHandler) DeleteTable(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	appErr := h.tableService.DeleteTable(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "table deleted successfully"})
}
