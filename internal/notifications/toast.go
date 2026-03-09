package notifications

import "github.com/RayZ3R0/sonami-gtk/internal/signals"

var OnToast = signals.NewStatelessSignal[string]()
