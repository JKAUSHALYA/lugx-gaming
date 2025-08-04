# Lugx Gaming Frontend Docker

This directory contains the Docker configuration for the Lugx Gaming frontend application.

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Build and start the frontend
docker-compose up --build

# Run in background
docker-compose up -d --build

# Stop the service
docker-compose down
```

### Using Docker Commands

```bash
# Build the image
docker build -t lugx-gaming-frontend .

# Run the container
docker run -d --name lugx-gaming-frontend -p 3000:80 lugx-gaming-frontend

# Stop and remove
docker stop lugx-gaming-frontend
docker rm lugx-gaming-frontend
```

### Using Makefile

```bash
# Build and run
make up

# Stop and clean
make clean

# Rebuild everything
make rebuild

# View logs
make logs
```

## Access

Once running, the frontend will be available at:

- http://localhost:3000

## Files

- `Dockerfile`: Docker image configuration
- `docker-compose.yml`: Docker Compose service configuration
- `.dockerignore`: Files to exclude from Docker build context
- `Makefile`: Convenient build and run commands

## Configuration

The frontend is served using nginx on port 80 inside the container, mapped to port 3000 on the host.

To change the port, modify the docker-compose.yml file:

```yaml
ports:
  - "8080:80" # Change 8080 to your desired port
```
