// Package migrations embeds the SQL migration files.
package migrations

import "embed"

// FS holds all migration SQL files.
//
//go:embed *.sql
var FS embed.FS
