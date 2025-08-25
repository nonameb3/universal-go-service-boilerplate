# Multi-stage build for Universal Go Service
# Production-ready Docker image with minimal footprint

# Build stage
FROM golang:1.21-alpine AS builder

# Install ca-certificates and git for downloading dependencies
RUN apk add --no-cache ca-certificates git

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o universal-service cmd/server/main.go

# Final stage - minimal image
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/universal-service .
COPY --from=builder /app/config ./config

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default environment
ENV GO_ENV=production
ENV LOG_LEVEL=info

# Run the service
CMD ["./universal-service"]