package service

import (
	"context"

	"restobill/internal/domain"
	"restobill/internal/dto"
	"restobill/internal/errs"
	"restobill/internal/logger"
	"restobill/internal/repository"
)

// valid table statuses
var validTableStatuses = map[string]bool{
	"available": true,
	"occupied":  true,
	"reserved":  true,
}

// RestaurantTableService defines the interface for restaurant table business logic
type RestaurantTableService interface {
	CreateTable(ctx context.Context, req *dto.CreateRestaurantTableRequest) (*dto.RestaurantTableResponse, *errs.AppError)
	GetAllTables(ctx context.Context) ([]dto.RestaurantTableResponse, *errs.AppError)
	GetTableByID(ctx context.Context, id string) (*dto.RestaurantTableResponse, *errs.AppError)
	UpdateTable(ctx context.Context, id string, req *dto.UpdateRestaurantTableRequest) (*dto.RestaurantTableResponse, *errs.AppError)
	UpdateTableStatus(ctx context.Context, id string, req *dto.UpdateTableStatusRequest) (*dto.RestaurantTableResponse, *errs.AppError)
	DeleteTable(ctx context.Context, id string) *errs.AppError
}

// restaurantTableService implements RestaurantTableService
type restaurantTableService struct {
	repo repository.RestaurantTableRepository
}

// NewRestaurantTableService creates a new RestaurantTableService instance
func NewRestaurantTableService(repo repository.RestaurantTableRepository) RestaurantTableService {
	return &restaurantTableService{repo: repo}
}

// CreateTable handles the business logic for creating a table
func (s *restaurantTableService) CreateTable(ctx context.Context, req *dto.CreateRestaurantTableRequest) (*dto.RestaurantTableResponse, *errs.AppError) {
	table := &domain.RestaurantTable{
		Number:   req.Number,
		Capacity: req.Capacity,
	}

	created, appErr := s.repo.Create(ctx, table)
	if appErr != nil {
		logger.ErrorF("service: failed to create table: %v", appErr.Message)
		return nil, appErr
	}

	return toTableResponse(created), nil
}

// GetAllTables retrieves all tables
func (s *restaurantTableService) GetAllTables(ctx context.Context) ([]dto.RestaurantTableResponse, *errs.AppError) {
	tables, appErr := s.repo.GetAll(ctx)
	if appErr != nil {
		logger.ErrorF("service: failed to get tables: %v", appErr.Message)
		return nil, appErr
	}

	response := make([]dto.RestaurantTableResponse, 0)
	for _, t := range tables {
		response = append(response, *toTableResponse(&t))
	}

	return response, nil
}

// GetTableByID retrieves a single table by ID
func (s *restaurantTableService) GetTableByID(ctx context.Context, id string) (*dto.RestaurantTableResponse, *errs.AppError) {
	table, appErr := s.repo.GetByID(ctx, id)
	if appErr != nil {
		logger.ErrorF("service: failed to get table: %v", appErr.Message)
		return nil, appErr
	}

	return toTableResponse(table), nil
}

// UpdateTable handles the business logic for updating a table
func (s *restaurantTableService) UpdateTable(ctx context.Context, id string, req *dto.UpdateRestaurantTableRequest) (*dto.RestaurantTableResponse, *errs.AppError) {
	table := &domain.RestaurantTable{
		ID:       id,
		Number:   req.Number,
		Capacity: req.Capacity,
	}

	updated, appErr := s.repo.Update(ctx, table)
	if appErr != nil {
		logger.ErrorF("service: failed to update table: %v", appErr.Message)
		return nil, appErr
	}

	return toTableResponse(updated), nil
}

// UpdateTableStatus handles the business logic for updating table status
func (s *restaurantTableService) UpdateTableStatus(ctx context.Context, id string, req *dto.UpdateTableStatusRequest) (*dto.RestaurantTableResponse, *errs.AppError) {
	if !validTableStatuses[req.Status] {
		return nil, errs.NewBadRequestError("invalid status: must be 'available', 'occupied', or 'reserved'")
	}

	updated, appErr := s.repo.UpdateStatus(ctx, id, req.Status)
	if appErr != nil {
		logger.ErrorF("service: failed to update table status: %v", appErr.Message)
		return nil, appErr
	}

	return toTableResponse(updated), nil
}

// DeleteTable handles the business logic for deleting a table
func (s *restaurantTableService) DeleteTable(ctx context.Context, id string) *errs.AppError {
	appErr := s.repo.Delete(ctx, id)
	if appErr != nil {
		logger.ErrorF("service: failed to delete table: %v", appErr.Message)
		return appErr
	}

	return nil
}

// toTableResponse converts a domain RestaurantTable to a RestaurantTableResponse DTO
func toTableResponse(t *domain.RestaurantTable) *dto.RestaurantTableResponse {
	return &dto.RestaurantTableResponse{
		ID:        t.ID,
		Number:    t.Number,
		Capacity:  t.Capacity,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
