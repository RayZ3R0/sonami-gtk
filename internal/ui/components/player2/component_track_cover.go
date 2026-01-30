package player2

import (
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components"
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
			texture, err := injector.MustInject[*imgutil.ImgUtil]().LoadCropped(trackInfo.CoverURL)
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

func trackCover() schwifty.Picture {
	return components.SquirclePicture(
		Picture().
			FromPaintable(resources.MissingAlbum()).
			BindPaintable(coverState).HExpand(true).
			HAlign(gtk.AlignCenterValue),
	)
}
