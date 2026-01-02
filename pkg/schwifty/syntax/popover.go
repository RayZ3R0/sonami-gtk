package syntax

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Popover(child any) schwifty.Popover {
	return managed("Popover", func() *gtk.Popover {
		popover := gtk.NewPopover()
		popover.SetChild(schwifty.ResolveWidget(child))
		return popover
	})
}
