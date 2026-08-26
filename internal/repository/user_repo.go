package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"restobill/internal/db"
	"restobill/internal/domain"
	"restobill/internal/errs"
	"restobill/internal/logger"
)

// userRepo implements UserRepository using PostgreSQL
type userRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepo{pool: pool}
}

// Create inserts a new user into the database
func (r *userRepo) Create(ctx context.Context, user *domain.User) (*domain.User, *errs.AppError) {
	query := `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, password_hash, role, is_active, created_at, updated_at`

	var u domain.User
	err := r.pool.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash, user.Role).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if pgErr := db.AsPgError(err); pgErr != nil {
			if db.IsUniqueViolation(pgErr) {
				return nil, errs.NewConflictError("user with this email already exists")
			}
		}
		logger.ErrorF("failed to create user: %v", err)
		return nil, errs.NewInternalServerError("failed to create user")
	}

	return &u, nil
}

// GetByEmail retrieves a user by email
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, *errs.AppError) {
	query := `SELECT id, name, email, password_hash, role, is_active, created_at, updated_at FROM users WHERE email = $1`

	var u domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, errs.NewNotFoundError("user not found")
	}

	return &u, nil
}

// GetByID retrieves a user by ID
func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, *errs.AppError) {
	query := `SELECT id, name, email, password_hash, role, is_active, created_at, updated_at FROM users WHERE id = $1`

	var u domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, errs.NewNotFoundError("user not found")
	}

	return &u, nil
}

// GetUserCount returns the total number of users in the database
func (r *userRepo) GetUserCount(ctx context.Context) (int, *errs.AppError) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		logger.ErrorF("failed to get user count: %v", err)
		return 0, errs.NewInternalServerError("failed to get user count")
	}

	return count, nil
}
