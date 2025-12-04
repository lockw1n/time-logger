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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o timelogger-server ./cmd/server

# ==============================
# Stage 2 — runtime
# ==============================
FROM debian:bookworm-slim AS runtime

WORKDIR /app

# Install wkhtmltopdf + required fonts
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
       ca-certificates \
       curl \
       wkhtmltopdf \
       fonts-dejavu \
       fontconfig \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/timelogger-server .

RUN chmod +x timelogger-server

HEALTHCHECK --interval=5s --timeout=2s --retries=10 \
  CMD curl -fsS http://localhost:8080/health || exit 1

EXPOSE 8080
CMD ["./timelogger-server"]
