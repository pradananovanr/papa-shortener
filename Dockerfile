# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files first to leverage Docker cache
COPY go.mod ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# Stage 2: Final lightweight image
FROM alpine:3.19

WORKDIR /app

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Copy built binary from builder stage
COPY --from=builder /app/server .

# Copy templates (dynamic HTML files)
COPY templates ./templates

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
