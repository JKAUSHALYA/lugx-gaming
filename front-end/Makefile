# Variables
IMAGE_NAME = lugx-gaming-frontend
CONTAINER_NAME = lugx-gaming-frontend
PORT = 3000

# Build the Docker image
build:
	docker build -t $(IMAGE_NAME) .

# Run the container
run:
	docker run -d --name $(CONTAINER_NAME) -p $(PORT):80 $(IMAGE_NAME)

# Stop and remove the container
stop:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

# Remove the Docker image
clean: stop
	docker rmi $(IMAGE_NAME) || true

# Build and run
up: build run

# Stop, clean and rebuild
rebuild: clean build run

# View logs
logs:
	docker logs -f $(CONTAINER_NAME)

# Access container shell
shell:
	docker exec -it $(CONTAINER_NAME) /bin/sh

# Show running containers
ps:
	docker ps

.PHONY: build run stop clean up rebuild logs shell ps
