package internal

import (
	"log/slog"
)

var (
	logger             = slog.With("library", "schwifty", "module", "resolver")
	shouldLogLifecycle = false
)
