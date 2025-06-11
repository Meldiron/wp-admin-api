FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./src

FROM alpine:latest

WORKDIR /app

# Install dependencies for SQLite
RUN apk add --no-cache ca-certificates libc6-compat

# Copy the binary from the builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/src/index.html ./src/index.html
COPY --from=builder /app/public ./public

# Create directory for mock data
RUN mkdir -p /mnt/mock

# Expose the port
EXPOSE 3000

# Run the application
CMD ["./server"]