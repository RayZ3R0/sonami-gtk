package secrets

import "codeberg.org/dergs/tonearm/internal/signals"

var SignedInChanged *signals.StatefulSignal[bool]

func init() {
	SignedInChanged = signals.NewStatefulSignal(HasRefreshToken())
}

func triggerSignedInChanged() {
	SignedInChanged.Notify(func(oldValue bool) bool {
		return HasRefreshToken()
	})
}
