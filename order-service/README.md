# LUGX Gaming - Order Service

A microservice for managing orders in the LUGX Gaming platform. This service handles order creation, tracking, and management with PostgreSQL persistence.

## Features

- **Order Management**: Create, read, update, and delete orders
- **Cart Items**: Support for multiple game items per order
- **Order Tracking**: Status updates (pending, confirmed, processing, shipped, delivered, cancelled)
- **Customer Orders**: Retrieve all orders for a specific customer
- **Order Statistics**: Basic analytics and reporting
- **Database Persistence**: PostgreSQL with automatic table creation
- **RESTful API**: Clean REST endpoints with JSON responses
- **Docker Support**: Containerized deployment with Docker Compose

## API Endpoints

### Orders

- `POST /api/v1/orders` - Create a new order
- `GET /api/v1/orders` - Get all orders (with pagination)
- `GET /api/v1/orders/:id` - Get a specific order
- `PUT /api/v1/orders/:id/status` - Update order status
- `DELETE /api/v1/orders/:id` - Delete an order
- `GET /api/v1/orders/customer/:customer_id` - Get orders by customer
- `GET /api/v1/orders/stats` - Get order statistics

### Health Check

- `GET /health` - Service health check

## Data Models

### Order

```json
{
  "id": "uuid",
  "customer_id": "string",
  "total_price": 0.0,
  "status": "pending|confirmed|processing|shipped|delivered|cancelled",
  "order_date": "2024-01-01T00:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "items": []
}
```

### Order Item

```json
{
  "id": "uuid",
  "order_id": "uuid",
  "game_id": 1,
  "game_name": "Game Name",
  "price": 59.99,
  "quantity": 1,
  "subtotal": 59.99
}
```

## Setup and Installation

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker (optional)

### Local Development

1. **Clone and navigate to the order-service directory**

   ```bash
   cd order-service
   ```

2. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your database configuration
   ```

3. **Install dependencies**

   ```bash
   go mod tidy
   ```

4. **Start PostgreSQL** (if not using Docker)

   - Make sure PostgreSQL is running on port 5432
   - Create a database named `lugx_gaming`

5. **Run the service**
   ```bash
   go run main.go
   ```
   Or use the batch file:
   ```bash
   start.bat
   ```

### Docker Deployment

1. **Using Docker Compose** (recommended)

   ```bash
   docker-compose up -d
   ```

2. **Manual Docker build**
   ```bash
   docker build -t order-service .
   docker run -p 8081:8081 order-service
   ```

## Environment Variables

| Variable      | Description       | Default     |
| ------------- | ----------------- | ----------- |
| `DB_HOST`     | PostgreSQL host   | localhost   |
| `DB_PORT`     | PostgreSQL port   | 5432        |
| `DB_USER`     | Database user     | postgres    |
| `DB_PASSWORD` | Database password | password    |
| `DB_NAME`     | Database name     | lugx_gaming |
| `DB_SSLMODE`  | SSL mode          | disable     |
| `PORT`        | Service port      | 8081        |

## Database Schema

The service automatically creates the following tables:

### orders

- `id` (UUID, Primary Key)
- `customer_id` (VARCHAR)
- `total_price` (DECIMAL)
- `status` (VARCHAR)
- `order_date` (TIMESTAMP)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### order_items

- `id` (UUID, Primary Key)
- `order_id` (UUID, Foreign Key)
- `game_id` (INTEGER)
- `game_name` (VARCHAR)
- `price` (DECIMAL)
- `quantity` (INTEGER)
- `subtotal` (DECIMAL)

## Usage Examples

### Create an Order

```bash
curl -X POST http://localhost:8081/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "customer123",
    "items": [
      {
        "game_id": 1,
        "game_name": "Cyberpunk 2077",
        "price": 59.99,
        "quantity": 1
      },
      {
        "game_id": 2,
        "game_name": "The Witcher 3",
        "price": 39.99,
        "quantity": 2
      }
    ]
  }'
```

### Get Order by ID

```bash
curl http://localhost:8081/api/v1/orders/{order-id}
```

### Update Order Status

```bash
curl -X PUT http://localhost:8081/api/v1/orders/{order-id}/status \
  -H "Content-Type: application/json" \
  -d '{"status": "confirmed"}'
```

### Get Customer Orders

```bash
curl http://localhost:8081/api/v1/orders/customer/customer123
```

## Architecture

The service follows a clean architecture pattern:

```
order-service/
├── main.go                 # Application entry point
├── models/                 # Data models and DTOs
├── handlers/               # HTTP request handlers
├── service/                # Business logic layer
├── repository/             # Data access layer
├── database/               # Database connection and setup
└── routes/                 # Route definitions
```

## Testing

The service includes comprehensive error handling and validation:

- Input validation for all endpoints
- Database transaction safety
- Proper HTTP status codes
- Detailed error messages

## Integration

This service is designed to work alongside the game-service and can be integrated with:

- Authentication services
- Payment processing services
- Inventory management systems
- Notification services

## Contributing

1. Follow Go coding standards
2. Add tests for new features
3. Update documentation
4. Ensure Docker compatibility

## License

This project is part of the LUGX Gaming platform.
