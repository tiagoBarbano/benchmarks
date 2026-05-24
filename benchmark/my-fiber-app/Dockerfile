# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o my-fiber-app ./cmd/server/main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app


# Copy the built binary from builder
COPY --from=builder /app/my-fiber-app .

# Copy .env file
COPY .env .

# Expose port 8080
EXPOSE 8080

# Use a non-root user for security
RUN adduser -D appuser
USER appuser

# Start the application
ENTRYPOINT ["./my-fiber-app"]
