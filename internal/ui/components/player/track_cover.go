package player

import (
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var coverState = state.New[schwifty.Paintable](nil)

func init() {
	player.TrackChanged.On(func(trackInfo *player.Track) bool {
		if trackInfo != nil && trackInfo.CoverURL != "" {
			texture, err := injector.MustInject[*imgutil.ImgUtil]().Load(trackInfo.CoverURL)
			if err != nil {
				slog.Error("failed to load track cover", "error", err)
				return signals.Continue
			}
			schwifty.OnMainThreadOncePure(func() {
				coverState.SetValue(texture)
				texture.Unref()
			})
		} else {
			schwifty.OnMainThreadOncePure(func() {
				coverState.SetValue(resources.MissingAlbum())
			})
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
