// Package migrations embeds the app SQL migrations. Each feature adds its own
// file from 00002 on.
package migrations

import "embed"

// FS holds the migrations at its root, as goose expects. Read-only.
//
//go:embed *.sql
var FS embed.FS
