package service

import (
	"context"

	"restobill/internal/domain"
	"restobill/internal/dto"
	"restobill/internal/errs"
	"restobill/internal/logger"
	"restobill/internal/repository"
)

// MenuItemService defines the interface for menu item business logic
type MenuItemService interface {
	CreateMenuItem(ctx context.Context, req *dto.CreateMenuItemRequest) (*dto.MenuItemResponse, *errs.AppError)
	GetAllMenuItems(ctx context.Context) ([]dto.MenuItemResponse, *errs.AppError)
	GetMenuItemByID(ctx context.Context, id string) (*dto.MenuItemResponse, *errs.AppError)
	GetMenuItemsByCategory(ctx context.Context, categoryID string) ([]dto.MenuItemResponse, *errs.AppError)
	UpdateMenuItem(ctx context.Context, id string, req *dto.UpdateMenuItemRequest) (*dto.MenuItemResponse, *errs.AppError)
	UpdateAvailability(ctx context.Context, id string, isAvailable bool) (*dto.MenuItemResponse, *errs.AppError)
	DeleteMenuItem(ctx context.Context, id string) *errs.AppError
}

// menuItemService implements MenuItemService
type menuItemService struct {
	repo repository.MenuItemRepository
}

// NewMenuItemService creates a new MenuItemService instance
func NewMenuItemService(repo repository.MenuItemRepository) MenuItemService {
	return &menuItemService{repo: repo}
}

// CreateMenuItem handles the business logic for creating a menu item
func (s *menuItemService) CreateMenuItem(ctx context.Context, req *dto.CreateMenuItemRequest) (*dto.MenuItemResponse, *errs.AppError) {
	isAvailable := true
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	item := &domain.MenuItem{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Price:       req.Price,
		IsAvailable: isAvailable,
	}

	created, appErr := s.repo.Create(ctx, item)
	if appErr != nil {
		logger.ErrorF("service: failed to create menu item: %v", appErr.Message)
		return nil, appErr
	}

	return toMenuItemResponse(created), nil
}

// GetAllMenuItems retrieves all menu items
func (s *menuItemService) GetAllMenuItems(ctx context.Context) ([]dto.MenuItemResponse, *errs.AppError) {
	items, appErr := s.repo.GetAll(ctx)
	if appErr != nil {
		logger.ErrorF("service: failed to get menu items: %v", appErr.Message)
		return nil, appErr
	}

	response := make([]dto.MenuItemResponse, 0)
	for _, item := range items {
		response = append(response, *toMenuItemResponse(&item))
	}

	return response, nil
}

// GetMenuItemByID retrieves a single menu item by ID
func (s *menuItemService) GetMenuItemByID(ctx context.Context, id string) (*dto.MenuItemResponse, *errs.AppError) {
	item, appErr := s.repo.GetByID(ctx, id)
	if appErr != nil {
		logger.ErrorF("service: failed to get menu item: %v", appErr.Message)
		return nil, appErr
	}

	return toMenuItemResponse(item), nil
}

// GetMenuItemsByCategory retrieves all menu items for a category
func (s *menuItemService) GetMenuItemsByCategory(ctx context.Context, categoryID string) ([]dto.MenuItemResponse, *errs.AppError) {
	items, appErr := s.repo.GetByCategory(ctx, categoryID)
	if appErr != nil {
		logger.ErrorF("service: failed to get menu items by category: %v", appErr.Message)
		return nil, appErr
	}

	response := make([]dto.MenuItemResponse, 0)
	for _, item := range items {
		response = append(response, *toMenuItemResponse(&item))
	}

	return response, nil
}

// UpdateMenuItem handles the business logic for updating a menu item
func (s *menuItemService) UpdateMenuItem(ctx context.Context, id string, req *dto.UpdateMenuItemRequest) (*dto.MenuItemResponse, *errs.AppError) {
	isAvailable := true
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	item := &domain.MenuItem{
		ID:          id,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Price:       req.Price,
		IsAvailable: isAvailable,
	}

	updated, appErr := s.repo.Update(ctx, item)
	if appErr != nil {
		logger.ErrorF("service: failed to update menu item: %v", appErr.Message)
		return nil, appErr
	}

	return toMenuItemResponse(updated), nil
}

// UpdateAvailability toggles menu item availability
func (s *menuItemService) UpdateAvailability(ctx context.Context, id string, isAvailable bool) (*dto.MenuItemResponse, *errs.AppError) {
	updated, appErr := s.repo.UpdateAvailability(ctx, id, isAvailable)
	if appErr != nil {
		logger.ErrorF("service: failed to update availability: %v", appErr.Message)
		return nil, appErr
	}

	return toMenuItemResponse(updated), nil
}

// DeleteMenuItem handles the business logic for deleting a menu item
func (s *menuItemService) DeleteMenuItem(ctx context.Context, id string) *errs.AppError {
	appErr := s.repo.Delete(ctx, id)
	if appErr != nil {
		logger.ErrorF("service: failed to delete menu item: %v", appErr.Message)
		return appErr
	}

	return nil
}

// toMenuItemResponse converts a domain MenuItem to a MenuItemResponse DTO
func toMenuItemResponse(m *domain.MenuItem) *dto.MenuItemResponse {
	return &dto.MenuItemResponse{
		ID:          m.ID,
		CategoryID:  m.CategoryID,
		Name:        m.Name,
		Price:       m.Price,
		IsAvailable: m.IsAvailable,
		CreatedAt:   m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
