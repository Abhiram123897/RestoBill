package service

import (
	"context"

	"restobill/internal/domain"
	"restobill/internal/dto"
	"restobill/internal/errs"
	"restobill/internal/logger"
	"restobill/internal/repository"
)

// valid order types
var validOrderTypes = map[string]bool{
	"dine_in":  true,
	"takeaway": true,
	"delivery": true,
}

// valid order statuses
var validOrderStatuses = map[string]bool{
	"open":      true,
	"served":    true,
	"billed":    true,
	"paid":      true,
	"cancelled": true,
}

// OrderService defines the interface for order business logic
type OrderService interface {
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.OrderResponse, *errs.AppError)
	GetAllOrders(ctx context.Context, status string, orderType string) ([]dto.OrderResponse, *errs.AppError)
	GetOrderByID(ctx context.Context, id string) (*dto.OrderResponse, *errs.AppError)
	UpdateOrderStatus(ctx context.Context, id string, req *dto.UpdateOrderStatusRequest) (*dto.OrderResponse, *errs.AppError)
	AddItemToOrder(ctx context.Context, orderID string, req *dto.AddOrderItemRequest) (*dto.OrderItemResponse, *errs.AppError)
	RemoveItemFromOrder(ctx context.Context, orderID string, itemID string) *errs.AppError
}

// orderService implements OrderService
type orderService struct {
	orderRepo    repository.OrderRepository
	menuItemRepo repository.MenuItemRepository
	tableRepo    repository.RestaurantTableRepository
}

// NewOrderService creates a new OrderService instance
func NewOrderService(orderRepo repository.OrderRepository, menuItemRepo repository.MenuItemRepository, tableRepo repository.RestaurantTableRepository) OrderService {
	return &orderService{orderRepo: orderRepo, menuItemRepo: menuItemRepo, tableRepo: tableRepo}
}

// CreateOrder handles the business logic for creating an order
func (s *orderService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.OrderResponse, *errs.AppError) {
	if !validOrderTypes[req.OrderType] {
		return nil, errs.NewBadRequestError("invalid order_type: must be 'dine_in', 'takeaway', or 'delivery'")
	}

	// dine_in requires a table_id
	if req.OrderType == "dine_in" && req.TableID == nil {
		return nil, errs.NewBadRequestError("table_id is required for dine_in orders")
	}

	// takeaway/delivery should not have a table_id
	if req.OrderType != "dine_in" && req.TableID != nil {
		return nil, errs.NewBadRequestError("table_id should not be provided for takeaway/delivery orders")
	}

	// Check if table already has an active order
	if req.TableID != nil {
		activeOrder, _ := s.orderRepo.GetActiveOrderByTableID(ctx, *req.TableID)
		if activeOrder != nil {
			return nil, errs.NewConflictError("table already has an active order")
		}
	}

	order := &domain.Order{
		TableID:   req.TableID,
		OrderType: req.OrderType,
		CreatedBy: req.CreatedBy,
	}

	created, appErr := s.orderRepo.Create(ctx, order)
	if appErr != nil {
		logger.ErrorF("service: failed to create order: %v", appErr.Message)
		return nil, appErr
	}

	// Auto-update table status to "occupied" for dine-in orders
	if req.OrderType == "dine_in" && req.TableID != nil {
		_, appErr := s.tableRepo.UpdateStatus(ctx, *req.TableID, "occupied")
		if appErr != nil {
			logger.ErrorF("service: failed to update table status: %v", appErr.Message)
		}
	}

	return s.buildOrderResponse(ctx, created)
}

// GetAllOrders retrieves all orders with optional filters
func (s *orderService) GetAllOrders(ctx context.Context, status string, orderType string) ([]dto.OrderResponse, *errs.AppError) {
	orders, appErr := s.orderRepo.GetAll(ctx, status, orderType)
	if appErr != nil {
		logger.ErrorF("service: failed to get orders: %v", appErr.Message)
		return nil, appErr
	}

	response := make([]dto.OrderResponse, 0)
	for _, o := range orders {
		orderResp, appErr := s.buildOrderResponse(ctx, &o)
		if appErr != nil {
			return nil, appErr
		}
		response = append(response, *orderResp)
	}

	return response, nil
}

// GetOrderByID retrieves a single order with its items
func (s *orderService) GetOrderByID(ctx context.Context, id string) (*dto.OrderResponse, *errs.AppError) {
	order, appErr := s.orderRepo.GetByID(ctx, id)
	if appErr != nil {
		logger.ErrorF("service: failed to get order: %v", appErr.Message)
		return nil, appErr
	}

	return s.buildOrderResponse(ctx, order)
}

