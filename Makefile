# Load environment variables from .env
include .env
export $(shell sed 's/=.*//' .env)

# Run the application
run:
	@go run main.go

# Build the application
build:
	@go build -o mock-api main.go

# Run application using Docker
docker-run:
	@docker-compose up --build

# Clean build files
clean:
	@rm -rf mock-api
