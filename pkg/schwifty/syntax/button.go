package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Button() schwifty.Button {
	return managed("Button", func() *gtk.Button {
		btn := gtk.NewButton()
		btn.ConnectClicked(&callback.ButtonClickedCallback)
		return btn
	})
}
