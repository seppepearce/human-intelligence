# Multi-stage build for Go backend with Neo4j integration
FROM golang:1.21-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git ca-certificates tzdata curl

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main_neo4j.go

# Development stage
FROM golang:1.21-alpine AS development

# Install development tools and dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    curl \
    bash \
    make \
    gcc \
    musl-dev

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@latest

# Install additional Go tools for development
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Change ownership of go modules cache
RUN chown -R appuser:appgroup /go/pkg/mod

# Copy source code
COPY . .

# Create necessary directories
RUN mkdir -p /app/logs /app/tmp /app/uploads

# Change ownership of app directory
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080 9090

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Default command for development
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:3.18 AS production

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/server .

# Copy configuration files
COPY --from=builder /app/migrations ./migrations/

# Create necessary directories
RUN mkdir -p /app/logs /app/tmp /app/uploads

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080 9090

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./server"]
