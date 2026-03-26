package player

import (
	"strings"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"codeberg.org/puregotk/puregotk/v4/pango"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

var titleState = state.NewStateful("")
var artistState = state.NewStateful("")

func init() {
	player.TrackChanged.On(func(trackInfo sonami.Track) bool {
		schwifty.OnMainThreadOncePure(func() {
			if trackInfo == nil {
				titleState.SetValue(gettext.Get("No Track"))
				artistState.SetValue(gettext.Get("No Artist"))
			} else {
				titleState.SetValue(sonami.FormatTitle(trackInfo))
				artistState.SetValue(strings.Join(trackInfo.Artists().Names(), ", "))
			}
		})
		return signals.Continue
	})
}

func trackInfo() schwifty.Box {
	return VStack(
		Label("").
			BindText(titleState).
			HAlign(gtk.AlignCenterValue).
			Ellipsis(pango.EllipsizeEndValue).
			WithCSSClass("title-4"),
		Label("").
			BindText(artistState).
			HAlign(gtk.AlignCenterValue).
			Ellipsis(pango.EllipsizeEndValue).
			WithCSSClass("heading").WithCSSClass("dimmed"),
	).Spacing(5)
}
