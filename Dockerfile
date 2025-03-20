# Use Golang image
FROM golang:1.20-alpine

# Set working directory
WORKDIR /app

# Copy Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application files
COPY . .

# Build the application
RUN go build -o mock-api main.go

# Set environment variables
ENV PORT=8080
ENV MOCK_DATA_PATH=/app/data

# Expose port
EXPOSE 8080

# Run the application
CMD ["./mock-api"]
