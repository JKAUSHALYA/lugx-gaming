package repository

import (
	"database/sql"
	"fmt"
	"time"

	"game-service/database"
	"game-service/models"
)

type GameRepository struct {
	db *sql.DB
}

// NewGameRepository creates a new game repository
func NewGameRepository() *GameRepository {
	return &GameRepository{
		db: database.DB,
	}
}

// CreateGame creates a new game in the database
func (r *GameRepository) CreateGame(game *models.Game) (*models.Game, error) {
	query := `
		INSERT INTO games (name, category, released_date, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	
	now := time.Now()
	game.CreatedAt = now
	game.UpdatedAt = now

	err := r.db.QueryRow(query, game.Name, game.Category, game.ReleasedDate, game.Price, game.CreatedAt, game.UpdatedAt).
		Scan(&game.ID, &game.CreatedAt, &game.UpdatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %v", err)
	}

	return game, nil
}

// GetGameByID retrieves a game by its ID
func (r *GameRepository) GetGameByID(id int) (*models.Game, error) {
	query := `
		SELECT id, name, category, released_date, price, created_at, updated_at
		FROM games
		WHERE id = $1
	`

	game := &models.Game{}
	err := r.db.QueryRow(query, id).Scan(
		&game.ID,
		&game.Name,
		&game.Category,
		&game.ReleasedDate,
		&game.Price,
		&game.CreatedAt,
		&game.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get game: %v", err)
	}

	return game, nil
}

// GetAllGames retrieves all games from the database
func (r *GameRepository) GetAllGames() ([]*models.Game, error) {
	query := `
		SELECT id, name, category, released_date, price, created_at, updated_at
		FROM games
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get games: %v", err)
	}
	defer rows.Close()

	var games []*models.Game
	for rows.Next() {
		game := &models.Game{}
		err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.Category,
			&game.ReleasedDate,
			&game.Price,
			&game.CreatedAt,
			&game.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %v", err)
		}
		games = append(games, game)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate games: %v", err)
	}

	return games, nil
}

// UpdateGame updates an existing game
func (r *GameRepository) UpdateGame(id int, updates *models.UpdateGameRequest) (*models.Game, error) {
	// First, get the current game
	currentGame, err := r.GetGameByID(id)
	if err != nil {
		return nil, err
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *updates.Name)
		argIndex++
		currentGame.Name = *updates.Name
	}

	if updates.Category != nil {
		setParts = append(setParts, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, *updates.Category)
		argIndex++
		currentGame.Category = *updates.Category
	}

	if updates.ReleasedDate != nil {
		releaseDate, err := time.Parse("2006-01-02", *updates.ReleasedDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}
		setParts = append(setParts, fmt.Sprintf("released_date = $%d", argIndex))
		args = append(args, releaseDate)
		argIndex++
		currentGame.ReleasedDate = releaseDate
	}

	if updates.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *updates.Price)
		argIndex++
		currentGame.Price = *updates.Price
	}

	if len(setParts) == 0 {
		return currentGame, nil // No updates to perform
	}

	// Add updated_at to the update
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add ID for WHERE clause
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE games 
		SET %s
		WHERE id = $%d
		RETURNING updated_at
	`, fmt.Sprintf("%s", setParts[0]), argIndex)

	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf(`
			UPDATE games 
			SET %s
			WHERE id = $%d
			RETURNING updated_at
		`, fmt.Sprintf("%s, %s", setParts[0], setParts[i]), argIndex)
	}

	// Reconstruct the query properly
	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query = fmt.Sprintf(`
		UPDATE games 
		SET %s
		WHERE id = $%d
		RETURNING updated_at
	`, setClause, argIndex)

	err = r.db.QueryRow(query, args...).Scan(&currentGame.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update game: %v", err)
	}

	return currentGame, nil
}

// DeleteGame deletes a game by its ID
func (r *GameRepository) DeleteGame(id int) error {
	query := `DELETE FROM games WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete game: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("game with ID %d not found", id)
	}

	return nil
}

// GetGamesByCategory retrieves games by category
func (r *GameRepository) GetGamesByCategory(category string) ([]*models.Game, error) {
	query := `
		SELECT id, name, category, released_date, price, created_at, updated_at
		FROM games
		WHERE category = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get games by category: %v", err)
	}
	defer rows.Close()

	var games []*models.Game
	for rows.Next() {
		game := &models.Game{}
		err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.Category,
			&game.ReleasedDate,
			&game.Price,
			&game.CreatedAt,
			&game.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %v", err)
		}
		games = append(games, game)
	}

	return games, nil
}
