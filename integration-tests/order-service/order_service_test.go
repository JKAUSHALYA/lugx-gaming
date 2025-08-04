package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const orderServiceBaseURL = "http://localhost:30081"

type Order struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id"`
	TotalPrice  float64     `json:"total_price"`
	Status      string      `json:"status"`
	OrderDate   time.Time   `json:"order_date"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Items       []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID       string  `json:"id"`
	OrderID  string  `json:"order_id"`
	GameID   int     `json:"game_id"`
	GameName string  `json:"game_name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}

type CreateOrderRequest struct {
	CustomerID string `json:"customer_id"`
	Items      []struct {
		GameID   int     `json:"game_id"`
		GameName string  `json:"game_name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	} `json:"items"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type OrderStats struct {
	TotalOrders    int     `json:"total_orders"`
	TotalRevenue   float64 `json:"total_revenue"`
	AverageOrder   float64 `json:"average_order"`
	PendingOrders  int     `json:"pending_orders"`
	CompletedOrders int    `json:"completed_orders"`
}

// Response wrappers used by the order service
type CreateOrderResponse struct {
	Message string `json:"message"`
	Order   Order  `json:"order"`
}

type GetOrderResponse struct {
	Order Order `json:"order"`
}

type GetOrdersResponse struct {
	Orders []Order `json:"orders"`
	Total  int     `json:"total"`
}

type OrderStatsResponse struct {
	Stats OrderStats `json:"stats"`
}

func TestOrderServiceHealth(t *testing.T) {
	resp, err := http.Get(orderServiceBaseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to make health check request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestCreateOrder(t *testing.T) {
	orderRequest := CreateOrderRequest{
		CustomerID: "customer123",
		Items: []struct {
			GameID   int     `json:"game_id"`
			GameName string  `json:"game_name"`
			Price    float64 `json:"price"`
			Quantity int     `json:"quantity"`
		}{
			{
				GameID:   1,
				GameName: "Test Game 1",
				Price:    29.99,
				Quantity: 2,
			},
			{
				GameID:   2,
				GameName: "Test Game 2",
				Price:    39.99,
				Quantity: 1,
			},
		},
	}

	jsonData, err := json.Marshal(orderRequest)
	if err != nil {
		t.Fatalf("Failed to marshal order request: %v", err)
	}

	resp, err := http.Post(orderServiceBaseURL+"/api/v1/orders", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 201 or 200, got %d", resp.StatusCode)
	}

	var response CreateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	order := response.Order

	if order.CustomerID != orderRequest.CustomerID {
		t.Errorf("Expected customer ID %s, got %s", orderRequest.CustomerID, order.CustomerID)
	}

	if len(order.Items) != len(orderRequest.Items) {
		t.Errorf("Expected %d items, got %d", len(orderRequest.Items), len(order.Items))
	}

	expectedTotal := (29.99 * 2) + (39.99 * 1)
	if order.TotalPrice != expectedTotal {
		t.Errorf("Expected total price %.2f, got %.2f", expectedTotal, order.TotalPrice)
	}
}

func TestGetAllOrders(t *testing.T) {
	resp, err := http.Get(orderServiceBaseURL + "/api/v1/orders")
	if err != nil {
		t.Fatalf("Failed to get all orders: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// The response might be paginated or a simple array
	// Check if it has orders data
	if orders, exists := response["orders"]; exists {
		if ordersArray, ok := orders.([]interface{}); ok {
			t.Logf("Found %d orders in paginated response", len(ordersArray))
		}
	} else {
		// Try to decode as array directly
		var orders []Order
		resp2, err := http.Get(orderServiceBaseURL + "/api/v1/orders")
		if err != nil {
			t.Fatalf("Failed to get all orders (second attempt): %v", err)
		}
		defer resp2.Body.Close()
		
		if err := json.NewDecoder(resp2.Body).Decode(&orders); err == nil {
			t.Logf("Found %d orders in direct array response", len(orders))
		}
	}
}

func TestGetOrderStatistics(t *testing.T) {
	resp, err := http.Get(orderServiceBaseURL + "/api/v1/orders/stats")
	if err != nil {
		t.Fatalf("Failed to get order statistics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var stats OrderStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode statistics response: %v", err)
	}

	// Basic validation - stats should have non-negative values
	if stats.TotalOrders < 0 {
		t.Errorf("Total orders should be non-negative, got %d", stats.TotalOrders)
	}

	if stats.TotalRevenue < 0 {
		t.Errorf("Total revenue should be non-negative, got %.2f", stats.TotalRevenue)
	}

	t.Logf("Order Statistics: Total Orders: %d, Total Revenue: %.2f, Average Order: %.2f", 
		stats.TotalOrders, stats.TotalRevenue, stats.AverageOrder)
}

func TestCreateAndUpdateOrderStatus(t *testing.T) {
	// Create an order first
	orderRequest := CreateOrderRequest{
		CustomerID: "customer456",
		Items: []struct {
			GameID   int     `json:"game_id"`
			GameName string  `json:"game_name"`
			Price    float64 `json:"price"`
			Quantity int     `json:"quantity"`
		}{
			{
				GameID:   3,
				GameName: "Status Test Game",
				Price:    49.99,
				Quantity: 1,
			},
		},
	}

	jsonData, err := json.Marshal(orderRequest)
	if err != nil {
		t.Fatalf("Failed to marshal order request: %v", err)
	}

	resp, err := http.Post(orderServiceBaseURL+"/api/v1/orders", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	defer resp.Body.Close()

	var createResponse CreateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}

	createdOrder := createResponse.Order

	// Update the order status
	statusUpdate := UpdateStatusRequest{
		Status: "shipped",
	}

	updateData, err := json.Marshal(statusUpdate)
	if err != nil {
		t.Fatalf("Failed to marshal status update: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/orders/%s/status", orderServiceBaseURL, createdOrder.ID), bytes.NewBuffer(updateData))
	if err != nil {
		t.Fatalf("Failed to create status update request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	updateResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to update order status: %v", err)
	}
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200 for status update, got %d", updateResp.StatusCode)
	}

	var updateResponse map[string]interface{}
	if err := json.NewDecoder(updateResp.Body).Decode(&updateResponse); err != nil {
		t.Fatalf("Failed to decode update response: %v", err)
	}

	// Verify the update was successful by fetching the order again
	getResp, err := http.Get(fmt.Sprintf("%s/api/v1/orders/%s", orderServiceBaseURL, createdOrder.ID))
	if err != nil {
		t.Fatalf("Failed to get updated order: %v", err)
	}
	defer getResp.Body.Close()

	var getResponse GetOrderResponse
	if err := json.NewDecoder(getResp.Body).Decode(&getResponse); err != nil {
		t.Fatalf("Failed to decode get response: %v", err)
	}

	if getResponse.Order.Status != statusUpdate.Status {
		t.Errorf("Expected status %s, got %s", statusUpdate.Status, getResponse.Order.Status)
	}
}

func TestGetSpecificOrder(t *testing.T) {
	// Create an order first
	orderRequest := CreateOrderRequest{
		CustomerID: "customer789",
		Items: []struct {
			GameID   int     `json:"game_id"`
			GameName string  `json:"game_name"`
			Price    float64 `json:"price"`
			Quantity int     `json:"quantity"`
		}{
			{
				GameID:   4,
				GameName: "Specific Order Test Game",
				Price:    24.99,
				Quantity: 3,
			},
		},
	}

	jsonData, err := json.Marshal(orderRequest)
	if err != nil {
		t.Fatalf("Failed to marshal order request: %v", err)
	}

	resp, err := http.Post(orderServiceBaseURL+"/api/v1/orders", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	defer resp.Body.Close()

	var createResponse CreateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}

	createdOrder := createResponse.Order

	// Get the specific order
	getResp, err := http.Get(fmt.Sprintf("%s/api/v1/orders/%s", orderServiceBaseURL, createdOrder.ID))
	if err != nil {
		t.Fatalf("Failed to get specific order: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", getResp.StatusCode)
	}

	var getResponse GetOrderResponse
	if err := json.NewDecoder(getResp.Body).Decode(&getResponse); err != nil {
		t.Fatalf("Failed to decode get response: %v", err)
	}

	retrievedOrder := getResponse.Order

	if retrievedOrder.ID != createdOrder.ID {
		t.Errorf("Expected order ID %s, got %s", createdOrder.ID, retrievedOrder.ID)
	}

	if retrievedOrder.CustomerID != orderRequest.CustomerID {
		t.Errorf("Expected customer ID %s, got %s", orderRequest.CustomerID, retrievedOrder.CustomerID)
	}
}

func TestGetOrdersByCustomer(t *testing.T) {
	customerID := "customer_test_123"
	
	// Create a couple of orders for the same customer
	for i := 0; i < 2; i++ {
		orderRequest := CreateOrderRequest{
			CustomerID: customerID,
			Items: []struct {
				GameID   int     `json:"game_id"`
				GameName string  `json:"game_name"`
				Price    float64 `json:"price"`
				Quantity int     `json:"quantity"`
			}{
				{
					GameID:   i + 10,
					GameName: fmt.Sprintf("Customer Test Game %d", i+1),
					Price:    19.99,
					Quantity: 1,
				},
			},
		}

		jsonData, err := json.Marshal(orderRequest)
		if err != nil {
			t.Fatalf("Failed to marshal order request %d: %v", i, err)
		}

		resp, err := http.Post(orderServiceBaseURL+"/api/v1/orders", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create order %d: %v", i, err)
		}
		resp.Body.Close()
	}

	// Get orders by customer
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/orders/customer/%s", orderServiceBaseURL, customerID))
	if err != nil {
		t.Fatalf("Failed to get orders by customer: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response GetOrdersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode customer orders response: %v", err)
	}

	orders := response.Orders

	// We should have at least 2 orders for this customer
	if len(orders) < 2 {
		t.Errorf("Expected at least 2 orders for customer, got %d", len(orders))
	}

	// Verify all orders belong to the correct customer
	for _, order := range orders {
		if order.CustomerID != customerID {
			t.Errorf("Expected customer ID %s, got %s", customerID, order.CustomerID)
		}
	}
}

func TestDeleteOrder(t *testing.T) {
	// Create an order first
	orderRequest := CreateOrderRequest{
		CustomerID: "customer_delete_test",
		Items: []struct {
			GameID   int     `json:"game_id"`
			GameName string  `json:"game_name"`
			Price    float64 `json:"price"`
			Quantity int     `json:"quantity"`
		}{
			{
				GameID:   99,
				GameName: "Order To Delete",
				Price:    9.99,
				Quantity: 1,
			},
		},
	}

	jsonData, err := json.Marshal(orderRequest)
	if err != nil {
		t.Fatalf("Failed to marshal order request: %v", err)
	}

	resp, err := http.Post(orderServiceBaseURL+"/api/v1/orders", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	defer resp.Body.Close()

	var createResponse CreateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}

	createdOrder := createResponse.Order

	// Delete the order
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/orders/%s", orderServiceBaseURL, createdOrder.ID), nil)
	if err != nil {
		t.Fatalf("Failed to create delete request: %v", err)
	}

	deleteResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete order: %v", err)
	}
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode != http.StatusOK && deleteResp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code 200 or 204 for delete, got %d", deleteResp.StatusCode)
	}

	// Verify the order is deleted
	getResp, err := http.Get(fmt.Sprintf("%s/api/v1/orders/%s", orderServiceBaseURL, createdOrder.ID))
	if err != nil {
		t.Fatalf("Failed to verify order deletion: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404 after deletion, got %d", getResp.StatusCode)
	}
}

func TestInvalidOrderCreation(t *testing.T) {
	// Test with invalid data (missing required fields)
	invalidRequest := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"game_id": 1,
				// Missing game_name, price, quantity
			},
		},
		// Missing customer_id
	}

	jsonData, err := json.Marshal(invalidRequest)
	if err != nil {
		t.Fatalf("Failed to marshal invalid request: %v", err)
	}

	resp, err := http.Post(orderServiceBaseURL+"/api/v1/orders", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make invalid create request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for invalid request, got %d", resp.StatusCode)
	}
}
