package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/KoIvRo/backend/internal/database"
	"github.com/KoIvRo/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	config, err := database.PoolConfigFromEnv()
	if err != nil {
		return fmt.Errorf("read database configuration: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("create database pool: %w", err)
	}
	defer pool.Close()

	if err := seed(ctx, pool); err != nil {
		return err
	}

	return nil
}

func seed(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "LOCK TABLE classes IN EXCLUSIVE MODE"); err != nil {
		return fmt.Errorf("lock classes table: %w", err)
	}

	var hasClasses bool
	if err := tx.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM classes)").Scan(&hasClasses); err != nil {
		return fmt.Errorf("check existing classes: %w", err)
	}
	if hasClasses {
		fmt.Println("seed skipped: database already contains classes")
		return tx.Commit(ctx)
	}

	location := time.FixedZone("VLAT", 10*60*60)
	classes := []domain.Class{
		{Name: "Йога", Trainer: "Анна", StartsAt: time.Date(2026, 11, 10, 18, 30, 0, 0, location), Capacity: 15},
		{Name: "Силовая тренировка", Trainer: "Иван", StartsAt: time.Date(2026, 11, 11, 19, 0, 0, 0, location), Capacity: 10},
		{Name: "Пилатес", Trainer: "Мария", StartsAt: time.Date(2026, 11, 12, 17, 0, 0, 0, location), Capacity: 12},
	}

	classIDs := make([]int, 0, len(classes))
	for _, class := range classes {
		var id int
		err := tx.QueryRow(
			ctx,
			`INSERT INTO classes (name, trainer, starts_at, capacity)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			class.Name,
			class.Trainer,
			class.StartsAt,
			class.Capacity,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("seed class %q: %w", class.Name, err)
		}

		classIDs = append(classIDs, id)
	}

	signups := []struct {
		classIndex int
		client     string
		phone      string
	}{
		{0, "Иванов", "+79990000001"},
		{0, "Петров", "+79990000002"},
		{0, "Сидорова", "+79990000003"},
		{1, "Смирнов", "+79990000004"},
		{1, "Кузнецова", "+79990000005"},
		{2, "Попов", "+79990000006"},
	}

	for _, signup := range signups {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO signups (class_id, client, phone) VALUES ($1, $2, $3)`,
			classIDs[signup.classIndex],
			signup.client,
			signup.phone,
		)
		if err != nil {
			return fmt.Errorf("seed signup for %q: %w", signup.client, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	fmt.Printf("seeded %d classes and %d signups\n", len(classes), len(signups))

	return nil
}
