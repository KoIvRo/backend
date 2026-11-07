package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PoolConfigFromEnv() (*pgxpool.Config, error) {
	host := os.Getenv("POSTGRES_HOST")
	portValue := os.Getenv("POSTGRES_PORT")
	databaseName := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	if host == "" || portValue == "" || databaseName == "" || user == "" || password == "" {
		return nil, errors.New("database environment variables are not fully set")
	}

	port, err := strconv.ParseUint(portValue, 10, 16)
	if err != nil || port == 0 {
		return nil, fmt.Errorf("invalid POSTGRES_PORT: %q", portValue)
	}

	config, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("create database config: %w", err)
	}

	config.ConnConfig.Host = host
	config.ConnConfig.Port = uint16(port)
	config.ConnConfig.Database = databaseName
	config.ConnConfig.User = user
	config.ConnConfig.Password = password

	return config, nil
}
