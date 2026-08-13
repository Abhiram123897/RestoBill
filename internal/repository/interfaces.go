package repository

import (
	"context"

	"restobill/internal/domain"
	"restobill/internal/errs"
)

// CategoryRepository defines the contract for category database operations
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) (*domain.Category, *errs.AppError)
	GetAll(ctx context.Context) ([]domain.Category, *errs.AppError)
	GetByID(ctx context.Context, id string) (*domain.Category, *errs.AppError)
	Update(ctx context.Context, category *domain.Category) (*domain.Category, *errs.AppError)
	Delete(ctx context.Context, id string) *errs.AppError
}

// MenuItemRepository defines the contract for menu item database operations
type MenuItemRepository interface {
	Create(ctx context.Context, item *domain.MenuItem) (*domain.MenuItem, *errs.AppError)
	GetAll(ctx context.Context) ([]domain.MenuItem, *errs.AppError)
	GetByID(ctx context.Context, id string) (*domain.MenuItem, *errs.AppError)
	GetByCategory(ctx context.Context, categoryID string) ([]domain.MenuItem, *errs.AppError)
	Update(ctx context.Context, item *domain.MenuItem) (*domain.MenuItem, *errs.AppError)
	UpdateAvailability(ctx context.Context, id string, isAvailable bool) (*domain.MenuItem, *errs.AppError)
	Delete(ctx context.Context, id string) *errs.AppError
}
