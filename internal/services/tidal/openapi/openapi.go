package openapi

import "log/slog"

var logger = slog.With("service", "TIDAL", "version", "openapi")
