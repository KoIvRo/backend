package migrations

import "embed"

// Files contains all SQL migrations compiled into the migration binary.
//
//go:embed *.sql
var Files embed.FS
