package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"restobill/internal/dto"
	"restobill/internal/logger"
	"restobill/internal/service"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// CreateCategory handles POST /api/v1/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	response, appErr := h.categoryService.CreateCategory(c.Request.Context(), &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetAllCategories handles GET /api/v1/categories
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	response, appErr := h.categoryService.GetAllCategories(c.Request.Context())
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCategoryByID handles GET /api/v1/categories/:id
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	response, appErr := h.categoryService.GetCategoryByID(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCategory handles PUT /api/v1/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorF("handler: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "type": "Bad request"})
		return
	}

	response, appErr := h.categoryService.UpdateCategory(c.Request.Context(), id, &req)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteCategory handles DELETE /api/v1/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format", "type": "Bad request"})
		return
	}

	appErr := h.categoryService.DeleteCategory(c.Request.Context(), id)
	if appErr != nil {
		c.JSON(appErr.Code, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}
