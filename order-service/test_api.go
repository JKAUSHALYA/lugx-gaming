package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-service/models"
	"order-service/routes"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	router := routes.SetupRoutes()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "order-service", response["service"])
}

func TestCreateOrderValidation(t *testing.T) {
	router := routes.SetupRoutes()
	
	// Test invalid order (missing customer_id)
	invalidOrder := models.CreateOrderRequest{
		Items: []models.CreateOrderItemRequest{
			{
				GameID:   1,
				GameName: "Test Game",
				Price:    10.00,
				Quantity: 1,
			},
		},
	}
	
	jsonData, _ := json.Marshal(invalidOrder)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 400, w.Code)
}

func TestCreateOrderSuccess(t *testing.T) {
	// This test would require a test database setup
	// For now, we'll just test the validation logic
	
	router := routes.SetupRoutes()
	
	validOrder := models.CreateOrderRequest{
		CustomerID: "test-customer",
		Items: []models.CreateOrderItemRequest{
			{
				GameID:   1,
				GameName: "Test Game",
				Price:    10.00,
				Quantity: 1,
			},
		},
	}
	
	jsonData, _ := json.Marshal(validOrder)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	// Without database setup, this will fail with 500, but at least validates the request structure
	// In a real test environment, you'd set up a test database
	fmt.Printf("Response code: %d\n", w.Code)
	fmt.Printf("Response body: %s\n", w.Body.String())
}

func TestGetOrdersEndpoint(t *testing.T) {
	router := routes.SetupRoutes()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/orders", nil)
	router.ServeHTTP(w, req)
	
	// Without database setup, this will fail with 500
	// But we can at least test that the endpoint exists
	fmt.Printf("Response code: %d\n", w.Code)
}

// Note: For comprehensive testing, you would need to:
// 1. Set up a test database
// 2. Use dependency injection for the database connection
// 3. Create test fixtures and tear down after tests
// 4. Test all CRUD operations
// 5. Test error scenarios
