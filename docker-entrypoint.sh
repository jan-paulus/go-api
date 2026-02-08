#!/bin/sh
set -e

echo "========================================"
echo "Product API - Starting Container"
echo "========================================"

# Set database path
export DB_PATH="${DB_PATH:-/data/products.db}"

# Set goose environment variables for migrations
export GOOSE_DRIVER="sqlite3"
export GOOSE_DBSTRING="${DB_PATH}"
export GOOSE_MIGRATION_DIR="/app/migrations"

echo "Database path: ${DB_PATH}"
echo "Migration directory: ${GOOSE_MIGRATION_DIR}"
echo ""

# Ensure data directory exists
DATA_DIR=$(dirname "${DB_PATH}")
if [ ! -d "${DATA_DIR}" ]; then
    echo "Creating data directory: ${DATA_DIR}"
    mkdir -p "${DATA_DIR}"
fi

# Run database migrations using goose with environment variables
echo "Running database migrations..."
cd "${GOOSE_MIGRATION_DIR}"
if goose up; then
    echo "Migrations completed successfully"
else
    echo "ERROR: Migrations failed!"
    exit 1
fi
cd /app

echo ""
echo "Starting application..."
echo "========================================"
echo ""

# Execute the application (passed as CMD or argument)
exec "$@"
