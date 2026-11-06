package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/KoIvRo/backend/internal/database"
	"github.com/KoIvRo/backend/migrations"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const migrationsDirectory = "."

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: go run ./cmd/migrate <up|down|status>")
	}

	config, err := database.PoolConfigFromEnv()
	if err != nil {
		return fmt.Errorf("read database configuration: %w", err)
	}

	db := stdlib.OpenDB(*config.ConnConfig)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	goose.SetBaseFS(migrations.Files)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set migration dialect: %w", err)
	}

	if err := runCommand(args[0], db); err != nil {
		return fmt.Errorf("run migration command %q: %w", args[0], err)
	}

	return nil
}

func runCommand(command string, db *sql.DB) error {
	switch command {
	case "up":
		return goose.Up(db, migrationsDirectory)
	case "down":
		return goose.Down(db, migrationsDirectory)
	case "status":
		return goose.Status(db, migrationsDirectory)
	default:
		return fmt.Errorf("unknown command %q; expected up, down, or status", command)
	}
}
