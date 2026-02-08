# Go API

A simple RESTful API for managing products built with Go, SQLite, sqlc, and goose.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- [Go](https://go.dev/dl/) 1.25.4 or higher
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) - SQL code generator
- [goose](https://github.com/pressly/goose) - Database migration tool

### Installing sqlc

```bash
# macOS
brew install sqlc

# Linux
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Or download binary from: https://github.com/sqlc-dev/sqlc/releases
```

### Installing goose

```bash
# Using go install
go install github.com/pressly/goose/v3/cmd/goose@latest

# Or using brew (macOS)
brew install goose
```

## Project Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd shopping-list
```

### 2. Install Go Dependencies

```bash
go mod download
```

### 3. Configure Environment

Copy the example environment file:

```bash
cp .env.example .env
```

The `.env` file configures goose with the following settings:

- `GOOSE_DRIVER=sqlite3` - Database driver
- `GOOSE_DBSTRING=./products.db` - Database file path
- `GOOSE_MIGRATION_DIR=./internal/adapters/sqlite/migrations` - Migration files directory

### 4. Run Database Migrations

The project uses goose for database migrations. With the `.env` file configured, you can run goose commands without passing arguments:

```bash
# Load environment variables (for current shell session)
source .env

# Run migrations (creates products.db database)
goose up

# Check migration status
goose status
```

**Available goose commands:**

- `goose up` - Run all pending migrations
- `goose down` - Roll back the last migration
- `goose status` - Show migration status
- `goose reset` - Roll back all migrations

**Note:** If you prefer not to use the `.env` file, you can still pass arguments directly:

```bash
goose -dir internal/adapters/sqlite/migrations sqlite3 products.db up
```

### 5. Generate Go Code from SQL

The project uses sqlc to generate type-safe Go code from SQL queries.

```bash
# Generate Go code from SQL queries
sqlc generate
```

This will generate the following files in `internal/adapters/sqlite/sqlc/`:

- `db.go` - Database interface
- `models.go` - Go structs for database tables
- `queries.sql.go` - Type-safe query functions

### 6. Run the Application

```bash
# Run the API server
go run ./cmd

# Or build and run
go build -o bin/api ./cmd
./bin/api
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check

```http
GET /health
```

### Products

#### List all products

```http
GET /products
```

#### Get product by ID

```http
GET /products/{id}
```

#### Create product

```http
POST /products
Content-Type: application/json

{
  "id": "product-123",
  "name": "Apple",
  "price_in_cents": 250,
  "quantity": 10,
  "created_at": 1738944000
}
```

## Testing

You can test the API using the provided HTTP requests in `http/products.http`. These can be executed with:

- [Kulala](https://kulala.mwco.app/) (Neovim plugin)
- [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) (VS Code extension)
- cURL or any HTTP client

Example with cURL:

```bash
# Health check
curl http://localhost:8080/health

# List products
curl http://localhost:8080/products

# Create product
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "id": "product-123",
    "name": "Apple",
    "price_in_cents": 250,
    "quantity": 10,
    "created_at": 1738944000
  }'
```

## Development Workflow

### Adding New Database Tables

1. Make sure your `.env` is loaded:

```bash
source .env
```

2. Create a new migration file:

```bash
goose create add_new_table sql
```

3. Edit the generated migration file in `internal/adapters/sqlite/migrations/`

4. Run the migration:

```bash
goose up
```

### Adding New Queries

1. Add your SQL query to `internal/adapters/sqlite/sqlc/queries.sql`

2. Regenerate Go code:

```bash
sqlc generate
```

3. Use the generated functions in your Go code

### Rolling Back Migrations

```bash
# Make sure your .env is loaded
source .env

# Roll back the last migration
goose down

# Roll back all migrations
goose reset
```

## Common Issues

### sqlc: command not found

Make sure sqlc is installed and in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### goose: command not found

Make sure goose is installed and in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Database locked error

If you get a "database is locked" error, make sure no other process is accessing `products.db`.

## Tech Stack

- **Go** - Programming language
- **Chi** - HTTP router
- **SQLite** - Database
- **sqlc** - Type-safe SQL code generation
- **goose** - Database migrations
- **slog** - Structured logging
