package notifications

import "codeberg.org/dergs/tidalwave/internal/signals"

var OnToast = signals.NewStatelessSignal[string]()
