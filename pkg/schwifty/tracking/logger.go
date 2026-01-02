package tracking

import (
	"log/slog"
	"time"
)

var logger = slog.With("library", "schwifty", "module", "tracking")

func LogAliveWidgets() {
	lastCount := 0
	for {
		widgets := Alive()
		if lastCount != len(widgets) {
			lastCount = len(widgets)
			slog.Info("alive widget count", "count", len(widgets))
		}
		time.Sleep(1 * time.Second)
	}
}
