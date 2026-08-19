package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"restobill/internal/db"
	"restobill/internal/domain"
	"restobill/internal/errs"
	"restobill/internal/logger"
)

// paymentRepo implements PaymentRepository using PostgreSQL
type paymentRepo struct {
	pool *pgxpool.Pool
}

// NewPaymentRepository creates a new PaymentRepository instance
func NewPaymentRepository(pool *pgxpool.Pool) PaymentRepository {
	return &paymentRepo{pool: pool}
}

// Create inserts a new payment into the database
func (r *paymentRepo) Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, *errs.AppError) {
	query := `
		INSERT INTO payments (bill_id, amount, method, reference)
		VALUES ($1, $2, $3, $4)
		RETURNING id, bill_id, amount, method, reference, created_at`

	var p domain.Payment
	err := r.pool.QueryRow(ctx, query, payment.BillID, payment.Amount, payment.Method, payment.Reference).Scan(
		&p.ID, &p.BillID, &p.Amount, &p.Method, &p.Reference, &p.CreatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsForeignKeyViolation(pgErr) {
				return nil, errs.NewBadRequestError("bill does not exist")
			}
		}
		logger.ErrorF("failed to create payment: %v", err)
		return nil, errs.NewInternalServerError("failed to create payment")
	}

	return &p, nil
}

// GetByBillID retrieves all payments for a specific bill
func (r *paymentRepo) GetByBillID(ctx context.Context, billID string) ([]domain.Payment, *errs.AppError) {
	query := `SELECT id, bill_id, amount, method, reference, created_at FROM payments WHERE bill_id = $1 ORDER BY created_at ASC`

	rows, err := r.pool.Query(ctx, query, billID)
	if err != nil {
		logger.ErrorF("failed to get payments by bill id: %v", err)
		return nil, errs.NewInternalServerError("failed to get payments")
	}
	defer rows.Close()

	var payments []domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.BillID, &p.Amount, &p.Method, &p.Reference, &p.CreatedAt); err != nil {
			logger.ErrorF("failed to scan payment: %v", err)
			return nil, errs.NewInternalServerError("failed to scan payment")
		}
		payments = append(payments, p)
	}

	return payments, nil
}
