# Build stage
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go

# Run stage — use a minimal image
FROM alpine:3.19

# Set working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/server .

# Copy migrations so they can be run inside the container
COPY --from=builder /app/migrations ./migrations

# Expose the API port
EXPOSE 8080

# Run the binary
CMD ["./server"]