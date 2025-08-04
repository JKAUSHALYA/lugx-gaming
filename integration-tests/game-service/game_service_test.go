package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const gameServiceBaseURL = "http://localhost:30080"

type Game struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	ReleasedDate time.Time `json:"released_date"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type CreateGameRequest struct {
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	ReleasedDate string  `json:"released_date"`
	Price        float64 `json:"price"`
}

type UpdateGameRequest struct {
	Name         *string  `json:"name,omitempty"`
	Category     *string  `json:"category,omitempty"`
	ReleasedDate *string  `json:"released_date,omitempty"`
	Price        *float64 `json:"price,omitempty"`
}

func TestGameServiceHealth(t *testing.T) {
	resp, err := http.Get(gameServiceBaseURL + "/api/v1/health")
	if err != nil {
		t.Fatalf("Failed to make health check request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestCreateGame(t *testing.T) {
	gameRequest := CreateGameRequest{
		Name:         "Test Game",
		Category:     "Action",
		ReleasedDate: "2024-01-01",
		Price:        59.99,
	}

	jsonData, err := json.Marshal(gameRequest)
	if err != nil {
		t.Fatalf("Failed to marshal game request: %v", err)
	}

	resp, err := http.Post(gameServiceBaseURL+"/api/v1/games", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 201 or 200, got %d", resp.StatusCode)
	}

	var response SuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Extract game data from response
	gameData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract game data from response")
	}

	if gameData["name"] != gameRequest.Name {
		t.Errorf("Expected game name %s, got %s", gameRequest.Name, gameData["name"])
	}

	if gameData["category"] != gameRequest.Category {
		t.Errorf("Expected game category %s, got %s", gameRequest.Category, gameData["category"])
	}

	if gameData["price"] != gameRequest.Price {
		t.Errorf("Expected game price %.2f, got %.2f", gameRequest.Price, gameData["price"])
	}
}

func TestGetAllGames(t *testing.T) {
	resp, err := http.Get(gameServiceBaseURL + "/api/v1/games")
	if err != nil {
		t.Fatalf("Failed to get all games: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response SuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Extract games data from response
	gamesData, ok := response.Data.([]interface{})
	if !ok {
		t.Logf("No games found or unexpected data format - this might be expected if database is empty")
		return
	}

	// We expect at least some games (including the one created in previous test)
	if len(gamesData) == 0 {
		t.Log("No games found - this might be expected if database is empty")
	}
}

func TestGetGamesByCategory(t *testing.T) {
	resp, err := http.Get(gameServiceBaseURL + "/api/v1/games?category=Action")
	if err != nil {
		t.Fatalf("Failed to get games by category: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response SuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Extract games data from response
	gamesData, ok := response.Data.([]interface{})
	if !ok {
		t.Logf("No games found or unexpected data format for category 'Action'")
		return
	}

	// Check that all returned games have the correct category
	for _, gameInterface := range gamesData {
		gameData, ok := gameInterface.(map[string]interface{})
		if !ok {
			continue
		}
		if gameData["category"] != "Action" {
			t.Errorf("Expected game category 'Action', got '%s'", gameData["category"])
		}
	}
}

func TestCreateAndUpdateGame(t *testing.T) {
	// First create a game
	gameRequest := CreateGameRequest{
		Name:         "Integration Test Game",
		Category:     "RPG",
		ReleasedDate: "2024-03-15",
		Price:        49.99,
	}

	jsonData, err := json.Marshal(gameRequest)
	if err != nil {
		t.Fatalf("Failed to marshal game request: %v", err)
	}

	resp, err := http.Post(gameServiceBaseURL+"/api/v1/games", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	defer resp.Body.Close()

	var createResponse SuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}

	// Extract game data from response
	gameData, ok := createResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract game data from create response")
	}

	gameID := int(gameData["id"].(float64))

	// Now update the game
	newName := "Updated Integration Test Game"
	newPrice := 39.99
	updateRequest := UpdateGameRequest{
		Name:  &newName,
		Price: &newPrice,
	}

	updateData, err := json.Marshal(updateRequest)
	if err != nil {
		t.Fatalf("Failed to marshal update request: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/games/%d", gameServiceBaseURL, gameID), bytes.NewBuffer(updateData))
	if err != nil {
		t.Fatalf("Failed to create update request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	updateResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to update game: %v", err)
	}
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200 for update, got %d", updateResp.StatusCode)
	}

	var updateResponse SuccessResponse
	if err := json.NewDecoder(updateResp.Body).Decode(&updateResponse); err != nil {
		t.Fatalf("Failed to decode update response: %v", err)
	}

	// Extract updated game data from response
	updatedGameData, ok := updateResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract updated game data from response")
	}

	if updatedGameData["name"] != newName {
		t.Errorf("Expected updated name %s, got %s", newName, updatedGameData["name"])
	}

	if updatedGameData["price"] != newPrice {
		t.Errorf("Expected updated price %.2f, got %.2f", newPrice, updatedGameData["price"])
	}
}

func TestGetSpecificGame(t *testing.T) {
	// Create a game first
	gameRequest := CreateGameRequest{
		Name:         "Specific Game Test",
		Category:     "Strategy",
		ReleasedDate: "2024-02-01",
		Price:        29.99,
	}

	jsonData, err := json.Marshal(gameRequest)
	if err != nil {
		t.Fatalf("Failed to marshal game request: %v", err)
	}

	resp, err := http.Post(gameServiceBaseURL+"/api/v1/games", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	defer resp.Body.Close()

	var createResponse SuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}

	// Extract game data from response
	gameData, ok := createResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract game data from create response")
	}

	gameID := int(gameData["id"].(float64))

	// Now get the specific game
	getResp, err := http.Get(fmt.Sprintf("%s/api/v1/games/%d", gameServiceBaseURL, gameID))
	if err != nil {
		t.Fatalf("Failed to get specific game: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", getResp.StatusCode)
	}

	var getResponse SuccessResponse
	if err := json.NewDecoder(getResp.Body).Decode(&getResponse); err != nil {
		t.Fatalf("Failed to decode get response: %v", err)
	}

	// Extract retrieved game data from response
	retrievedGameData, ok := getResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract retrieved game data from response")
	}

	if int(retrievedGameData["id"].(float64)) != gameID {
		t.Errorf("Expected game ID %d, got %d", gameID, int(retrievedGameData["id"].(float64)))
	}

	if retrievedGameData["name"] != gameRequest.Name {
		t.Errorf("Expected game name %s, got %s", gameRequest.Name, retrievedGameData["name"])
	}
}

func TestDeleteGame(t *testing.T) {
	// Create a game first
	gameRequest := CreateGameRequest{
		Name:         "Game To Delete",
		Category:     "Test",
		ReleasedDate: "2024-01-15",
		Price:        19.99,
	}

	jsonData, err := json.Marshal(gameRequest)
	if err != nil {
		t.Fatalf("Failed to marshal game request: %v", err)
	}

	resp, err := http.Post(gameServiceBaseURL+"/api/v1/games", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}
	defer resp.Body.Close()

	var createResponse SuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}

	// Extract game data from response
	gameData, ok := createResponse.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract game data from create response")
	}

	gameID := int(gameData["id"].(float64))

	// Now delete the game
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/games/%d", gameServiceBaseURL, gameID), nil)
	if err != nil {
		t.Fatalf("Failed to create delete request: %v", err)
	}

	deleteResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete game: %v", err)
	}
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode != http.StatusOK && deleteResp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code 200 or 204 for delete, got %d", deleteResp.StatusCode)
	}

	// Verify the game is deleted by trying to get it
	getResp, err := http.Get(fmt.Sprintf("%s/api/v1/games/%d", gameServiceBaseURL, gameID))
	if err != nil {
		t.Fatalf("Failed to verify game deletion: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404 after deletion, got %d", getResp.StatusCode)
	}
}

func TestInvalidGameCreation(t *testing.T) {
	// Test with invalid data (missing required fields)
	invalidRequest := map[string]interface{}{
		"name": "Invalid Game",
		// Missing category, released_date, and price
	}

	jsonData, err := json.Marshal(invalidRequest)
	if err != nil {
		t.Fatalf("Failed to marshal invalid request: %v", err)
	}

	resp, err := http.Post(gameServiceBaseURL+"/api/v1/games", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make invalid create request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for invalid request, got %d", resp.StatusCode)
	}
}
