# Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s' -o ts-viewer ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/ts-viewer .

# Expose the default port
EXPOSE 8080

# Set default environment variables
ENV HTTP_ADDR=:8080
ENV LOG_LEVEL=info

# Run the application
ENTRYPOINT ["./ts-viewer"]
