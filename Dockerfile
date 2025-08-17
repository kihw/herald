# Multi-stage build for frontend
FROM node:18 as frontend-builder
WORKDIR /app/web

# Copy package files
COPY web/package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY web/ ./

# Build frontend
RUN npm run build

# Go backend build stage
FROM golang:1.23-alpine AS go-builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Copy the binary from go-builder stage
COPY --from=go-builder /app/server .

# Copy built frontend from frontend-builder stage
COPY --from=frontend-builder /app/web/dist ./web/dist

# Create non-root user for security
RUN adduser -D -s /bin/sh app && \
  chown -R app:app /app
USER app

# Expose port
EXPOSE 8000

# Add health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8000/api/health || exit 1

# Run server
CMD ["./server"]