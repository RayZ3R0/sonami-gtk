package player

import (
	"log/slog"
	"math"

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

var (
	coverState = state.New[schwifty.Paintable](nil)
)

func init() {
	player.TrackChanged.On(func(trackInfo tonearm.Track) bool {
		if trackInfo != nil {
			go func() {
				coverUrl := trackInfo.Cover(math.MaxInt)
				if coverUrl == "" {
					slog.Error("Failed to load hi-res cover URL")
					return
				}

				texture, err := injector.MustInject[*imgutil.ImgUtil]().LoadCropped(coverUrl)
				if err != nil {
					slog.Error("failed to load hi-res track cover", "error", err)
					return
				}
				coverState.SetValue(texture)
			}()
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
			HAlign(gtk.AlignCenterValue).
			ConnectConstruct(func(p *gtk.Picture) {
				controller := gtk.NewGestureClick()
				controller.ConnectPressed(new(func(gtk.GestureClick, int, float64, float64) {
					components.GetMediaViewer().ShowFile(coverState.Value())
				}))
				p.AddController(&controller.EventController)
			}),
	)
}
