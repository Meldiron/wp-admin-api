FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install necessary build tools for CGO
RUN apk add --no-cache gcc musl-dev

# Set environment variables for CGO
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application with explicit CGO flags
RUN go build -tags "linux" -o main ./src

FROM alpine:latest

WORKDIR /app

# Install runtime dependencies for SQLite, healthcheck, and Docker CLI for various actions
RUN apk add --no-cache gcc musl-dev curl docker-cli


# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/src/index.html ./src/index.html
COPY --from=builder /app/public ./public

# Create directory for mock data
RUN mkdir -p /app/mock

# Expose the port
EXPOSE 3000

# Run the application
CMD ["./main"]