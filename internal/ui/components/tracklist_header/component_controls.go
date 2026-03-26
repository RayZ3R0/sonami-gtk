package tracklist_header

import (
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
)

func componentControls(playFunc func(), shuffleFunc func()) schwifty.Box {
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
	).Spacing(12).HAlign(gtk.AlignEndValue).HExpand(true)
}
