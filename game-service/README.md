# Game Service API

A REST API microservice for managing games built with Go, Gin framework, and PostgreSQL.

## Features

- Create, read, update, and delete games
- Store game information: name, category, release date, and price
- PostgreSQL database integration
- RESTful API endpoints
- Docker support
- Input validation
- Error handling

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Docker (optional)

## Installation & Setup

### Method 1: Direct Run

1. **Install dependencies:**

   ```bash
   go mod tidy
   ```

2. **Set up environment variables:**
   Copy `.env` file and update the database configuration:

   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=lugx_gaming
   DB_SSLMODE=disable
   PORT=8080
   ```

3. **Run the service:**
   ```bash
   go run main.go
   ```

### Method 2: Docker Compose

1. **Start services:**
   ```bash
   docker-compose up -d
   ```

This will start both PostgreSQL database and the game service.

## API Endpoints

### Base URL: `http://localhost:8080/api/v1`

### Health Check

- **GET** `/health`
  - Returns service health status

### Game Management

#### Create Game

- **POST** `/games`
- **Request Body:**
  ```json
  {
    "name": "The Witcher 3",
    "category": "RPG",
    "released_date": "2015-05-19",
    "price": 29.99
  }
  ```

#### Get All Games

- **GET** `/games`
- **Query Parameters:**
  - `category` (optional): Filter by category

#### Get Game by ID

- **GET** `/games/{id}`

#### Update Game

- **PUT** `/games/{id}`
- **Request Body:** (all fields optional)
  ```json
  {
    "name": "Updated Game Name",
    "category": "Action",
    "released_date": "2024-01-01",
    "price": 39.99
  }
  ```

#### Delete Game

- **DELETE** `/games/{id}`

## Example API Calls

### Create a new game

```bash
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cyberpunk 2077",
    "category": "RPG",
    "released_date": "2020-12-10",
    "price": 59.99
  }'
```

### Get all games

```bash
curl http://localhost:8080/api/v1/games
```

### Get games by category

```bash
curl http://localhost:8080/api/v1/games?category=RPG
```

### Get specific game

```bash
curl http://localhost:8080/api/v1/games/1
```

### Update a game

```bash
curl -X PUT http://localhost:8080/api/v1/games/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 49.99
  }'
```

### Delete a game

```bash
curl -X DELETE http://localhost:8080/api/v1/games/1
```

## Response Format

### Success Response

```json
{
  "message": "Operation successful",
  "data": { ... }
}
```

### Error Response

```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

## Database Schema

### Games Table

```sql
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    released_date DATE NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Project Structure

```
game-service/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies
├── .env                    # Environment variables
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose configuration
├── models/
│   └── game.go            # Data models
├── database/
│   └── connection.go      # Database connection and setup
├── repository/
│   └── game_repository.go # Data access layer
├── service/
│   └── game_service.go    # Business logic layer
├── handlers/
│   └── game_handler.go    # HTTP request handlers
└── routes/
    └── routes.go          # Route definitions
```

## Error Handling

The API includes comprehensive error handling for:

- Invalid request formats
- Database connection issues
- Resource not found
- Validation errors
- Internal server errors

## CORS Support

The API includes CORS headers to allow cross-origin requests from frontend applications.

## Logging

The service includes structured logging for:

- Database connection status
- Server startup information
- Error tracking

## Development

To add new features:

1. Define new models in `models/`
2. Add repository methods in `repository/`
3. Implement business logic in `service/`
4. Create handlers in `handlers/`
5. Define routes in `routes/`

## Testing

You can test the API using:

- cURL commands (examples above)
- Postman
- Any HTTP client

## License

This project is created for educational purposes as part of the CMM 707 Cloud Computing coursework.
