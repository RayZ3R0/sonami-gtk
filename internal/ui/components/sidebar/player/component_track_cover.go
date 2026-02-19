package player

import (
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var coverState = state.New[schwifty.Paintable](nil)

func init() {
	player.TrackChanged.On(func(trackInfo tonearm.Track) bool {
		if trackInfo != nil {
			coverUrl := trackInfo.Cover(320)
			if coverUrl == "" {
				slog.Error("Failed to load cover URL")
				return signals.Continue
			}

			texture, err := injector.MustInject[*imgutil.ImgUtil]().LoadCropped(coverUrl)
			if err != nil {
				slog.Error("failed to load track cover", "error", err)
				return signals.Continue
			}
			coverState.SetValue(texture)
		} else {
			coverState.SetValue(resources.MissingAlbum())
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
