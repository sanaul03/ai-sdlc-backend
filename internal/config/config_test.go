package config_test

import (
	"os"
	"testing"

	"github.com/sanaul03/ai-sdlc-backend/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_NAME", "testdb")
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected default host 'localhost', got %q", cfg.Database.Host)
	}
	if cfg.Database.Port != "5432" {
		t.Errorf("expected default port '5432', got %q", cfg.Database.Port)
	}
	if cfg.Database.SSLMode != "disable" {
		t.Errorf("expected default sslmode 'disable', got %q", cfg.Database.SSLMode)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("expected default server port '8080', got %q", cfg.Server.Port)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_NAME", "mydb")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("SERVER_PORT", "9090")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("SERVER_PORT")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Database.Host != "db.example.com" {
		t.Errorf("expected host 'db.example.com', got %q", cfg.Database.Host)
	}
	if cfg.Database.Port != "5433" {
		t.Errorf("expected port '5433', got %q", cfg.Database.Port)
	}
	if cfg.Database.User != "admin" {
		t.Errorf("expected user 'admin', got %q", cfg.Database.User)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("expected server port '9090', got %q", cfg.Server.Port)
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	db := config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "pass",
		Name:     "db",
		SSLMode:  "disable",
	}
	want := "host=localhost port=5432 user=user password=pass dbname=db sslmode=disable"
	got := db.DSN()
	if got != want {
		t.Errorf("DSN mismatch:\n  want: %q\n   got: %q", want, got)
	}
}
