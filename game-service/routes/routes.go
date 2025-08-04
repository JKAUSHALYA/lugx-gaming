package routes

import (
	"game-service/handlers"

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
	gameHandler := handlers.NewGameHandler()

	// API version 1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", gameHandler.HealthCheck)

		// Game routes
		games := v1.Group("/games")
		{
			games.POST("", gameHandler.CreateGame)           // Create a new game
			games.GET("", gameHandler.GetAllGames)           // Get all games (with optional category filter)
			games.GET("/:id", gameHandler.GetGame)           // Get game by ID
			games.PUT("/:id", gameHandler.UpdateGame)        // Update game by ID
			games.DELETE("/:id", gameHandler.DeleteGame)     // Delete game by ID
		}
	}

	return router
}
