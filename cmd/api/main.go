package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sanaul03/ai-sdlc-backend/internal/config"
	"github.com/sanaul03/ai-sdlc-backend/internal/database"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet/handler"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet/postgres"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := runMigrations(cfg); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	db, err := database.New(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	// Repositories
	carGroupRepo := postgres.NewCarGroupRepository(db)
	vehicleRepo := postgres.NewVehicleRepository(db)

	// Services
	carGroupSvc := fleet.NewCarGroupService(carGroupRepo)
	vehicleSvc := fleet.NewVehicleService(vehicleRepo)

	// HTTP handlers
	mux := http.NewServeMux()
	handler.NewCarGroupHandler(carGroupSvc).RegisterRoutes(mux)
	handler.NewVehicleHandler(vehicleSvc).RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-stop
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}

func runMigrations(cfg *config.Config) error {
	m, err := migrate.New("file://migrations", cfg.Database.MigrateDSN())
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}
	log.Println("migrations applied successfully")
	return nil
}
