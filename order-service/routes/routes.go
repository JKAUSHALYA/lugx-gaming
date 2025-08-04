package routes

import (
	"order-service/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler()

	// Health check endpoint
	router.GET("/health", orderHandler.HealthCheck)

	// API version 1 routes
	v1 := router.Group("/api/v1")
	{
		// Order routes
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)                                 // Create new order
			orders.GET("", orderHandler.GetAllOrders)                               // Get all orders with pagination
			orders.GET("/stats", orderHandler.GetOrderStatistics)                   // Get order statistics
			orders.GET("/:id", orderHandler.GetOrderByID)                          // Get specific order
			orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)              // Update order status
			orders.DELETE("/:id", orderHandler.DeleteOrder)                        // Delete order
			orders.GET("/customer/:customer_id", orderHandler.GetOrdersByCustomerID) // Get orders by customer
		}
	}

	return router
}
