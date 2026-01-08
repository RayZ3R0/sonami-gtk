package player

import (
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

var titleState = state.NewStateful("No Track")

func init() {
	player.TrackChanged.On(func(trackInfo *player.Track) bool {
		schwifty.OnMainThreadOncePure(func() {
			if trackInfo == nil {
				titleState.SetValue("No Track")
			} else {
				titleState.SetValue(trackInfo.Title)
			}
		})
		return signals.Continue
	})
}

func trackTitle() schwifty.Label {
	return Label("No Track").
		BindText(titleState).
		FontSize(24).
		FontWeight(800).
		LineHeight(1.2).
		HMargin(32).
		HAlign(gtk.AlignCenterValue).
		Ellipsis(pango.EllipsizeEndValue).
		MarginTop(35)
}
