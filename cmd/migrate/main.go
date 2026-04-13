// Command migrate is a standalone CLI tool for managing database schema migrations.
//
// Usage:
//
//	migrate <command> [args]
//
// Commands:
//
//	up              Apply all pending migrations
//	down [N]        Roll back N migrations (default: 1)
//	version         Print the current migration version
//	force <version> Force-set the migration version without running any SQL
//
// The tool reads database connection settings from the same environment
// variables used by the API server (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD,
// DB_NAME, DB_SSLMODE).  The migrations directory is resolved relative to the
// working directory, so run the binary from the repository root.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sanaul03/ai-sdlc-backend/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	m, err := migrate.New("file://migrations", cfg.Database.MigrateDSN())
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}
	defer m.Close()

	command := os.Args[1]

	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate up: %v", err)
		}
		log.Println("migrations applied successfully")

	case "down":
		steps := 1
		if len(os.Args) >= 3 {
			n, err := strconv.Atoi(os.Args[2])
			if err != nil || n < 1 {
				log.Fatalf("down: invalid step count %q (must be a positive integer)", os.Args[2])
			}
			steps = n
		}
		if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate down: %v", err)
		}
		log.Printf("rolled back %d migration(s)", steps)

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				log.Println("version: no migrations applied yet")
				return
			}
			log.Fatalf("migrate version: %v", err)
		}
		if dirty {
			log.Printf("version: %d (dirty — last migration did not complete cleanly)", version)
		} else {
			log.Printf("version: %d", version)
		}

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("force: version argument required")
		}
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("force: invalid version %q (must be an integer)", os.Args[2])
		}
		if err := m.Force(v); err != nil {
			log.Fatalf("migrate force: %v", err)
		}
		log.Printf("forced migration version to %d", v)

	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", command)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: migrate <command> [args]

Commands:
  up              Apply all pending migrations
  down [N]        Roll back N migrations (default: 1)
  version         Print the current migration version
  force <version> Force-set migration version without running SQL

Environment variables (all required unless a default is shown):
  DB_HOST      PostgreSQL host     (default: localhost)
  DB_PORT      PostgreSQL port     (default: 5432)
  DB_USER      Database user
  DB_PASSWORD  Database password
  DB_NAME      Database name
  DB_SSLMODE   SSL mode            (default: disable)

Run from the repository root so that the migrations/ directory is found.
`)
}
