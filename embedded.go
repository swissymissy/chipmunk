package chipmunk

import "embed"

// FrontendFS holds the static HTML/CSS/JS/images served at /.
// Embedded so the binary is fully self-contained.
//
//go:embed all:cmd/frontend
var FrontendFS embed.FS

// SchemaFS holds the goose migration files run at startup.
//
//go:embed sql/schema/*.sql
var SchemaFS embed.FS
