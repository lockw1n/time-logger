# ==============================
# Stage 1 — build application
# ==============================
FROM golang:1.25.2 AS builder

WORKDIR /app

# Install git for private modules
RUN apt-get update && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build server binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o timelogger-server ./cmd/server

# Build migration binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o timelogger-migrate ./cmd/migrate


# ==============================
# Stage 2 — runtime
# ==============================
FROM alpine:3.20 AS runtime

WORKDIR /app

# Install only what a static Go binary actually needs
RUN apk add --no-cache ca-certificates

# Copy both binaries from builder
COPY --from=builder /app/timelogger-server .
COPY --from=builder /app/timelogger-migrate .

RUN chmod +x timelogger-server timelogger-migrate

EXPOSE 8080

# Default command (can be overridden by docker-compose)
CMD ["./timelogger-server"]
