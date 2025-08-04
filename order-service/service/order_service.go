package service

import (
	"fmt"

	"order-service/models"
	"order-service/repository"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
}

// NewOrderService creates a new instance of OrderService
func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo: repository.NewOrderRepository(),
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(request *models.CreateOrderRequest) (*models.Order, error) {
	// Validate request
	if len(request.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	// Convert request to order model
	order := &models.Order{
		CustomerID: request.CustomerID,
		Items:      make([]models.OrderItem, len(request.Items)),
	}

	for i, item := range request.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0 for game %s", item.GameName)
		}
		if item.Price < 0 {
			return nil, fmt.Errorf("price cannot be negative for game %s", item.GameName)
		}

		order.Items[i] = models.OrderItem{
			GameID:   item.GameID,
			GameName: item.GameName,
			Price:    item.Price,
			Quantity: item.Quantity,
		}
	}

	// Create order in repository
	err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	return order, nil
}

// GetOrderByID retrieves an order by its ID
func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	order, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrdersByCustomerID retrieves all orders for a specific customer
func (s *OrderService) GetOrdersByCustomerID(customerID string) ([]models.Order, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID is required")
	}

	orders, err := s.orderRepo.GetOrdersByCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// GetAllOrders retrieves all orders with pagination
func (s *OrderService) GetAllOrders(page, pageSize int) (*models.OrdersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	orders, total, err := s.orderRepo.GetAllOrders(pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	orderResponses := make([]models.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = models.OrderResponse{
			ID:         order.ID,
			CustomerID: order.CustomerID,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			OrderDate:  order.OrderDate,
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
			Items:      order.Items,
		}
	}

	return &models.OrdersListResponse{
		Orders: orderResponses,
		Total:  total,
	}, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(id string, request *models.UpdateOrderStatusRequest) error {
	if id == "" {
		return fmt.Errorf("order ID is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":    true,
		"confirmed":  true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
	}

	if !validStatuses[request.Status] {
		return fmt.Errorf("invalid status: %s", request.Status)
	}

	err := s.orderRepo.UpdateOrderStatus(id, request.Status)
	if err != nil {
		return err
	}

	return nil
}

// DeleteOrder deletes an order
func (s *OrderService) DeleteOrder(id string) error {
	if id == "" {
		return fmt.Errorf("order ID is required")
	}

	err := s.orderRepo.DeleteOrder(id)
	if err != nil {
		return err
	}

	return nil
}

// GetOrderStatistics provides basic order statistics
func (s *OrderService) GetOrderStatistics() (map[string]interface{}, error) {
	// This could be expanded to provide more detailed statistics
	orders, total, err := s.orderRepo.GetAllOrders(1000, 0) // Get recent orders for stats
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for statistics: %v", err)
	}

	statusCounts := make(map[string]int)
	var totalRevenue float64

	for _, order := range orders {
		statusCounts[order.Status]++
		if order.Status == "delivered" {
			totalRevenue += order.TotalPrice
		}
	}

	stats := map[string]interface{}{
		"total_orders":   total,
		"total_revenue":  totalRevenue,
		"status_counts":  statusCounts,
	}

	return stats, nil
}
