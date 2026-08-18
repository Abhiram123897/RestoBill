package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"restobill/internal/db"
	"restobill/internal/domain"
	"restobill/internal/errs"
	"restobill/internal/logger"
)

// restaurantTableRepo implements RestaurantTableRepository using PostgreSQL
type restaurantTableRepo struct {
	pool *pgxpool.Pool
}

// NewRestaurantTableRepository creates a new RestaurantTableRepository instance
func NewRestaurantTableRepository(pool *pgxpool.Pool) RestaurantTableRepository {
	return &restaurantTableRepo{pool: pool}
}

// Create inserts a new restaurant table into the database
func (r *restaurantTableRepo) Create(ctx context.Context, table *domain.RestaurantTable) (*domain.RestaurantTable, *errs.AppError) {
	query := `
		INSERT INTO restaurant_tables (number, capacity)
		VALUES ($1, $2)
		RETURNING id, number, capacity, status, created_at`

	var t domain.RestaurantTable
	err := r.pool.QueryRow(ctx, query, table.Number, table.Capacity).Scan(
		&t.ID, &t.Number, &t.Capacity, &t.Status, &t.CreatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsUniqueViolation(pgErr) {
				return nil, errs.NewConflictError("table with this number already exists")
			}
		}
		logger.ErrorF("failed to create restaurant table: %v", err)
		return nil, errs.NewInternalServerError("failed to create restaurant table")
	}

	return &t, nil
}

// GetAll retrieves all restaurant tables ordered by number
func (r *restaurantTableRepo) GetAll(ctx context.Context) ([]domain.RestaurantTable, *errs.AppError) {
	query := `SELECT id, number, capacity, status, created_at FROM restaurant_tables ORDER BY number ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		logger.ErrorF("failed to get restaurant tables: %v", err)
		return nil, errs.NewInternalServerError("failed to get restaurant tables")
	}
	defer rows.Close()

	var tables []domain.RestaurantTable
	for rows.Next() {
		var t domain.RestaurantTable
		if err := rows.Scan(&t.ID, &t.Number, &t.Capacity, &t.Status, &t.CreatedAt); err != nil {
			logger.ErrorF("failed to scan restaurant table: %v", err)
			return nil, errs.NewInternalServerError("failed to scan restaurant table")
		}
		tables = append(tables, t)
	}

	return tables, nil
}

// GetByID retrieves a single restaurant table by its ID
func (r *restaurantTableRepo) GetByID(ctx context.Context, id string) (*domain.RestaurantTable, *errs.AppError) {
	query := `SELECT id, number, capacity, status, created_at FROM restaurant_tables WHERE id = $1`

	var t domain.RestaurantTable
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Number, &t.Capacity, &t.Status, &t.CreatedAt,
	)
	if err != nil {
		logger.ErrorF("failed to get restaurant table by id: %v", err)
		return nil, errs.NewNotFoundError("table not found")
	}

	return &t, nil
}

// Update modifies an existing restaurant table
func (r *restaurantTableRepo) Update(ctx context.Context, table *domain.RestaurantTable) (*domain.RestaurantTable, *errs.AppError) {
	query := `
		UPDATE restaurant_tables
		SET number = $1, capacity = $2
		WHERE id = $3
		RETURNING id, number, capacity, status, created_at`

	var t domain.RestaurantTable
	err := r.pool.QueryRow(ctx, query, table.Number, table.Capacity, table.ID).Scan(
		&t.ID, &t.Number, &t.Capacity, &t.Status, &t.CreatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsUniqueViolation(pgErr) {
				return nil, errs.NewConflictError("table with this number already exists")
			}
		}
		logger.ErrorF("failed to update restaurant table: %v", err)
		return nil, errs.NewNotFoundError("table not found")
	}

	return &t, nil
}

// UpdateStatus changes the status of a restaurant table
func (r *restaurantTableRepo) UpdateStatus(ctx context.Context, id string, status string) (*domain.RestaurantTable, *errs.AppError) {
	query := `
		UPDATE restaurant_tables
		SET status = $1
		WHERE id = $2
		RETURNING id, number, capacity, status, created_at`

	var t domain.RestaurantTable
	err := r.pool.QueryRow(ctx, query, status, id).Scan(
		&t.ID, &t.Number, &t.Capacity, &t.Status, &t.CreatedAt,
	)
	if err != nil {
		logger.ErrorF("failed to update table status: %v", err)
		return nil, errs.NewNotFoundError("table not found")
	}

	return &t, nil
}

// Delete removes a restaurant table by its ID
func (r *restaurantTableRepo) Delete(ctx context.Context, id string) *errs.AppError {
	query := `DELETE FROM restaurant_tables WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsForeignKeyViolation(pgErr) {
				return errs.NewConflictError("cannot delete table: orders still linked to this table")
			}
		}
		logger.ErrorF("failed to delete restaurant table: %v", err)
		return errs.NewInternalServerError("failed to delete restaurant table")
	}

	if result.RowsAffected() == 0 {
		return errs.NewNotFoundError("table not found")
	}

	return nil
}
