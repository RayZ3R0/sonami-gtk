package secrets

import (
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
)

// SignedInChanged and SignedInState are always true in account-free mode.
// The app behaves as if the user is always authenticated.
var SignedInChanged *signals.StatefulSignal[bool]
var SignedInState *state.State[bool]

func init() {
	SignedInState = state.NewStateful(true)
	SignedInChanged = signals.NewStatefulSignal(true)
	SignedInChanged.On(func(b bool) bool {
		SignedInState.SetValue(b)
		return signals.Continue
	})
}

func triggerSignedInChanged() {
	// No-op in account-free mode
}
