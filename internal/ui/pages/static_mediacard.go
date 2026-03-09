package pages

import (
	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
)

// NewStaticMediaCardPage builds a scrollable WrapBox grid from a pre-loaded slice.
// If items is empty, the emptyState widget is returned instead.
func NewStaticMediaCardPage[T any](
	items []T,
	emptyState schwifty.BaseWidgetable,
	factory func(T) schwifty.BaseWidgetable,
) schwifty.BaseWidgetable {
	if len(items) == 0 {
		return emptyState
	}

	list := WrapBox().
		VMargin(20).
		HMargin(40).
		VAlign(gtk.AlignStartValue).
		Justify(adw.JustifyFillValue).
		JustifyLastLine(true)()

	for _, item := range items {
		list.Append(CenterBox().CenterWidget(factory(item)).ToGTK())
	}

	return ScrolledWindow().
		Child(components.MainContent(Widget(&list.Widget))).
		Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue)
}
