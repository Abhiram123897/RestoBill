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

// RestaurantTableRepository defines the contract for restaurant table database operations
type RestaurantTableRepository interface {
	Create(ctx context.Context, table *domain.RestaurantTable) (*domain.RestaurantTable, *errs.AppError)
	GetAll(ctx context.Context) ([]domain.RestaurantTable, *errs.AppError)
	GetByID(ctx context.Context, id string) (*domain.RestaurantTable, *errs.AppError)
	Update(ctx context.Context, table *domain.RestaurantTable) (*domain.RestaurantTable, *errs.AppError)
	UpdateStatus(ctx context.Context, id string, status string) (*domain.RestaurantTable, *errs.AppError)
	Delete(ctx context.Context, id string) *errs.AppError
}

// OrderRepository defines the contract for order database operations
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) (*domain.Order, *errs.AppError)
	GetAll(ctx context.Context, status string, orderType string) ([]domain.Order, *errs.AppError)
	GetByID(ctx context.Context, id string) (*domain.Order, *errs.AppError)
	GetActiveOrderByTableID(ctx context.Context, tableID string) (*domain.Order, *errs.AppError)
	UpdateStatus(ctx context.Context, id string, status string) (*domain.Order, *errs.AppError)
	AddItem(ctx context.Context, item *domain.OrderItem) (*domain.OrderItem, *errs.AppError)
	GetItemsByOrderID(ctx context.Context, orderID string) ([]domain.OrderItem, *errs.AppError)
	RemoveItem(ctx context.Context, orderID string, itemID string) *errs.AppError
}

// BillRepository defines the contract for bill database operations
type BillRepository interface {
	Create(ctx context.Context, bill *domain.Bill) (*domain.Bill, *errs.AppError)
	GetByID(ctx context.Context, id string) (*domain.Bill, *errs.AppError)
	GetByOrderID(ctx context.Context, orderID string) (*domain.Bill, *errs.AppError)
	Update(ctx context.Context, bill *domain.Bill) (*domain.Bill, *errs.AppError)
	UpdateStatus(ctx context.Context, id string, status string) (*domain.Bill, *errs.AppError)
}

// PaymentRepository defines the contract for payment database operations
type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, *errs.AppError)
	GetByBillID(ctx context.Context, billID string) ([]domain.Payment, *errs.AppError)
}

// UserRepository defines the contract for user database operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, *errs.AppError)
	GetByEmail(ctx context.Context, email string) (*domain.User, *errs.AppError)
	GetByID(ctx context.Context, id string) (*domain.User, *errs.AppError)
	GetUserCount(ctx context.Context) (int, *errs.AppError)
}
