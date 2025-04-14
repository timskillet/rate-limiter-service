# Use official Golang image as builder
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Install Git (required for go mod)
RUN apk add --no-cache git

# Copy go mod and download deps
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -o rate-limiter-service ./cmd/server

# Final image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/rate-limiter-service .

# Expose app port
EXPOSE 8080

# Run the app
CMD ["./rate-limiter-service"]
