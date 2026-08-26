package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"restobill/internal/domain"
	"restobill/internal/dto"
	"restobill/internal/errs"
	"restobill/internal/logger"
	"restobill/internal/repository"
)

// valid user roles
var validRoles = map[string]bool{
	"admin":   true,
	"cashier": true,
	"waiter":  true,
}

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest, callerRole string) (*dto.AuthResponse, *errs.AppError)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, *errs.AppError)
	GetCurrentUser(ctx context.Context, userID string) (*dto.UserResponse, *errs.AppError)
}

// authService implements AuthService
type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

// Register creates a new user with hashed password and returns a JWT token
func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest, callerRole string) (*dto.AuthResponse, *errs.AppError) {
	// Check if this is the first user (auto-admin)
	userCount, appErr := s.userRepo.GetUserCount(ctx)
	if appErr != nil {
		return nil, appErr
	}

	if userCount == 0 {
		// First user is always admin
		req.Role = "admin"
	} else {
		// After first user, only admin can register
		if callerRole != "admin" {
			return nil, errs.NewForbiddenError("only admin can register new users")
		}
	}

	// Validate role
	if !validRoles[req.Role] {
		return nil, errs.NewBadRequestError("invalid role: must be 'admin', 'cashier', or 'waiter'")
	}

	// Validate password length
	if len(req.Password) < 6 {
		return nil, errs.NewBadRequestError("password must be at least 6 characters")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorF("failed to hash password: %v", err)
		return nil, errs.NewInternalServerError("failed to process registration")
	}

	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
	}

	created, appErr := s.userRepo.Create(ctx, user)
	if appErr != nil {
		logger.ErrorF("service: failed to register user: %v", appErr.Message)
		return nil, appErr
	}

	// Generate JWT token
	token, err := s.generateToken(created)
	if err != nil {
		logger.ErrorF("failed to generate token: %v", err)
		return nil, errs.NewInternalServerError("failed to generate token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  toUserResponse(created),
	}, nil
}

// Login verifies credentials and returns a JWT token
func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, *errs.AppError) {
	// Get user by email
	user, appErr := s.userRepo.GetByEmail(ctx, req.Email)
	if appErr != nil {
		return nil, errs.NewUnauthorizedError("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errs.NewUnauthorizedError("account is deactivated")
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errs.NewUnauthorizedError("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		logger.ErrorF("failed to generate token: %v", err)
		return nil, errs.NewInternalServerError("failed to generate token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

// GetCurrentUser retrieves the current user from token
func (s *authService) GetCurrentUser(ctx context.Context, userID string) (*dto.UserResponse, *errs.AppError) {
	user, appErr := s.userRepo.GetByID(ctx, userID)
	if appErr != nil {
		return nil, appErr
	}

	response := toUserResponse(user)
	return &response, nil
}

// generateToken creates a JWT token for a user
func (s *authService) generateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// toUserResponse converts domain User to UserResponse DTO
func toUserResponse(u *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Role:     u.Role,
		IsActive: u.IsActive,
	}
}
