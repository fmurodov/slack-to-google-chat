# Stage 1: Build the Go application
FROM golang:latest AS builder

WORKDIR /app

COPY . .

# Build the Go application statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Stage 2: Create a minimal container to run the application
FROM scratch

WORKDIR /app

# Copy the binary and ca-certs from the builder stage
COPY --from=builder /app/app .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose the port your application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
