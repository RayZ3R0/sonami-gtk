package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func ScrolledWindow() schwifty.ScrolledWindow {
	return managed("ScrolledWindow", func() *gtk.ScrolledWindow {
		scrolledWindow := gtk.NewScrolledWindow()
		scrolledWindow.ConnectEdgeReached(&callback.ScrolledWindowEdgeReachedCallback)
		return scrolledWindow
	})
}
