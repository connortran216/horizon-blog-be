# Build stage
FROM golang:1.25.3-alpine AS builder

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
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o main .

# Final stage
FROM ubuntu:latest

# Install ca-certificates for HTTPS requests
RUN apt-get update && apt-get install -y ca-certificates tzdata wget golang-go && apt-get clean

# Create app user
RUN groupadd -r appgroup && useradd -r -g appgroup appuser

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy generated swagger docs
COPY --from=builder /app/docs ./docs

# Copy migration files
COPY --from=builder /app/migration ./migration

# Create home directory for appuser
RUN mkdir -p /home/appuser && chown -R appuser:appgroup /home/appuser

# Copy the whole repo to /home/appuser
COPY . /home/appuser

# Change ownership to app user
RUN chown -R appuser:appgroup /root/

# # Switch to non-root user
# USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
