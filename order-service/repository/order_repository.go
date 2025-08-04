package repository

import (
	"database/sql"
	"fmt"
	"time"

	"order-service/database"
	"order-service/models"

	"github.com/google/uuid"
)

type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new instance of OrderRepository
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		db: database.DB,
	}
}

// CreateOrder creates a new order with its items
func (r *OrderRepository) CreateOrder(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Generate UUID for the order
	order.ID = uuid.New().String()
	order.Status = "pending"
	order.OrderDate = time.Now()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Calculate total price
	var totalPrice float64
	for i := range order.Items {
		order.Items[i].Subtotal = order.Items[i].Price * float64(order.Items[i].Quantity)
		totalPrice += order.Items[i].Subtotal
	}
	order.TotalPrice = totalPrice

	// Insert order
	query := `INSERT INTO orders (id, customer_id, total_price, status, order_date, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err = tx.Exec(query, order.ID, order.CustomerID, order.TotalPrice, order.Status, 
					order.OrderDate, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}

	// Insert order items
	itemQuery := `INSERT INTO order_items (id, order_id, game_id, game_name, price, quantity, subtotal) 
				  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	for i := range order.Items {
		order.Items[i].ID = uuid.New().String()
		order.Items[i].OrderID = order.ID
		
		_, err = tx.Exec(itemQuery, order.Items[i].ID, order.Items[i].OrderID, 
						order.Items[i].GameID, order.Items[i].GameName, 
						order.Items[i].Price, order.Items[i].Quantity, order.Items[i].Subtotal)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %v", err)
		}
	}

	return tx.Commit()
}

// GetOrderByID retrieves an order by its ID
func (r *OrderRepository) GetOrderByID(id string) (*models.Order, error) {
	order := &models.Order{}
	
	query := `SELECT id, customer_id, total_price, status, order_date, created_at, updated_at 
			  FROM orders WHERE id = $1`
	
	err := r.db.QueryRow(query, id).Scan(
		&order.ID, &order.CustomerID, &order.TotalPrice, &order.Status,
		&order.OrderDate, &order.CreatedAt, &order.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	// Get order items
	items, err := r.getOrderItems(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %v", err)
	}
	order.Items = items

	return order, nil
}

// GetOrdersByCustomerID retrieves all orders for a specific customer
func (r *OrderRepository) GetOrdersByCustomerID(customerID string) ([]models.Order, error) {
	query := `SELECT id, customer_id, total_price, status, order_date, created_at, updated_at 
			  FROM orders WHERE customer_id = $1 ORDER BY order_date DESC`
	
	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %v", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.CustomerID, &order.TotalPrice, &order.Status,
			&order.OrderDate, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %v", err)
		}

		// Get items for this order
		items, err := r.getOrderItems(order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items for order %s: %v", order.ID, err)
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

// GetAllOrders retrieves all orders with pagination
func (r *OrderRepository) GetAllOrders(limit, offset int) ([]models.Order, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM orders`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders count: %v", err)
	}

	// Get orders with pagination
	query := `SELECT id, customer_id, total_price, status, order_date, created_at, updated_at 
			  FROM orders ORDER BY order_date DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query orders: %v", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.CustomerID, &order.TotalPrice, &order.Status,
			&order.OrderDate, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %v", err)
		}

		// Get items for this order
		items, err := r.getOrderItems(order.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get order items for order %s: %v", order.ID, err)
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, total, nil
}

// UpdateOrderStatus updates the status of an order
func (r *OrderRepository) UpdateOrderStatus(id string, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	
	result, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// DeleteOrder deletes an order and its items
func (r *OrderRepository) DeleteOrder(id string) error {
	// The order_items will be deleted automatically due to CASCADE
	query := `DELETE FROM orders WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// getOrderItems retrieves all items for a specific order
func (r *OrderRepository) getOrderItems(orderID string) ([]models.OrderItem, error) {
	query := `SELECT id, order_id, game_id, game_name, price, quantity, subtotal 
			  FROM order_items WHERE order_id = $1 ORDER BY id`
	
	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %v", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.GameID, &item.GameName,
			&item.Price, &item.Quantity, &item.Subtotal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %v", err)
		}
		items = append(items, item)
	}

	return items, nil
}
