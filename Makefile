.PHONY: build run test clean docker-build docker-run docker-stop help

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker Compose services"
	@echo "  deps         - Download dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"

# Build the application
build:
	go build -o bin/swipe-sports-backend .

# Run the application locally
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Build Docker image
docker-build:
	docker build -t swipe-sports-backend .

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose services
docker-stop:
	docker-compose down

# Run with Docker Compose and rebuild
docker-run-build:
	docker-compose up -d --build

# View logs
logs:
	docker-compose logs -f app

# Access database
db:
	docker-compose exec mysql mysql -u root -ppassword swipe_sports

# Access Redis CLI
redis:
	docker-compose exec redis redis-cli

# Generate mock data
mock-data:
	@echo "Mock data is already included in scripts/init.sql"
	@echo "Run 'docker-compose up -d' to start services with sample data"

# Setup development environment
setup: deps docker-run-build
	@echo "Development environment is ready!"
	@echo "API: http://localhost:8080"
	@echo "phpMyAdmin: http://localhost:8081"
	@echo "Health check: http://localhost:8080/health" 