package openapi

import "log/slog"

var logger = slog.With("service", "TIDAL").WithGroup("tidal").With("backend", "openapi").WithGroup("openapi")
