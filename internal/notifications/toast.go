package notifications

import "codeberg.org/dergs/tidalwave/internal/signals"

var OnToast = toastOnNotifySignal{
	signals.NewSignal[func(path string) bool](),
}

type toastOnNotifySignal struct {
	signals.Signal[func(title string) bool]
}

func (r *toastOnNotifySignal) Notify(title string) {
	r.Signal.Notify(title)
}
