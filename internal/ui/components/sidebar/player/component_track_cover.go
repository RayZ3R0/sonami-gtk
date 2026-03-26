package player

import (
	"log/slog"
	"math"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/resources"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/imgutil"
	"github.com/infinytum/injector"
)

var (
	coverState = state.New[schwifty.Paintable](nil)
)

func init() {
	player.TrackChanged.On(func(trackInfo sonami.Track) bool {
		if trackInfo != nil {
			coverUrl := trackInfo.Cover(math.MaxInt)
			if coverUrl == "" {
				slog.Error("Failed to load hi-res cover URL")
				return signals.Continue
			}

			texture, err := injector.MustInject[*imgutil.ImgUtil]().LoadCropped(coverUrl)
			if err != nil {
				slog.Error("failed to load hi-res track cover", "error", err)
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
			HAlign(gtk.AlignCenterValue).
			ConnectConstruct(func(p *gtk.Picture) {
				controller := gtk.NewGestureClick()
				controller.ConnectPressed(new(func(gtk.GestureClick, int32, float64, float64) {
					components.GetMediaViewer().ShowFile(coverState.Value())
				}))
				p.AddController(&controller.EventController)
			}),
	)
}
