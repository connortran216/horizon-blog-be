# Build stage
FROM golang:1.24.5-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install swag for swagger generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
RUN swag init -g main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create app user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy generated swagger docs
COPY --from=builder /app/docs ./docs

# Copy migration files
COPY --from=builder /app/migration ./migration

# Change ownership to app user
RUN chown -R appuser:appgroup /root/

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
