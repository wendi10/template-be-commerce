# ─── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Download dependencies first (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/bin/api ./cmd/api

# ─── Stage 2: Runtime ──────────────────────────────────────────────────────────
FROM scratch

# Copy timezone data and CA certs from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy migrations
COPY --from=builder /app/migrations /migrations

# Copy the binary
COPY --from=builder /app/bin/api /api

EXPOSE 8080

ENTRYPOINT ["/api"]
