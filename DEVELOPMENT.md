# Technology Framework
- Use the Go programming language.
- Use PostgreSQL database.
- Don't use any ORM (Object Relational Mapping) framework. Use only the built-in `database/sql` Go package and [pgx connector](https://github.com/jackc/pgx)
- Use [golang-migrate with pgx](https://github.com/golang-migrate/migrate/tree/master/database/pgx/v5) for database migration

--- 

# Environment Configuration
The application uses environment variables for configuration. You can set these variables in a `.env` file in the root directory.

1. Copy the `.env.example` file to `.env`.
2. Update the `.env` file with your configuration values.

---

# REST API Convention
- Use `lower-kebab-case` for variable name
- Use the following path: `/api/vN/...` where `vN` is the version number, e.g. `v1` or `v2`
- Always define the valid HTTP method, i.e. for creating data only allow `POST`, and raise HTTP status 405 for method other than `POST`

---

# Database Design
- Always use audit trail fields:
  + `created_by` (text)
  + `updated_by` (text)
  + `created_at` (timestamp with timezone)
  + `updated_at` (timestamp with timezone)
  + `deleted` (bool)
  + `deleted_at` (timestamp with timezone)
- For datetime fields, always use `TIMESTAMP WITH TIME ZONE` data type

---

# Unit Testing
- Use Go table-based unit testing.
- Use [mockery](https://github.com/vektra/mockery) to generate mock objects for unit testing. The repository already contains `.mockery.yml` that can be used as a base mockery configuration.
- Use [assertion from testify](https://github.com/stretchr/testify?tab=readme-ov-file#assert-package) for asserting test result (wanted vs actual).
- Always create unit tests whenever possible.
- Ensure the unit test passes before submitting a pull request.

--- 

# Folder Structure
## Overview
- `cmd`: application entry point
- `db`  
  - `migration`: SQL files for database migration/initialization
- `internal`: working folder for feature code, should not be used/imported by other Go apps
    - `entity`: database table representation
    - `handler`: REST API handler
    - `model`: REST API request/response representation
    - `server`: 
      + `server.go`: contains Go multiplexer and API handler registration
      + `middleware.go`: HTTP middleware/interceptor
    `store`: database access layer for CRUD operation

For more information, see the `Detailed Description` section below.


## Detailed Description
### `cmd`
Application entry point. This folder should contain only one file: `app.go`, which has the `main()` function.

### `db/migration`
Contains all SQL files for database migration, includes:
- database object definition, such as create table, create index, etc.
- data initialization
  
The general rule is that one object definition should be in one SQL file.
If there is an additional definition (e.g., additional constraints for table `mytable`), you should put the additional definition in the same SQL file as the main object. In this sample scenario, the file `001_mytable.sql` can contain:
- DDL for `CREATE TABLE IF NOT EXISTS mytable ...`
- additional constraints for `mytable`, e.g., `ALTER TABLE mytable ...`

The naming convention for SQL files is lower snake case with the format below.
- `[3 digit sequence]_descriptive_sql_goal.up.sql` (for example: `001_create_table_vehicles.up.sql`) for object creation
- `[3 digit sequence]_descriptive_sql_goal.down.sql` (for example: `001_create_table_vehicles.down.sql`) for object destruction
- The naming conventions are aligned with the standard for [Go migrate](https://github.com/golang-migrate/migrate?tab=readme-ov-file#migration-files)

The object creation/destruction must be safe:
- Use `IF EXISTS` when creating an object
- Use `IF NOT EXISTS` when destroying an object

### `internal/entity`
Go structs that represent database tables or DTO (data transfer object).

### `internal/handler`
REST API handler. Any request/response body must be taken from the `internal/model`.

### `internal/model`
Go structs that represent API request and response.

### `internal/server`
All API handlers on `internal/handler` must be registered in the file `internal/server/server.go`.  
Any required middleware (such as authorization) must be placed on the file `internal/server/middleware.go`.

### `internal/store`
Go files that represent the data access layer to handle database operations (CRUD).

---

# Code Example
The default code contains files named `exampleXxx.go`.
Refer to those files to understand the basic structure of creating a REST API with a PostgreSQL database.
In a nutshell, follow the steps below.

1. Create SQL migration file under `db/migration`. Always create `.up.sql` and `down.sql` files. Example: `db/migration/000_create_table_my_examples.up.sql` and `db/migration/000_create_table_my_examples.down.sql`.
2. Create the struct to represent the table. Example: see `type Example struct` on `internal/entity/example.go`.
3. Create an interface for the data store layer. Example: see `type ExampleStoreProvider interface` on `internal/store/example.go`.
4. Create the data store implementation. Example: see `type ExampleDatabaseStore struct` on `internal/store/exampledb.go`.
5. Create the REST API request and response. Example: Go structs on `internal/model/example.go`.
6. Create the REST API handlers. Example: Go structs on `internal/handler/example.go`.
7. Register the REST API handlers on the multiplexer. If needed, wrap the handler with middleware. Example: see `func Mux()` on `internal/server/server.go`.
