# Multi-stage Dockerfile for Go microtester
# Stage 1: Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Copy go mod files first for better caching
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./
COPY internals internals/

# Build the application
# CGO_ENABLED=0 for static binary
# GOOS=linux for Linux target
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o microtester .

# Stage 2: Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/microtester .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose default port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./microtester"]