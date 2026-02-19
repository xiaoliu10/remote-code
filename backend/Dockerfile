# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache \
    tmux \
    ca-certificates \
    bash

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .

# Create workspace directory
RUN mkdir -p /workspace

EXPOSE 8080

CMD ["./server"]
