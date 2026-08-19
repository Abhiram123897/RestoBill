package service

import (
	"context"

	"restobill/internal/domain"
	"restobill/internal/dto"
	"restobill/internal/errs"
	"restobill/internal/logger"
	"restobill/internal/repository"
)

// valid payment methods
var validPaymentMethods = map[string]bool{
	"cash": true,
	"card": true,
	"upi":  true,
}

// PaymentService defines the interface for payment business logic
type PaymentService interface {
	RecordPayment(ctx context.Context, req *dto.RecordPaymentRequest) (*dto.PaymentResponse, *errs.AppError)
	GetPaymentsByBillID(ctx context.Context, billID string) ([]dto.PaymentResponse, *errs.AppError)
}

// paymentService implements PaymentService
type paymentService struct {
	paymentRepo repository.PaymentRepository
	billRepo    repository.BillRepository
	orderRepo   repository.OrderRepository
	tableRepo   repository.RestaurantTableRepository
}

// NewPaymentService creates a new PaymentService instance
func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	billRepo repository.BillRepository,
	orderRepo repository.OrderRepository,
	tableRepo repository.RestaurantTableRepository,
) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		billRepo:    billRepo,
		orderRepo:   orderRepo,
		tableRepo:   tableRepo,
	}
}

// RecordPayment records a payment against a bill
func (s *paymentService) RecordPayment(ctx context.Context, req *dto.RecordPaymentRequest) (*dto.PaymentResponse, *errs.AppError) {
	// Validate payment method
	if !validPaymentMethods[req.Method] {
		return nil, errs.NewBadRequestError("invalid method: must be 'cash', 'card', or 'upi'")
	}

	// Validate amount
	if req.Amount <= 0 {
		return nil, errs.NewBadRequestError("amount must be greater than 0")
	}

	// Get bill
	bill, appErr := s.billRepo.GetByID(ctx, req.BillID)
	if appErr != nil {
		return nil, appErr
	}

	// Can only pay pending bills
	if bill.Status != "pending" {
		return nil, errs.NewBadRequestError("cannot record payment: bill is already " + bill.Status)
	}

	// Check if payment amount exceeds remaining balance
	existingPayments, appErr := s.paymentRepo.GetByBillID(ctx, req.BillID)
	if appErr != nil {
		return nil, appErr
	}

	totalPaid := 0.0
	for _, p := range existingPayments {
		totalPaid += p.Amount
	}

	remaining := bill.Total - totalPaid
	if req.Amount > remaining {
		return nil, errs.NewBadRequestError("payment amount exceeds remaining balance")
	}

	// Record payment
	payment := &domain.Payment{
		BillID:    req.BillID,
		Amount:    req.Amount,
		Method:    req.Method,
		Reference: req.Reference,
	}

	created, appErr := s.paymentRepo.Create(ctx, payment)
	if appErr != nil {
		logger.ErrorF("service: failed to record payment: %v", appErr.Message)
		return nil, appErr
	}

	// Check if bill is fully paid
	totalPaid += req.Amount
	if totalPaid >= bill.Total {
		// Mark bill as paid
		_, appErr = s.billRepo.UpdateStatus(ctx, bill.ID, "paid")
		if appErr != nil {
			logger.ErrorF("service: failed to mark bill as paid: %v", appErr.Message)
		}

		// Mark order as paid
		_, appErr = s.orderRepo.UpdateStatus(ctx, bill.OrderID, "paid")
		if appErr != nil {
			logger.ErrorF("service: failed to mark order as paid: %v", appErr.Message)
		}

		// Free the table (if dine-in)
		order, _ := s.orderRepo.GetByID(ctx, bill.OrderID)
		if order != nil && order.TableID != nil {
			_, appErr = s.tableRepo.UpdateStatus(ctx, *order.TableID, "available")
			if appErr != nil {
				logger.ErrorF("service: failed to free table: %v", appErr.Message)
			}
		}
	}

	return toPaymentResponse(created), nil
}

// GetPaymentsByBillID retrieves all payments for a bill
func (s *paymentService) GetPaymentsByBillID(ctx context.Context, billID string) ([]dto.PaymentResponse, *errs.AppError) {
	payments, appErr := s.paymentRepo.GetByBillID(ctx, billID)
	if appErr != nil {
		logger.ErrorF("service: failed to get payments: %v", appErr.Message)
		return nil, appErr
	}

	response := make([]dto.PaymentResponse, 0)
	for _, p := range payments {
		response = append(response, *toPaymentResponse(&p))
	}

	return response, nil
}

// toPaymentResponse converts a domain Payment to a PaymentResponse DTO
func toPaymentResponse(p *domain.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:        p.ID,
		BillID:    p.BillID,
		Amount:    p.Amount,
		Method:    p.Method,
		Reference: p.Reference,
		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
