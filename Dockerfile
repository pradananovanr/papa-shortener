FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy source code
COPY . .

# Build with auto toolchain
ENV GOTOOLCHAIN=auto
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/server .

# Copy templates and config
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/config.yaml .

# Expose port
EXPOSE 8080

# Run
CMD ["./server"]