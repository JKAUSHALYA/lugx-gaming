package service

import (
	"fmt"
	"time"

	"game-service/models"
	"game-service/repository"
)

type GameService struct {
	repo *repository.GameRepository
}

// NewGameService creates a new game service
func NewGameService() *GameService {
	return &GameService{
		repo: repository.NewGameRepository(),
	}
}

// CreateGame creates a new game
func (s *GameService) CreateGame(req *models.CreateGameRequest) (*models.Game, error) {
	// Validate and parse the release date
	releaseDate, err := time.Parse("2006-01-02", req.ReleasedDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format. Use YYYY-MM-DD: %v", err)
	}

	// Create game object
	game := &models.Game{
		Name:         req.Name,
		Category:     req.Category,
		ReleasedDate: releaseDate,
		Price:        req.Price,
	}

	// Save to database
	createdGame, err := s.repo.CreateGame(game)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %v", err)
	}

	return createdGame, nil
}

// GetGameByID retrieves a game by its ID
func (s *GameService) GetGameByID(id int) (*models.Game, error) {
	game, err := s.repo.GetGameByID(id)
	if err != nil {
		return nil, err
	}
	return game, nil
}

// GetAllGames retrieves all games
func (s *GameService) GetAllGames() ([]*models.Game, error) {
	games, err := s.repo.GetAllGames()
	if err != nil {
		return nil, fmt.Errorf("failed to get games: %v", err)
	}
	return games, nil
}

// UpdateGame updates an existing game
func (s *GameService) UpdateGame(id int, req *models.UpdateGameRequest) (*models.Game, error) {
	// Validate date format if provided
	if req.ReleasedDate != nil {
		_, err := time.Parse("2006-01-02", *req.ReleasedDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format. Use YYYY-MM-DD: %v", err)
		}
	}

	// Validate price if provided
	if req.Price != nil && *req.Price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}

	updatedGame, err := s.repo.UpdateGame(id, req)
	if err != nil {
		return nil, err
	}

	return updatedGame, nil
}

// DeleteGame deletes a game by its ID
func (s *GameService) DeleteGame(id int) error {
	err := s.repo.DeleteGame(id)
	if err != nil {
		return err
	}
	return nil
}

// GetGamesByCategory retrieves games by category
func (s *GameService) GetGamesByCategory(category string) ([]*models.Game, error) {
	games, err := s.repo.GetGamesByCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to get games by category: %v", err)
	}
	return games, nil
}
