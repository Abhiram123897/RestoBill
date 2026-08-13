package service

import (
	"context"

	"restobill/internal/domain"
	"restobill/internal/dto"
	"restobill/internal/errs"
	"restobill/internal/logger"
	"restobill/internal/repository"
)

// CategoryService defines the interface for category business logic
type CategoryService interface {
	CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, *errs.AppError)
	GetAllCategories(ctx context.Context) ([]dto.CategoryResponse, *errs.AppError)
	GetCategoryByID(ctx context.Context, id string) (*dto.CategoryResponse, *errs.AppError)
	UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, *errs.AppError)
	DeleteCategory(ctx context.Context, id string) *errs.AppError
}

// categoryService implements CategoryService
type categoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new CategoryService instance
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// CreateCategory handles the business logic for creating a category
func (s *categoryService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, *errs.AppError) {
	category := &domain.Category{
		Name:         req.Name,
		DisplayOrder: req.DisplayOrder,
	}

	created, appErr := s.repo.Create(ctx, category)
	if appErr != nil {
		logger.ErrorF("service: failed to create category: %v", appErr.Message)
		return nil, appErr
	}

	return toCategoryResponse(created), nil
}

// GetAllCategories retrieves all categories
func (s *categoryService) GetAllCategories(ctx context.Context) ([]dto.CategoryResponse, *errs.AppError) {
	categories, appErr := s.repo.GetAll(ctx)
	if appErr != nil {
		logger.ErrorF("service: failed to get categories: %v", appErr.Message)
		return nil, appErr
	}

	response := make([]dto.CategoryResponse, 0)
	for _, c := range categories {
		response = append(response, *toCategoryResponse(&c))
	}

	return response, nil
}

// GetCategoryByID retrieves a single category by ID
func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*dto.CategoryResponse, *errs.AppError) {
	category, appErr := s.repo.GetByID(ctx, id)
	if appErr != nil {
		logger.ErrorF("service: failed to get category: %v", appErr.Message)
		return nil, appErr
	}

	return toCategoryResponse(category), nil
}

// UpdateCategory handles the business logic for updating a category
func (s *categoryService) UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, *errs.AppError) {
	category := &domain.Category{
		ID:           id,
		Name:         req.Name,
		DisplayOrder: req.DisplayOrder,
	}

	updated, appErr := s.repo.Update(ctx, category)
	if appErr != nil {
		logger.ErrorF("service: failed to update category: %v", appErr.Message)
		return nil, appErr
	}

	return toCategoryResponse(updated), nil
}

// DeleteCategory handles the business logic for deleting a category
func (s *categoryService) DeleteCategory(ctx context.Context, id string) *errs.AppError {
	appErr := s.repo.Delete(ctx, id)
	if appErr != nil {
		logger.ErrorF("service: failed to delete category: %v", appErr.Message)
		return appErr
	}

	return nil
}

// toCategoryResponse converts a domain Category to a CategoryResponse DTO
func toCategoryResponse(c *domain.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:           c.ID,
		Name:         c.Name,
		DisplayOrder: c.DisplayOrder,
		CreatedAt:    c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
