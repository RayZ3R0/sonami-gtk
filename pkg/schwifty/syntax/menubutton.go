package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func MenuButton() schwifty.MenuButton {
	return managed("MenuButton", func() *gtk.MenuButton {
		return gtk.NewMenuButton()
	})
}
