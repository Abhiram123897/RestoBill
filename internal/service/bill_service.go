package service

import (
	"context"
	"math"

	"restobill/internal/domain"
	"restobill/internal/dto"
	"restobill/internal/errs"
	"restobill/internal/logger"
	"restobill/internal/repository"
)

// BillService defines the interface for bill business logic
type BillService interface {
	GenerateBill(ctx context.Context, req *dto.GenerateBillRequest) (*dto.BillResponse, *errs.AppError)
	GetBillByID(ctx context.Context, id string) (*dto.BillResponse, *errs.AppError)
	GetBillByOrderID(ctx context.Context, orderID string) (*dto.BillResponse, *errs.AppError)
	ApplyDiscount(ctx context.Context, id string, req *dto.ApplyDiscountRequest) (*dto.BillResponse, *errs.AppError)
	VoidBill(ctx context.Context, id string) (*dto.BillResponse, *errs.AppError)
}

// billService implements BillService
type billService struct {
	billRepo  repository.BillRepository
	orderRepo repository.OrderRepository
}

// NewBillService creates a new BillService instance
func NewBillService(billRepo repository.BillRepository, orderRepo repository.OrderRepository) BillService {
	return &billService{billRepo: billRepo, orderRepo: orderRepo}
}

// GenerateBill creates a bill for an order
func (s *billService) GenerateBill(ctx context.Context, req *dto.GenerateBillRequest) (*dto.BillResponse, *errs.AppError) {
	// Validate tax rate
	if req.TaxRate < 0 || req.TaxRate > 100 {
		return nil, errs.NewBadRequestError("tax_rate must be between 0 and 100")
	}

	// Verify order exists
	order, appErr := s.orderRepo.GetByID(ctx, req.OrderID)
	if appErr != nil {
		return nil, errs.NewBadRequestError("order not found")
	}

	// Order must be open or served to generate bill
	if order.Status != "open" && order.Status != "served" {
		return nil, errs.NewBadRequestError("cannot generate bill: order is already " + order.Status)
	}

	// Check if bill already exists for this order (ignore void bills)
	existingBill, _ := s.billRepo.GetByOrderID(ctx, req.OrderID)
	if existingBill != nil && existingBill.Status != "void" {
		return nil, errs.NewConflictError("bill already exists for this order")
	}

	// Get order items to calculate subtotal
	items, appErr := s.orderRepo.GetItemsByOrderID(ctx, req.OrderID)
	if appErr != nil {
		return nil, appErr
	}

	if len(items) == 0 {
		return nil, errs.NewBadRequestError("cannot generate bill: order has no items")
	}

	// Calculate subtotal
	subtotal := 0.0
	for _, item := range items {
		subtotal += float64(item.Quantity) * item.UnitPrice
	}

	// Calculate tax
	taxAmount := roundToTwo(subtotal * (req.TaxRate / 100))

	// Calculate total
	discountAmount := req.DiscountAmount
	if discountAmount < 0 {
		return nil, errs.NewBadRequestError("discount_amount cannot be negative")
	}
	if discountAmount > (subtotal + taxAmount) {
		return nil, errs.NewBadRequestError("discount_amount cannot exceed bill total")
	}

	total := math.Round(subtotal + taxAmount - discountAmount)

	bill := &domain.Bill{
		OrderID:        req.OrderID,
		Subtotal:       subtotal,
		TaxAmount:      taxAmount,
		TaxRate:        req.TaxRate,
		DiscountAmount: discountAmount,
		DiscountReason: req.DiscountReason,
		Total:          total,
		GeneratedBy:    req.GeneratedBy,
	}

	created, appErr := s.billRepo.Create(ctx, bill)
	if appErr != nil {
		logger.ErrorF("service: failed to create bill: %v", appErr.Message)
		return nil, appErr
	}

	// Update order status to "billed"
	_, appErr = s.orderRepo.UpdateStatus(ctx, req.OrderID, "billed")
	if appErr != nil {
		logger.ErrorF("service: failed to update order status to billed: %v", appErr.Message)
	}

	return toBillResponse(created), nil
}

