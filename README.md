# AI SDLC Backend

REST API backend for the AI SDLC fleet management system, implementing the Fleet Structure module.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Environment Configuration](#environment-configuration)
3. [Database Setup](#database-setup)
4. [Database Migrations](#database-migrations)
5. [Running the Application](#running-the-application)
6. [Running Tests](#running-tests)
7. [Project Structure](#project-structure)
8. [API Reference](#api-reference)

---

## Prerequisites

| Tool | Minimum Version | Notes |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.25 | Language runtime |
| [PostgreSQL](https://www.postgresql.org/download/) | 14 | Primary database |

---

## Environment Configuration

The application is configured entirely via environment variables. Copy the table below into a `.env` file (the file is already in `.gitignore`) and fill in the values before starting.

```dotenv
# ── Database ──────────────────────────────────────────────────────────────────
DB_HOST=localhost        # PostgreSQL host          (default: localhost)
DB_PORT=5432             # PostgreSQL port          (default: 5432)
DB_USER=your_db_user     # Database user            (required)
DB_PASSWORD=your_db_pass # Database password        (required)
DB_NAME=ai_sdlc          # Database name            (required)
DB_SSLMODE=disable       # SSL mode                 (default: disable)
                         # Use 'require' for production

# ── HTTP server ───────────────────────────────────────────────────────────────
SERVER_PORT=8080         # Port the API listens on  (default: 8080)
```

Export the variables before running the binary, for example:

```bash
export $(grep -v '^#' .env | xargs)
```

### Variable Reference

| Variable | Required | Default | Description |
|---|---|---|---|
| `DB_HOST` | No | `localhost` | PostgreSQL server hostname or IP |
| `DB_PORT` | No | `5432` | PostgreSQL server port |
| `DB_USER` | **Yes** | — | PostgreSQL username |
| `DB_PASSWORD` | **Yes** | — | PostgreSQL password |
| `DB_NAME` | **Yes** | — | PostgreSQL database name |
| `DB_SSLMODE` | No | `disable` | pgx SSL mode (`disable`, `require`, `verify-full`, …) |
| `SERVER_PORT` | No | `8080` | TCP port the HTTP server binds to |

---

## Database Setup

### 1. Create the database

Connect to PostgreSQL as a superuser and create the application database and user:

```sql
CREATE DATABASE ai_sdlc;
CREATE USER ai_sdlc_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE ai_sdlc TO ai_sdlc_user;
```

---

## Database Migrations

Schema migrations are managed with [golang-migrate](https://github.com/golang-migrate/migrate) using raw SQL files in the `migrations/` directory — one table per file, named `NNNNNN_<description>.<up|down>.sql`.

| File | Description |
|---|---|
| `000001_create_car_groups_table.up.sql` | Creates the `car_groups` table |
| `000002_create_vehicles_table.up.sql` | Creates the `vehicles` table with FK, indices, and CHECK constraints |

### Automatic migrations

The API server runs `migrate up` automatically before accepting requests, so **no manual migration step is required** for normal operation.

### Standalone migration CLI

A dedicated CLI binary at `cmd/migrate` lets you run, roll back, or inspect migrations independently of the API server. Run all commands from the **repository root** so the `migrations/` directory is found.

```bash
# Build the migration binary
go build -o bin/migrate ./cmd/migrate
```

| Command | Description |
|---|---|
| `./bin/migrate up` | Apply all pending migrations |
| `./bin/migrate down [N]` | Roll back N migrations (default: 1) |
| `./bin/migrate version` | Print the current schema version |
| `./bin/migrate force <version>` | Force-set version without running SQL (use after manual fixes) |

```bash
# Export environment variables first (see Environment Configuration above)
export DB_USER=ai_sdlc_user
export DB_PASSWORD=your_secure_password
export DB_NAME=ai_sdlc

# Apply all pending migrations
./bin/migrate up

# Roll back the last migration
./bin/migrate down

# Roll back the last 2 migrations
./bin/migrate down 2

# Check the current schema version
./bin/migrate version

# Force-set version to 1 (after a manual schema fix)
./bin/migrate force 1
```

You can also run without building first:

```bash
go run ./cmd/migrate up
go run ./cmd/migrate down
go run ./cmd/migrate version
```

---

## Running the Application

### From source

```bash
# 1. Clone the repository
git clone https://github.com/sanaul03/ai-sdlc-backend.git
cd ai-sdlc-backend

# 2. Install Go module dependencies
go mod download

# 3. Set required environment variables
export DB_USER=ai_sdlc_user
export DB_PASSWORD=your_secure_password
export DB_NAME=ai_sdlc
# Optional overrides (shown with their defaults):
# export DB_HOST=localhost
# export DB_PORT=5432
# export DB_SSLMODE=disable
# export SERVER_PORT=8080

# 4. (Optional) Run migrations manually before starting the server
go run ./cmd/migrate up

# 5. Start the API server
go run ./cmd/api
```

The server automatically applies any pending migrations before accepting requests, so step 4 is only needed if you want to manage migrations separately.

### Build binaries

```bash
# Build both the API server and the migration tool
go build -o bin/api     ./cmd/api
go build -o bin/migrate ./cmd/migrate

# Run migrations then start the server
./bin/migrate up
./bin/api
```

### Expected startup output

```
2026/04/13 12:00:00 migrations applied successfully
2026/04/13 12:00:00 server listening on :8080
```

The API is now reachable at `http://localhost:8080/api/v1`.

---

## Running Tests

```bash
# Run all unit tests
go test ./...

# Run with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./internal/fleet/...
go test -v ./internal/fleet/handler/...
```

All tests are pure unit tests (no database connection required).

---

## Project Structure

```
.
├── cmd/
│   ├── api/
│   │   └── main.go                  # Entry point: config, migrations, HTTP server
│   └── migrate/
│       └── main.go                  # Standalone migration CLI (up/down/version/force)
├── internal/
│   ├── config/
│   │   └── config.go                # Environment variable loading
│   ├── database/
│   │   └── database.go              # pgxpool connection factory
│   └── fleet/                       # Fleet Structure domain
│       ├── car_group.go             # CarGroup model and input types
│       ├── vehicle.go               # Vehicle model, enums, and input types
│       ├── errors.go                # Domain error sentinels
│       ├── repository.go            # Repository interfaces
│       ├── car_group_service.go     # CarGroup business logic
│       ├── vehicle_service.go       # Vehicle business logic and validation
│       ├── postgres/
│       │   ├── car_group_repository.go  # PostgreSQL CarGroup repository
│       │   └── vehicle_repository.go   # PostgreSQL Vehicle repository
│       └── handler/
│           ├── car_group_handler.go # HTTP handlers for /api/v1/car-groups
│           └── vehicle_handler.go   # HTTP handlers for /api/v1/vehicles
└── migrations/
    ├── 000001_create_car_groups_table.up.sql
    ├── 000001_create_car_groups_table.down.sql
    ├── 000002_create_vehicles_table.up.sql
    └── 000002_create_vehicles_table.down.sql
```

---

## API Reference

Base URL: `http://localhost:8080/api/v1`

The caller identity is read from the `X-User-ID` request header (falls back to `"system"` if absent).

### Car Groups

| Method | Path | Description |
|---|---|---|
| `POST` | `/car-groups` | Create a car group |
| `GET` | `/car-groups` | List car groups |
| `GET` | `/car-groups/{id}` | Get a car group by ID |
| `PUT` | `/car-groups/{id}` | Update a car group |
| `DELETE` | `/car-groups/{id}` | Soft-delete a car group |

#### Create car group — `POST /car-groups`

```json
{
  "name": "Economy Sedan",
  "description": "Small, fuel-efficient sedans",
  "size_category": "compact"
}
```

**Query parameters for `GET /car-groups`:**

| Parameter | Type | Description |
|---|---|---|
| `q` | string | Case-insensitive name search |
| `deleted` | bool | Include soft-deleted records (default `false`) |

---

### Vehicles

| Method | Path | Description |
|---|---|---|
| `POST` | `/vehicles` | Create a vehicle |
| `GET` | `/vehicles` | List vehicles (paginated) |
| `GET` | `/vehicles/{id}` | Get a vehicle by ID |
| `PUT` | `/vehicles/{id}` | Update a vehicle |
| `DELETE` | `/vehicles/{id}` | Soft-delete a vehicle |
| `PATCH` | `/vehicles/{id}/designation` | Update designation only |

#### Create vehicle — `POST /vehicles`

```json
{
  "car_group_id": "550e8400-e29b-41d4-a716-446655440000",
  "branch_id":    "550e8400-e29b-41d4-a716-446655440001",
  "vin":          "1HGBH41JXMN109186",
  "licence_plate": "ABC-1234",
  "brand":        "Toyota",
  "model":        "Corolla",
  "year":         2022,
  "fuel_type":    "petrol",
  "transmission_type": "automatic",
  "current_mileage":   0,
  "status":       "unavailable",
  "designation":  "rental_only",
  "acquisition_date": "2022-01-15",
  "ownership_type": "owned"
}
```

#### Update designation — `PATCH /vehicles/{id}/designation`

```json
{
  "designation": "sales_only"
}
```

**Query parameters for `GET /vehicles`:**

| Parameter | Type | Description |
|---|---|---|
| `car_group_id` | UUID | Filter by car group |
| `branch_id` | UUID | Filter by branch |
| `status` | string | Filter by status (`available`, `on_rent`, `needs_cleaning`, `needs_inspection`, `under_maintenance`, `unavailable`, `decommissioned`) |
| `designation` | string | Filter by designation (`rental_only`, `sales_only`, `shared`) |
| `fuel_type` | string | Filter by fuel type (`petrol`, `diesel`, `electric`, `hybrid`) |
| `transmission_type` | string | Filter by transmission (`manual`, `automatic`) |
| `expiry_warning` | bool | `true` to return only vehicles with insurance/registration expiry ≤ 30 days |
| `page` | int | Page number (default `1`) |
| `page_size` | int | Results per page (default `20`) |

**List response envelope:**

```json
{
  "data": [ { "...": "..." } ],
  "total": 42
}
```

Each vehicle object includes an `expiry_warning: true` flag when insurance or registration expires within 30 days.

### Common HTTP status codes

| Code | Meaning |
|---|---|
| `200 OK` | Successful read or update |
| `201 Created` | Resource created |
| `204 No Content` | Successful delete |
| `400 Bad Request` | Invalid request body or path parameter |
| `404 Not Found` | Resource does not exist |
| `409 Conflict` | Operation blocked by existing data (e.g., deleting a car group with active vehicles) |
| `500 Internal Server Error` | Unexpected server error |
