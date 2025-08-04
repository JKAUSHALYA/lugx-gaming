package handlers

import (
	"net/http"
	"strconv"

	"game-service/models"
	"game-service/service"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService *service.GameService
}

// NewGameHandler creates a new game handler
func NewGameHandler() *GameHandler {
	return &GameHandler{
		gameService: service.NewGameService(),
	}
}

// CreateGame handles POST /games
func (h *GameHandler) CreateGame(c *gin.Context) {
	var req models.CreateGameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	game, err := h.gameService.CreateGame(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to create game",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Game created successfully",
		Data:    game,
	})
}

// GetGame handles GET /games/:id
func (h *GameHandler) GetGame(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid game ID",
			Message: "Game ID must be a number",
		})
		return
	}

	game, err := h.gameService.GetGameByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Game not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Game retrieved successfully",
		Data:    game,
	})
}

// GetAllGames handles GET /games
func (h *GameHandler) GetAllGames(c *gin.Context) {
	// Check if category filter is provided
	category := c.Query("category")
	
	var games []*models.Game
	var err error

	if category != "" {
		games, err = h.gameService.GetGamesByCategory(category)
	} else {
		games, err = h.gameService.GetAllGames()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to retrieve games",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Games retrieved successfully",
		Data:    games,
	})
}

// UpdateGame handles PUT /games/:id
func (h *GameHandler) UpdateGame(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid game ID",
			Message: "Game ID must be a number",
		})
		return
	}

	var req models.UpdateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	game, err := h.gameService.UpdateGame(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to update game",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Game updated successfully",
		Data:    game,
	})
}

// DeleteGame handles DELETE /games/:id
func (h *GameHandler) DeleteGame(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid game ID",
			Message: "Game ID must be a number",
		})
		return
	}

	err = h.gameService.DeleteGame(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Failed to delete game",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Game deleted successfully",
	})
}

// HealthCheck handles GET /health
func (h *GameHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "game-service",
		"version": "1.0.0",
	})
}
