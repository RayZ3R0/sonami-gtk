package player

import (
	"log/slog"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/resources"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var coverState = state.New[gdk.Paintable](nil)

func init() {
	player.OnTrackChanged.On(func(trackInfo player.TrackInformation) bool {
		if trackInfo.CoverURL != "" {
			texture, err := injector.MustInject[*imgutil.ImgUtil]().Load(trackInfo.CoverURL)
			if err != nil {
				slog.Error("failed to load track cover", "error", err)
				return signals.Continue
			}
			coverState.SetValue(texture)
			texture.Unref()
		} else {
			coverState.SetValue(resources.MissingAlbum())
		}
		return signals.Continue
	})
}

func trackCover() schwifty.AspectFrame {
	return AspectFrame(
		Image().
			PixelSize(380).
			Overflow(gtk.OverflowHiddenValue).
			FromPaintable(resources.MissingAlbum()).
			BindPaintable(coverState),
	).
		CornerRadius(10).
		Overflow(gtk.OverflowHiddenValue).
		HAlign(gtk.AlignCenterValue).
		Background("alpha(var(--view-fg-color), 0.1)")
}
