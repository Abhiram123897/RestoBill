package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"restobill/internal/db"
	"restobill/internal/domain"
	"restobill/internal/errs"
	"restobill/internal/logger"
)

// billRepo implements BillRepository using PostgreSQL
type billRepo struct {
	pool *pgxpool.Pool
}

// NewBillRepository creates a new BillRepository instance
func NewBillRepository(pool *pgxpool.Pool) BillRepository {
	return &billRepo{pool: pool}
}

// Create inserts a new bill into the database
func (r *billRepo) Create(ctx context.Context, bill *domain.Bill) (*domain.Bill, *errs.AppError) {
	query := `
		INSERT INTO bills (order_id, subtotal, tax_amount, tax_rate, discount_amount, discount_reason, total, generated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, order_id, subtotal, tax_amount, tax_rate, discount_amount, discount_reason, total, status, generated_by, created_at`

	var b domain.Bill
	err := r.pool.QueryRow(ctx, query,
		bill.OrderID, bill.Subtotal, bill.TaxAmount, bill.TaxRate,
		bill.DiscountAmount, bill.DiscountReason, bill.Total, bill.GeneratedBy,
	).Scan(
		&b.ID, &b.OrderID, &b.Subtotal, &b.TaxAmount, &b.TaxRate,
		&b.DiscountAmount, &b.DiscountReason, &b.Total, &b.Status, &b.GeneratedBy, &b.CreatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsUniqueViolation(pgErr) {
				return nil, errs.NewConflictError("bill already exists for this order")
			}
			if db.IsForeignKeyViolation(pgErr) {
				return nil, errs.NewBadRequestError("invalid order_id or generated_by")
			}
		}
		logger.ErrorF("failed to create bill: %v", err)
		return nil, errs.NewInternalServerError("failed to create bill")
	}

	return &b, nil
}

// GetByID retrieves a bill by its ID
func (r *billRepo) GetByID(ctx context.Context, id string) (*domain.Bill, *errs.AppError) {
	query := `SELECT id, order_id, subtotal, tax_amount, tax_rate, discount_amount, discount_reason, total, status, generated_by, created_at FROM bills WHERE id = $1`

	var b domain.Bill
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.OrderID, &b.Subtotal, &b.TaxAmount, &b.TaxRate,
		&b.DiscountAmount, &b.DiscountReason, &b.Total, &b.Status, &b.GeneratedBy, &b.CreatedAt,
	)
	if err != nil {
		logger.ErrorF("failed to get bill by id: %v", err)
		return nil, errs.NewNotFoundError("bill not found")
	}

	return &b, nil
}

// GetByOrderID retrieves the current (non-void) bill for an order.
// An order can have several bills after void-and-rebill, so void bills are skipped
// and the newest one wins.
func (r *billRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.Bill, *errs.AppError) {
	query := `SELECT id, order_id, subtotal, tax_amount, tax_rate, discount_amount, discount_reason, total, status, generated_by, created_at FROM bills WHERE order_id = $1 AND status != 'void' ORDER BY created_at DESC LIMIT 1`

	var b domain.Bill
	err := r.pool.QueryRow(ctx, query, orderID).Scan(
		&b.ID, &b.OrderID, &b.Subtotal, &b.TaxAmount, &b.TaxRate,
		&b.DiscountAmount, &b.DiscountReason, &b.Total, &b.Status, &b.GeneratedBy, &b.CreatedAt,
	)
	if err != nil {
		return nil, nil // No bill found for this order — not an error
	}

	return &b, nil
}

// Update modifies an existing bill (used for applying discounts)
func (r *billRepo) Update(ctx context.Context, bill *domain.Bill) (*domain.Bill, *errs.AppError) {
	query := `
		UPDATE bills
		SET discount_amount = $1, discount_reason = $2, total = $3
		WHERE id = $4
		RETURNING id, order_id, subtotal, tax_amount, tax_rate, discount_amount, discount_reason, total, status, generated_by, created_at`

	var b domain.Bill
	err := r.pool.QueryRow(ctx, query, bill.DiscountAmount, bill.DiscountReason, bill.Total, bill.ID).Scan(
		&b.ID, &b.OrderID, &b.Subtotal, &b.TaxAmount, &b.TaxRate,
		&b.DiscountAmount, &b.DiscountReason, &b.Total, &b.Status, &b.GeneratedBy, &b.CreatedAt,
	)
	if err != nil {
		logger.ErrorF("failed to update bill: %v", err)
		return nil, errs.NewNotFoundError("bill not found")
	}

	return &b, nil
}

// UpdateStatus changes the status of a bill
func (r *billRepo) UpdateStatus(ctx context.Context, id string, status string) (*domain.Bill, *errs.AppError) {
	query := `
		UPDATE bills
		SET status = $1
		WHERE id = $2
		RETURNING id, order_id, subtotal, tax_amount, tax_rate, discount_amount, discount_reason, total, status, generated_by, created_at`

	var b domain.Bill
	err := r.pool.QueryRow(ctx, query, status, id).Scan(
		&b.ID, &b.OrderID, &b.Subtotal, &b.TaxAmount, &b.TaxRate,
		&b.DiscountAmount, &b.DiscountReason, &b.Total, &b.Status, &b.GeneratedBy, &b.CreatedAt,
	)
	if err != nil {
		logger.ErrorF("failed to update bill status: %v", err)
		return nil, errs.NewNotFoundError("bill not found")
	}

	return &b, nil
}
