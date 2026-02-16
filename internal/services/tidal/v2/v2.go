package v2

import "log/slog"

var logger = slog.With("service", "TIDAL").WithGroup("tidal").With("backend", "v2").WithGroup("v2")
