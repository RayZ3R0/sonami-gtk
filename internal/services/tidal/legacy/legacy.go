package legacy

import "log/slog"

var logger = slog.With("service", "TIDAL", "version", "legacy")
