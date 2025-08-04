package main

import (
	"log"
	"os"

	"game-service/database"
	"game-service/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Setup routes
	router := routes.SetupRoutes()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Game Service starting on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET    /api/v1/health")
	log.Printf("  POST   /api/v1/games")
	log.Printf("  GET    /api/v1/games")
	log.Printf("  GET    /api/v1/games/:id")
	log.Printf("  PUT    /api/v1/games/:id")
	log.Printf("  DELETE /api/v1/games/:id")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