// UpdateOrderStatus handles order status transitions
func (s *orderService) UpdateOrderStatus(ctx context.Context, id string, req *dto.UpdateOrderStatusRequest) (*dto.OrderResponse, *errs.AppError) {
	if !validOrderStatuses[req.Status] {
		return nil, errs.NewBadRequestError("invalid status: must be 'open', 'served', 'billed', or 'paid'")
	}

	// Get current order to validate status transition
	currentOrder, appErr := s.orderRepo.GetByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}

	if !isValidStatusTransition(currentOrder.Status, req.Status) {
		return nil, errs.NewBadRequestError("invalid status transition: cannot change from '" + currentOrder.Status + "' to '" + req.Status + "'")
	}

	updated, appErr := s.orderRepo.UpdateStatus(ctx, id, req.Status)
	if appErr != nil {
		logger.ErrorF("service: failed to update order status: %v", appErr.Message)
		return nil, appErr
	}

	// Auto-free the table when a dine-in order is paid
	if req.Status == "paid" && updated.TableID != nil {
		_, appErr := s.tableRepo.UpdateStatus(ctx, *updated.TableID, "available")
		if appErr != nil {
			logger.ErrorF("service: failed to free table: %v", appErr.Message)
		}
	}

	// Auto-free the table when a dine-in order is cancelled
	if req.Status == "cancelled" && updated.TableID != nil {
		_, appErr := s.tableRepo.UpdateStatus(ctx, *updated.TableID, "available")
		if appErr != nil {
			logger.ErrorF("service: failed to free table: %v", appErr.Message)
		}
	}

	return s.buildOrderResponse(ctx, updated)
}

// AddItemToOrder adds a menu item to an existing order
func (s *orderService) AddItemToOrder(ctx context.Context, orderID string, req *dto.AddOrderItemRequest) (*dto.OrderItemResponse, *errs.AppError) {
	// Verify order exists and is still open
	order, appErr := s.orderRepo.GetByID(ctx, orderID)
	if appErr != nil {
		return nil, appErr
	}

	if order.Status != "open" {
		return nil, errs.NewBadRequestError("cannot add items: order is already " + order.Status)
	}

	// Get menu item to snapshot the price
	menuItem, appErr := s.menuItemRepo.GetByID(ctx, req.MenuItemID)
	if appErr != nil {
		return nil, errs.NewBadRequestError("menu item not found")
	}

	if !menuItem.IsAvailable {
		return nil, errs.NewBadRequestError("menu item '" + menuItem.Name + "' is currently unavailable")
	}

	orderItem := &domain.OrderItem{
		OrderID:    orderID,
		MenuItemID: req.MenuItemID,
		Quantity:   req.Quantity,
		UnitPrice:  menuItem.Price, // Snapshot the current price
		Notes:      req.Notes,
	}

	created, appErr := s.orderRepo.AddItem(ctx, orderItem)
	if appErr != nil {
		logger.ErrorF("service: failed to add order item: %v", appErr.Message)
		return nil, appErr
	}

	return toOrderItemResponse(created), nil
}

// RemoveItemFromOrder removes an item from an order
func (s *orderService) RemoveItemFromOrder(ctx context.Context, orderID string, itemID string) *errs.AppError {
	// Verify order exists
	order, appErr := s.orderRepo.GetByID(ctx, orderID)
	if appErr != nil {
		return appErr
	}

	// Once billed, the bill total depends on these items. Void the bill first.
	if order.Status == "billed" || order.Status == "paid" || order.Status == "cancelled" {
		return errs.NewBadRequestError("cannot remove items: order is already " + order.Status + ", void the bill first")
	}

	appErr = s.orderRepo.RemoveItem(ctx, orderID, itemID)
	if appErr != nil {
		logger.ErrorF("service: failed to remove order item: %v", appErr.Message)
		return appErr
	}

	return nil
}

// buildOrderResponse constructs a full order response with items
func (s *orderService) buildOrderResponse(ctx context.Context, order *domain.Order) (*dto.OrderResponse, *errs.AppError) {
	items, appErr := s.orderRepo.GetItemsByOrderID(ctx, order.ID)
	if appErr != nil {
		return nil, appErr
	}

	itemResponses := make([]dto.OrderItemResponse, 0)
	for _, item := range items {
		itemResponses = append(itemResponses, *toOrderItemResponse(&item))
	}

	return &dto.OrderResponse{
		ID:        order.ID,
		TableID:   order.TableID,
		Status:    order.Status,
		OrderType: order.OrderType,
		CreatedBy: order.CreatedBy,
		Items:     itemResponses,
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// isValidStatusTransition checks if a manual status change is allowed.
// "billed" and "paid" are never set manually: the bill service sets "billed"
// when a bill is generated, and the payment service sets "paid" when the bill
// is fully paid. Those services update the repository directly, not this map.
func isValidStatusTransition(current string, next string) bool {
	transitions := map[string][]string{
		"open":      {"served", "cancelled"},
		"served":    {"cancelled"},
		"billed":    {},
		"paid":      {}, // terminal state
		"cancelled": {}, // terminal state
	}

	allowed, exists := transitions[current]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}

// toOrderItemResponse converts a domain OrderItem to a response DTO
func toOrderItemResponse(item *domain.OrderItem) *dto.OrderItemResponse {
	return &dto.OrderItemResponse{
		ID:         item.ID,
		OrderID:    item.OrderID,
		MenuItemID: item.MenuItemID,
		Quantity:   item.Quantity,
		UnitPrice:  item.UnitPrice,
		Notes:      item.Notes,
		CreatedAt:  item.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
