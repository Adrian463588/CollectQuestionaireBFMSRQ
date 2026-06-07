# Stage 1: Build the Go binary securely
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary with stripped debug symbols (-w -s) for security and size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/server/main.go

# Stage 2: Minimal runtime image
FROM alpine:3.20

# Add non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install CA certificates for secure HTTPS outbound connections
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Set ownership to non-root user
RUN chown appuser:appgroup /app/server

# Expose the configured port
EXPOSE 8081

# Run as non-root user
USER appuser

# Start application
CMD ["./server"]
