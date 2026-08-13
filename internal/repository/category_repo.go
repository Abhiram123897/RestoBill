package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"restobill/internal/db"
	"restobill/internal/domain"
	"restobill/internal/errs"
	"restobill/internal/logger"
)

// categoryRepo implements CategoryRepository using PostgreSQL
type categoryRepo struct {
	pool *pgxpool.Pool
}

// NewCategoryRepository creates a new CategoryRepository instance
func NewCategoryRepository(pool *pgxpool.Pool) CategoryRepository {
	return &categoryRepo{pool: pool}
}

// Create inserts a new category into the database
func (r *categoryRepo) Create(ctx context.Context, category *domain.Category) (*domain.Category, *errs.AppError) {
	query := `
		INSERT INTO categories (name, display_order)
		VALUES ($1, $2)
		RETURNING id, name, display_order, created_at, updated_at`

	var c domain.Category
	err := r.pool.QueryRow(ctx, query, category.Name, category.DisplayOrder).Scan(
		&c.ID, &c.Name, &c.DisplayOrder, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsUniqueViolation(pgErr) {
				return nil, errs.NewConflictError("category with this name already exists")
			}
		}
		logger.ErrorF("failed to create category: %v", err)
		return nil, errs.NewInternalServerError("failed to create category")
	}

	return &c, nil
}

// GetAll retrieves all categories ordered by display_order
func (r *categoryRepo) GetAll(ctx context.Context) ([]domain.Category, *errs.AppError) {
	query := `SELECT id, name, display_order, created_at, updated_at FROM categories ORDER BY display_order ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		logger.ErrorF("failed to get categories: %v", err)
		return nil, errs.NewInternalServerError("failed to get categories")
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.DisplayOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			logger.ErrorF("failed to scan category: %v", err)
			return nil, errs.NewInternalServerError("failed to scan category")
		}
		categories = append(categories, c)
	}

	return categories, nil
}

// GetByID retrieves a single category by its ID
func (r *categoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, *errs.AppError) {
	query := `SELECT id, name, display_order, created_at, updated_at FROM categories WHERE id = $1`

	var c domain.Category
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.DisplayOrder, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		logger.ErrorF("failed to get category by id: %v", err)
		return nil, errs.NewNotFoundError("category not found")
	}

	return &c, nil
}

// Update modifies an existing category
func (r *categoryRepo) Update(ctx context.Context, category *domain.Category) (*domain.Category, *errs.AppError) {
	query := `
		UPDATE categories
		SET name = $1, display_order = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, display_order, created_at, updated_at`

	var c domain.Category
	err := r.pool.QueryRow(ctx, query, category.Name, category.DisplayOrder, category.ID).Scan(
		&c.ID, &c.Name, &c.DisplayOrder, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsUniqueViolation(pgErr) {
				return nil, errs.NewConflictError("category with this name already exists")
			}
		}
		logger.ErrorF("failed to update category: %v", err)
		return nil, errs.NewNotFoundError("category not found")
	}

	return &c, nil
}

// Delete removes a category by its ID
func (r *categoryRepo) Delete(ctx context.Context, id string) *errs.AppError {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsForeignKeyViolation(pgErr) {
				return errs.NewConflictError("cannot delete category: menu items still exist in this category")
			}
		}
		logger.ErrorF("failed to delete category: %v", err)
		return errs.NewInternalServerError("failed to delete category")
	}

	if result.RowsAffected() == 0 {
		return errs.NewNotFoundError("category not found")
	}

	return nil
}
