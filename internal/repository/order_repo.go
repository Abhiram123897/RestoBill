package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"restobill/internal/db"
	"restobill/internal/domain"
	"restobill/internal/errs"
	"restobill/internal/logger"
)

// orderRepo implements OrderRepository using PostgreSQL
type orderRepo struct {
	pool *pgxpool.Pool
}

// NewOrderRepository creates a new OrderRepository instance
func NewOrderRepository(pool *pgxpool.Pool) OrderRepository {
	return &orderRepo{pool: pool}
}

// Create inserts a new order into the database
func (r *orderRepo) Create(ctx context.Context, order *domain.Order) (*domain.Order, *errs.AppError) {
	query := `
		INSERT INTO orders (table_id, order_type, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, table_id, status, order_type, created_by, created_at, updated_at`

	var o domain.Order
	err := r.pool.QueryRow(ctx, query, order.TableID, order.OrderType, order.CreatedBy).Scan(
		&o.ID, &o.TableID, &o.Status, &o.OrderType, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsForeignKeyViolation(pgErr) {
				return nil, errs.NewBadRequestError("invalid table_id or created_by: referenced resource does not exist")
			}
		}
		logger.ErrorF("failed to create order: %v", err)
		return nil, errs.NewInternalServerError("failed to create order")
	}

	return &o, nil
}

// GetAll retrieves all orders with optional filters
func (r *orderRepo) GetAll(ctx context.Context, status string, orderType string) ([]domain.Order, *errs.AppError) {
	query := `SELECT id, table_id, status, order_type, created_by, created_at, updated_at FROM orders`

	var args []interface{}
	var conditions []string
	argIndex := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	if orderType != "" {
		conditions = append(conditions, fmt.Sprintf("order_type = $%d", argIndex))
		args = append(args, orderType)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		logger.ErrorF("failed to get orders: %v", err)
		return nil, errs.NewInternalServerError("failed to get orders")
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.TableID, &o.Status, &o.OrderType, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt); err != nil {
			logger.ErrorF("failed to scan order: %v", err)
			return nil, errs.NewInternalServerError("failed to scan order")
		}
		orders = append(orders, o)
	}

	return orders, nil
}

// GetByID retrieves a single order by its ID
func (r *orderRepo) GetByID(ctx context.Context, id string) (*domain.Order, *errs.AppError) {
	query := `SELECT id, table_id, status, order_type, created_by, created_at, updated_at FROM orders WHERE id = $1`

	var o domain.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.TableID, &o.Status, &o.OrderType, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		logger.ErrorF("failed to get order by id: %v", err)
		return nil, errs.NewNotFoundError("order not found")
	}

	return &o, nil
}

// GetActiveOrderByTableID checks if a table has an active order (open, served, or billed)
func (r *orderRepo) GetActiveOrderByTableID(ctx context.Context, tableID string) (*domain.Order, *errs.AppError) {
	query := `SELECT id, table_id, status, order_type, created_by, created_at, updated_at FROM orders WHERE table_id = $1 AND status IN ('open', 'served', 'billed') LIMIT 1`

	var o domain.Order
	err := r.pool.QueryRow(ctx, query, tableID).Scan(
		&o.ID, &o.TableID, &o.Status, &o.OrderType, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		// No active order found — this is expected, not an error
		return nil, nil
	}

	return &o, nil
}

// UpdateStatus changes the status of an order
func (r *orderRepo) UpdateStatus(ctx context.Context, id string, status string) (*domain.Order, *errs.AppError) {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, table_id, status, order_type, created_by, created_at, updated_at`

	var o domain.Order
	err := r.pool.QueryRow(ctx, query, status, id).Scan(
		&o.ID, &o.TableID, &o.Status, &o.OrderType, &o.CreatedBy, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		logger.ErrorF("failed to update order status: %v", err)
		return nil, errs.NewNotFoundError("order not found")
	}

	return &o, nil
}

// AddItem inserts a new item into an order
func (r *orderRepo) AddItem(ctx context.Context, item *domain.OrderItem) (*domain.OrderItem, *errs.AppError) {
	query := `
		INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, order_id, menu_item_id, quantity, unit_price, notes, created_at`

	var oi domain.OrderItem
	err := r.pool.QueryRow(ctx, query, item.OrderID, item.MenuItemID, item.Quantity, item.UnitPrice, item.Notes).Scan(
		&oi.ID, &oi.OrderID, &oi.MenuItemID, &oi.Quantity, &oi.UnitPrice, &oi.Notes, &oi.CreatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsForeignKeyViolation(pgErr) {
				return nil, errs.NewBadRequestError("invalid order_id or menu_item_id: referenced resource does not exist")
			}
		}
		logger.ErrorF("failed to add order item: %v", err)
		return nil, errs.NewInternalServerError("failed to add order item")
	}

	return &oi, nil
}

// GetItemsByOrderID retrieves all items for a specific order
func (r *orderRepo) GetItemsByOrderID(ctx context.Context, orderID string) ([]domain.OrderItem, *errs.AppError) {
	query := `SELECT id, order_id, menu_item_id, quantity, unit_price, notes, created_at FROM order_items WHERE order_id = $1 ORDER BY created_at ASC`

	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		logger.ErrorF("failed to get order items: %v", err)
		return nil, errs.NewInternalServerError("failed to get order items")
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var oi domain.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.MenuItemID, &oi.Quantity, &oi.UnitPrice, &oi.Notes, &oi.CreatedAt); err != nil {
			logger.ErrorF("failed to scan order item: %v", err)
			return nil, errs.NewInternalServerError("failed to scan order item")
		}
		items = append(items, oi)
	}

	return items, nil
}

// RemoveItem deletes an item from an order
func (r *orderRepo) RemoveItem(ctx context.Context, orderID string, itemID string) *errs.AppError {
	query := `DELETE FROM order_items WHERE id = $1 AND order_id = $2`

	result, err := r.pool.Exec(ctx, query, itemID, orderID)
	if err != nil {
		logger.ErrorF("failed to remove order item: %v", err)
		return errs.NewInternalServerError("failed to remove order item")
	}

	if result.RowsAffected() == 0 {
		return errs.NewNotFoundError("order item not found")
	}

	return nil
}
