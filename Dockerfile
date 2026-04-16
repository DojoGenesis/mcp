# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Create default decisions directory
RUN mkdir -p /app/decisions

# Copy the binary from builder
COPY --from=builder /build/server .

# Run as non-root user
RUN adduser -D -u 1000 dojo
USER dojo

# Default env vars
ENV DOJO_ADR_PATH=/app/decisions

# Run the server
ENTRYPOINT ["./server"]
