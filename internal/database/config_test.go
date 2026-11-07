package database

import "testing"

func TestPoolConfigFromEnv(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "db")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_DB", "fitness")
	t.Setenv("POSTGRES_USER", "admin")
	t.Setenv("POSTGRES_PASSWORD", "secret")

	config, err := PoolConfigFromEnv()
	if err != nil {
		t.Fatalf("PoolConfigFromEnv() error = %v", err)
	}

	if config.ConnConfig.Host != "db" {
		t.Errorf("host = %q, want %q", config.ConnConfig.Host, "db")
	}
	if config.ConnConfig.Port != 5432 {
		t.Errorf("port = %d, want %d", config.ConnConfig.Port, 5432)
	}
	if config.ConnConfig.Database != "fitness" {
		t.Errorf("database = %q, want %q", config.ConnConfig.Database, "fitness")
	}
	if config.ConnConfig.User != "admin" {
		t.Errorf("user = %q, want %q", config.ConnConfig.User, "admin")
	}
	if config.ConnConfig.Password != "secret" {
		t.Errorf("password was not copied to the connection config")
	}
}

func TestPoolConfigFromEnvRejectsInvalidPort(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "db")
	t.Setenv("POSTGRES_PORT", "invalid")
	t.Setenv("POSTGRES_DB", "fitness")
	t.Setenv("POSTGRES_USER", "admin")
	t.Setenv("POSTGRES_PASSWORD", "secret")

	if _, err := PoolConfigFromEnv(); err == nil {
		t.Fatal("PoolConfigFromEnv() error = nil, want an invalid port error")
	}
}