// GetBillByID retrieves a bill by its ID
func (s *billService) GetBillByID(ctx context.Context, id string) (*dto.BillResponse, *errs.AppError) {
	bill, appErr := s.billRepo.GetByID(ctx, id)
	if appErr != nil {
		logger.ErrorF("service: failed to get bill: %v", appErr.Message)
		return nil, appErr
	}

	return toBillResponse(bill), nil
}

// GetBillByOrderID retrieves a bill by order ID
func (s *billService) GetBillByOrderID(ctx context.Context, orderID string) (*dto.BillResponse, *errs.AppError) {
	bill, _ := s.billRepo.GetByOrderID(ctx, orderID)
	if bill == nil {
		return nil, errs.NewNotFoundError("no bill found for this order")
	}

	return toBillResponse(bill), nil
}

// ApplyDiscount applies or updates discount on a pending bill
func (s *billService) ApplyDiscount(ctx context.Context, id string, req *dto.ApplyDiscountRequest) (*dto.BillResponse, *errs.AppError) {
	bill, appErr := s.billRepo.GetByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}

	if bill.Status != "pending" {
		return nil, errs.NewBadRequestError("cannot apply discount: bill is already " + bill.Status)
	}

	if req.DiscountAmount < 0 {
		return nil, errs.NewBadRequestError("discount_amount cannot be negative")
	}

	if req.DiscountAmount > (bill.Subtotal + bill.TaxAmount) {
		return nil, errs.NewBadRequestError("discount_amount cannot exceed bill total")
	}

	// Recalculate total
	bill.DiscountAmount = req.DiscountAmount
	bill.DiscountReason = req.DiscountReason
	bill.Total = math.Round(bill.Subtotal + bill.TaxAmount - req.DiscountAmount)

	updated, appErr := s.billRepo.Update(ctx, bill)
	if appErr != nil {
		logger.ErrorF("service: failed to apply discount: %v", appErr.Message)
		return nil, appErr
	}

	return toBillResponse(updated), nil
}

// VoidBill marks a bill as void (cancelled)
func (s *billService) VoidBill(ctx context.Context, id string) (*dto.BillResponse, *errs.AppError) {
	bill, appErr := s.billRepo.GetByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}

	if bill.Status == "paid" {
		return nil, errs.NewBadRequestError("cannot void a paid bill")
	}

	if bill.Status == "void" {
		return nil, errs.NewBadRequestError("bill is already void")
	}

	updated, appErr := s.billRepo.UpdateStatus(ctx, id, "void")
	if appErr != nil {
		logger.ErrorF("service: failed to void bill: %v", appErr.Message)
		return nil, appErr
	}

	// Revert order status back to "served" or "open"
	_, appErr = s.orderRepo.UpdateStatus(ctx, bill.OrderID, "open")
	if appErr != nil {
		logger.ErrorF("service: failed to revert order status: %v", appErr.Message)
	}

	return toBillResponse(updated), nil
}

// toBillResponse converts a domain Bill to a BillResponse DTO
func toBillResponse(b *domain.Bill) *dto.BillResponse {
	return &dto.BillResponse{
		ID:             b.ID,
		OrderID:        b.OrderID,
		Subtotal:       b.Subtotal,
		TaxAmount:      b.TaxAmount,
		TaxRate:        b.TaxRate,
		DiscountAmount: b.DiscountAmount,
		DiscountReason: b.DiscountReason,
		Total:          b.Total,
		Status:         b.Status,
		GeneratedBy:    b.GeneratedBy,
		CreatedAt:      b.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// roundToTwo rounds a float to 2 decimal places
func roundToTwo(val float64) float64 {
	return math.Round(val*100) / 100
}
