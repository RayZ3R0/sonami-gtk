package tracklist

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func ControlsColumn(track tonearm.Track, grid *gtk.Grid, position int, column int) int {
	grid.Attach(
		HStack(
			Button().
				TooltipText(gettext.Get("Add to Collection")).
				IconName("heart-outline-thick-symbolic").
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue).
				WithCSSClass("flat").Sensitive(false),
			Button().
				TooltipText(gettext.Get("Add to Queue")).
				IconName("plus-symbolic").
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue).
				ActionName("win.player.queue-track").
				ActionTargetValue(glib.NewVariantString(track.ID())).
				WithCSSClass("flat"),
		).
			Margin(10).
			HAlign(gtk.AlignEndValue).
			ToGTK(),
		column,
		0,
		1,
		1,
	)
	return 1
}
