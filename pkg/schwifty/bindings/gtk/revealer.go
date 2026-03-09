package gtk

import "codeberg.org/puregotk/puregotk/v4/gtk"

//go:generate go run github.com/RayZ3R0/sonami-gtk/pkg/schwifty/gen Revealer *gtk.Revealer gtk

func (r Revealer) TransitionType(t gtk.RevealerTransitionType) Revealer {
	return func() *gtk.Revealer {
		w := r()
		w.SetTransitionType(t)
		return w
	}
}
