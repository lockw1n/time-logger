# ==============================
# Stage 1 — dependencies
# ==============================
FROM golang:1.25.2-alpine AS deps
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

# ==============================
# Stage 2 — build application
# ==============================
FROM deps AS builder
WORKDIR /app
COPY . .
RUN go build -o timelogger-server ./cmd/server

# ==============================
# Stage 3 — runtime
# ==============================
FROM alpine:latest AS runtime
WORKDIR /app
RUN apk add --no-cache ca-certificates wget
COPY --from=builder /app/timelogger-server .
RUN chmod +x timelogger-server
HEALTHCHECK --interval=5s --timeout=2s --retries=10 \
  CMD wget --spider -q http://localhost:8080/health || exit 1
EXPOSE 8080
CMD ["./timelogger-server"]
