package migrations

import "embed"

// Files contains SQL migration files embedded into the API binary.
//
//go:embed *.sql
var Files embed.FS
