package tracklist_header

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func componentControls(playFunc func(), shuffleFunc func(), popover *gtk.PopoverMenu) schwifty.Box {
	return HStack(
		Button().
			TooltipText(gettext.Get("Shuffle Album")).
			IconName("playlist-shuffle-symbolic").
			WithCSSClass("pill").
			BindSensitive(playerControlsAvailableState).
			ConnectClicked(func(b gtk.Button) {
				shuffleFunc()
			}),
		Button().
			TooltipText(gettext.Get("Play Album")).
			IconName("play-symbolic").
			WithCSSClass("pill").
			WithCSSClass("suggested-action").
			BindSensitive(playerControlsAvailableState).
			ConnectClicked(func(b gtk.Button) {
				playFunc()
			}),
		MenuButton().
			TooltipText(gettext.Get("More…")).
			Popover(popover).
			WithCSSClass("flat").
			WithCSSClass("circular").
			IconName("view-more-symbolic"),
	).Spacing(12).HAlign(gtk.AlignEndValue).HExpand(true)
}
