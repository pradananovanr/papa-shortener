FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy pre-built binary (build on host first!)
COPY server ./server

# Copy templates and config
COPY templates ./templates
COPY config.yaml .

# Expose port
EXPOSE 8080

# Run
CMD ["./server"]
