package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sanaul03/ai-sdlc-backend/internal/config"
	"github.com/sanaul03/ai-sdlc-backend/internal/server"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	if err := runMigrations(cfg.DB); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, server.New(pool)); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func runMigrations(db config.DBConfig) error {
	dsn := fmt.Sprintf(
		"pgx5://%s:%s@%s:%s/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode,
	)
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
