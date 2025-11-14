package main

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("Connecting to database...")

	// Database
	conn, err := sql.Open("sqlite3", "products.db")
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
