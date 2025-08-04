package models

import (
	"time"
)

// Game represents a game entity
type Game struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name" binding:"required"`
	Category     string    `json:"category" db:"category" binding:"required"`
	ReleasedDate time.Time `json:"released_date" db:"released_date" binding:"required"`
	Price        float64   `json:"price" db:"price" binding:"required,min=0"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateGameRequest represents the request body for creating a game
type CreateGameRequest struct {
	Name         string  `json:"name" binding:"required"`
	Category     string  `json:"category" binding:"required"`
	ReleasedDate string  `json:"released_date" binding:"required"` // Format: "2006-01-02"
	Price        float64 `json:"price" binding:"required,min=0"`
}

// UpdateGameRequest represents the request body for updating a game
type UpdateGameRequest struct {
	Name         *string  `json:"name,omitempty"`
	Category     *string  `json:"category,omitempty"`
	ReleasedDate *string  `json:"released_date,omitempty"` // Format: "2006-01-02"
	Price        *float64 `json:"price,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
