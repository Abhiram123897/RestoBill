package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"restobill/internal/dto"
	"restobill/internal/logger"
	"restobill/internal/service"
)

// MenuItemHandler handles menu item HTTP requests
type MenuItemHandler struct {
	menuItemService service.MenuItemService
}

// NewMenuItemHandler creates a new MenuItemHandler
func NewMenuItemHandler(menuItemService service.MenuItemService) *MenuItemHandler {
	return &MenuItemHandler{menuItemService: menuItemService}
}

// CreateMenuItem handles POST /api/v1/menu-items
func (h *MenuItemHandler) CreateMenuItem(c *gin.Context) {
	var req dto.CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	if _, err := uuid.Parse(req.CategoryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid category_id format", "type": "Bad request"})
		return
	}

	response, appErr := h.menuItemService.CreateMenuItem(c.Request.Context(), &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetAllMenuItems handles GET /api/v1/menu-items
func (h *MenuItemHandler) GetAllMenuItems(c *gin.Context) {
	categoryID := c.Query("category_id")
	if categoryID != "" {
		if _, err := uuid.Parse(categoryID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid category_id format", "type": "Bad request"})
			return
		}
		response, appErr := h.menuItemService.GetMenuItemsByCategory(c.Request.Context(), categoryID)
		if appErr != nil {
			c.JSON(appErr.Code, appErr)
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}

	response, appErr := h.menuItemService.GetAllMenuItems(c.Request.Context())
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMenuItemByID handles GET /api/v1/menu-items/:id
func (h *MenuItemHandler) GetMenuItemByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	response, appErr := h.menuItemService.GetMenuItemByID(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateMenuItem handles PUT /api/v1/menu-items/:id
func (h *MenuItemHandler) UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	var req dto.UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	if _, err := uuid.Parse(req.CategoryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid category_id format", "type": "Bad request"})
		return
	}

	response, appErr := h.menuItemService.UpdateMenuItem(c.Request.Context(), id, &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateAvailability handles PATCH /api/v1/menu-items/:id/availability
func (h *MenuItemHandler) UpdateAvailability(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	var req struct {
		IsAvailable bool `json:"is_available"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	response, appErr := h.menuItemService.UpdateAvailability(c.Request.Context(), id, req.IsAvailable)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteMenuItem handles DELETE /api/v1/menu-items/:id
func (h *MenuItemHandler) DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	appErr := h.menuItemService.DeleteMenuItem(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "menu item deleted successfully"})
}
