package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"restobill/internal/db"
	"restobill/internal/domain"
	"restobill/internal/errs"
	"restobill/internal/logger"
)

// menuItemRepo implements MenuItemRepository using PostgreSQL
type menuItemRepo struct {
	pool *pgxpool.Pool
}

// NewMenuItemRepository creates a new MenuItemRepository instance
func NewMenuItemRepository(pool *pgxpool.Pool) MenuItemRepository {
	return &menuItemRepo{pool: pool}
}

// Create inserts a new menu item into the database
func (r *menuItemRepo) Create(ctx context.Context, item *domain.MenuItem) (*domain.MenuItem, *errs.AppError) {
	query := `
		INSERT INTO menu_items (category_id, name, price, is_available)
		VALUES ($1, $2, $3, $4)
		RETURNING id, category_id, name, price, is_available, created_at, updated_at`

	var m domain.MenuItem
	err := r.pool.QueryRow(ctx, query, item.CategoryID, item.Name, item.Price, item.IsAvailable).Scan(
		&m.ID, &m.CategoryID, &m.Name, &m.Price, &m.IsAvailable, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsForeignKeyViolation(pgErr) {
				return nil, errs.NewBadRequestError("category does not exist")
			}
			if db.IsUniqueViolation(pgErr) {
				return nil, errs.NewConflictError("menu item with this name already exists in this category")
			}
		}
		logger.ErrorF("failed to create menu item: %v", err)
		return nil, errs.NewInternalServerError("failed to create menu item")
	}

	return &m, nil
}

// GetAll retrieves all menu items
func (r *menuItemRepo) GetAll(ctx context.Context) ([]domain.MenuItem, *errs.AppError) {
	query := `SELECT id, category_id, name, price, is_available, created_at, updated_at FROM menu_items ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		logger.ErrorF("failed to get menu items: %v", err)
		return nil, errs.NewInternalServerError("failed to get menu items")
	}
	defer rows.Close()

	var items []domain.MenuItem
	for rows.Next() {
		var m domain.MenuItem
		if err := rows.Scan(&m.ID, &m.CategoryID, &m.Name, &m.Price, &m.IsAvailable, &m.CreatedAt, &m.UpdatedAt); err != nil {
			logger.ErrorF("failed to scan menu item: %v", err)
			return nil, errs.NewInternalServerError("failed to scan menu item")
		}
		items = append(items, m)
	}

	return items, nil
}

// GetByID retrieves a single menu item by its ID
func (r *menuItemRepo) GetByID(ctx context.Context, id string) (*domain.MenuItem, *errs.AppError) {
	query := `SELECT id, category_id, name, price, is_available, created_at, updated_at FROM menu_items WHERE id = $1`

	var m domain.MenuItem
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.CategoryID, &m.Name, &m.Price, &m.IsAvailable, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		logger.ErrorF("failed to get menu item by id: %v", err)
		return nil, errs.NewNotFoundError("menu item not found")
	}

	return &m, nil
}

// GetByCategory retrieves all menu items for a specific category
func (r *menuItemRepo) GetByCategory(ctx context.Context, categoryID string) ([]domain.MenuItem, *errs.AppError) {
	query := `SELECT id, category_id, name, price, is_available, created_at, updated_at FROM menu_items WHERE category_id = $1 ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query, categoryID)
	if err != nil {
		logger.ErrorF("failed to get menu items by category: %v", err)
		return nil, errs.NewInternalServerError("failed to get menu items by category")
	}
	defer rows.Close()

	var items []domain.MenuItem
	for rows.Next() {
		var m domain.MenuItem
		if err := rows.Scan(&m.ID, &m.CategoryID, &m.Name, &m.Price, &m.IsAvailable, &m.CreatedAt, &m.UpdatedAt); err != nil {
			logger.ErrorF("failed to scan menu item: %v", err)
			return nil, errs.NewInternalServerError("failed to scan menu item")
		}
		items = append(items, m)
	}

	return items, nil
}

// Update modifies an existing menu item
func (r *menuItemRepo) Update(ctx context.Context, item *domain.MenuItem) (*domain.MenuItem, *errs.AppError) {
	query := `
		UPDATE menu_items
		SET category_id = $1, name = $2, price = $3, is_available = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, category_id, name, price, is_available, created_at, updated_at`

	var m domain.MenuItem
	err := r.pool.QueryRow(ctx, query, item.CategoryID, item.Name, item.Price, item.IsAvailable, item.ID).Scan(
		&m.ID, &m.CategoryID, &m.Name, &m.Price, &m.IsAvailable, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsForeignKeyViolation(pgErr) {
				return nil, errs.NewBadRequestError("category does not exist")
			}
			if db.IsUniqueViolation(pgErr) {
				return nil, errs.NewConflictError("menu item with this name already exists in this category")
			}
		}
		logger.ErrorF("failed to update menu item: %v", err)
		return nil, errs.NewNotFoundError("menu item not found")
	}

	return &m, nil
}

// UpdateAvailability toggles the availability of a menu item
func (r *menuItemRepo) UpdateAvailability(ctx context.Context, id string, isAvailable bool) (*domain.MenuItem, *errs.AppError) {
	query := `
		UPDATE menu_items
		SET is_available = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, category_id, name, price, is_available, created_at, updated_at`

	var m domain.MenuItem
	err := r.pool.QueryRow(ctx, query, isAvailable, id).Scan(
		&m.ID, &m.CategoryID, &m.Name, &m.Price, &m.IsAvailable, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		logger.ErrorF("failed to update menu item availability: %v", err)
		return nil, errs.NewNotFoundError("menu item not found")
	}

	return &m, nil
}

// Delete removes a menu item by its ID
func (r *menuItemRepo) Delete(ctx context.Context, id string) *errs.AppError {
	query := `DELETE FROM menu_items WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		logger.ErrorF("failed to delete menu item: %v", err)
		return errs.NewInternalServerError("failed to delete menu item")
	}

	if result.RowsAffected() == 0 {
		return errs.NewNotFoundError("menu item not found")
	}

	return nil
}
