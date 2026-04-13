package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sanaul03/ai-sdlc-backend/internal/cargroup"
	"github.com/sanaul03/ai-sdlc-backend/internal/platform/database"
	"github.com/sanaul03/ai-sdlc-backend/internal/platform/middleware"
	"github.com/sanaul03/ai-sdlc-backend/internal/vehicle"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("startup error: %v", err)
	}
}

func run() error {
	ctx := context.Background()

	// --- Database ---
	dbCfg := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "ai_sdlc"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	pool, err := database.Connect(ctx, dbCfg)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	// --- Migrations ---
	migrationsPath := getEnv("MIGRATIONS_PATH", "file://migrations")
	m, err := migrate.New(migrationsPath, "pgx5://"+dbCfg.DSN())
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	// --- JWT public key ---
	publicKey, err := loadRSAPublicKey(getEnv("JWT_PUBLIC_KEY_PATH", ""))
	if err != nil {
		return fmt.Errorf("load JWT public key: %w", err)
	}

	// --- Dependencies ---
	cgRepo := cargroup.NewRepository(pool)
	cgSvc := cargroup.NewService(cgRepo)
	cgHandler := cargroup.NewHandler(cgSvc)

	vRepo := vehicle.NewRepository(pool)
	vSvc := vehicle.NewService(vRepo)
	vHandler := vehicle.NewHandler(vSvc)

	// --- Router ---
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(publicKey))

		// Write endpoints – FLEET_MANAGER only
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("FLEET_MANAGER"))
			cgHandler.RegisterWriteRoutes(r)
			vHandler.RegisterWriteRoutes(r)
		})

		// Read-only endpoints – all authenticated roles
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("FLEET_MANAGER", "BRANCH_STAFF", "OPERATIONS_MANAGER"))
			cgHandler.RegisterReadRoutes(r)
			vHandler.RegisterReadRoutes(r)
		})
	})

	// --- Server ---
	addr := ":" + getEnv("PORT", "8080")
	log.Printf("server listening on %s", addr)
	return http.ListenAndServe(addr, r)
}

func loadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	if path == "" {
		return nil, fmt.Errorf("JWT_PUBLIC_KEY_PATH environment variable not set")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read public key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from public key file")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	rsaKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA")
	}

	return rsaKey, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
