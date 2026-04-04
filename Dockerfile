# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from backend directory
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the backend source
COPY backend/ .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go

# Run stage — use a minimal image
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
