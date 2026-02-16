package v1

import "log/slog"

var logger = slog.With("service", "TIDAL").WithGroup("tidal").With("backend", "v1").WithGroup("v1")
