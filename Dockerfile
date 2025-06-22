# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Create data directory for SQLite
RUN mkdir -p /root/data

# Expose port
EXPOSE 9001

# Set environment variables
ENV API_PORT=9001
ENV DB_DRIVER=sqlite3
ENV DB_DSN=/root/data/admoai.db
ENV ENVIRONMENT=production

# Run the application
CMD ["./server"] 