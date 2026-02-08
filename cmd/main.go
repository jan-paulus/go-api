package main

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Get database path from environment or use default
	dbPath := getEnv("DB_PATH", "products.db")

	logger.Info("Connecting to database...", "path", dbPath)

	// Database
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error("Failed to connect to database.", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	logger.Info("Successfully connected to database.")

	api := application{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Failed to start server.", "error", err)
		os.Exit(1)
	}
}
