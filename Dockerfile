FROM golang:1.25-alpine AS builder

RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    git

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /build

RUN go mod download

COPY sqlc.yaml ./

COPY internal/adapters/sqlite/migrations ./internal/adapters/sqlite/migrations
COPY internal/adapters/sqlite/sqlc/queries.sql ./internal/adapters/sqlite/sqlc/queries.sql

RUN sqlc generate

COPY cmd ./cmd
COPY internal ./internal

# Build the application
# CGO_ENABLED=1 is required for go-sqlite3
# -ldflags="-w -s" strips debug info for smaller binary
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /build/api \
    ./cmd

FROM alpine:latest

RUN apk add --no-cache \
    sqlite-libs \
    ca-certificates \
    tzdata

WORKDIR /app
RUN mkdir -p /data

COPY --from=builder /build/api /app/api

# Copy goose binary for migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy migration files
COPY --from=builder /build/internal/adapters/sqlite/migrations /app/migrations

# Copy entrypoint script
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

# Expose application port
EXPOSE 8080

# Set default environment variables
ENV DB_PATH=/data/products.db \
    GOOSE_DRIVER=sqlite3 \
    GOOSE_MIGRATION_DIR=/app/migrations

# Use entrypoint script to run migrations before starting app
ENTRYPOINT ["/app/docker-entrypoint.sh"]

# Default command (can be overridden)
CMD ["/app/api"]
