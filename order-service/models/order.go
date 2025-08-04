package models

import (
	"time"
)

// Order represents an order entity
type Order struct {
	ID          string      `json:"id" db:"id"`
	CustomerID  string      `json:"customer_id" db:"customer_id" binding:"required"`
	TotalPrice  float64     `json:"total_price" db:"total_price"`
	Status      string      `json:"status" db:"status"`
	OrderDate   time.Time   `json:"order_date" db:"order_date"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
	Items       []OrderItem `json:"items,omitempty"`
}

// OrderItem represents an item within an order
type OrderItem struct {
	ID       string  `json:"id" db:"id"`
	OrderID  string  `json:"order_id" db:"order_id"`
	GameID   int     `json:"game_id" db:"game_id" binding:"required"`
	GameName string  `json:"game_name" db:"game_name"`
	Price    float64 `json:"price" db:"price" binding:"required,min=0"`
	Quantity int     `json:"quantity" db:"quantity" binding:"required,min=1"`
	Subtotal float64 `json:"subtotal" db:"subtotal"`
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	CustomerID string                   `json:"customer_id" binding:"required"`
	Items      []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
}

// CreateOrderItemRequest represents an item in the order creation request
type CreateOrderItemRequest struct {
	GameID   int     `json:"game_id" binding:"required"`
	GameName string  `json:"game_name" binding:"required"`
	Price    float64 `json:"price" binding:"required,min=0"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
}

// UpdateOrderStatusRequest represents the request body for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed processing shipped delivered cancelled"`
}

// OrderResponse represents the response structure for order queries
type OrderResponse struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	OrderDate  time.Time   `json:"order_date"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Items      []OrderItem `json:"items"`
}

// OrdersListResponse represents the response for listing orders
type OrdersListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
}
