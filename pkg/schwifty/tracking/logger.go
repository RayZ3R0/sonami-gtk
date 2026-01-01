package tracking

import (
	"log/slog"
	"time"
)

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
