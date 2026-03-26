package player

import (
	"log/slog"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"codeberg.org/puregotk/puregotk/v4/pango"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/resources"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/imgutil"
	"github.com/infinytum/injector"
)

var playingFromCoverState = state.NewStateful[schwifty.Paintable](resources.MissingAlbum())
var playingFromTitleState = state.NewStateful[string]("Nothing")
var playingFromCanNavigateState = state.NewStateful[bool](false)
var playingFromNavTargetState = state.NewStateful[string]("")

func init() {
	player.SourceChanged.OnLazy(func(s sonami.PlaybackSource) bool {
		if s != nil {
			playingFromTitleState.SetValue(s.Title())
			playingFromNavTargetState.SetValue(s.Route())
			playingFromCanNavigateState.SetValue(s.Route() != "")

			coverUrl := s.Cover(80)
			if coverUrl == "" {
				schwifty.OnMainThreadOncePure(func() {
					playingFromCoverState.SetValue(resources.MissingAlbum())
				})
				return signals.Continue
			}

			texture, err := injector.MustInject[*imgutil.ImgUtil]().LoadCropped(coverUrl)
			if err != nil {
				slog.Error("failed to load source cover", "error", err)
				return signals.Continue
			}
			schwifty.OnMainThreadOncePure(func() {
				playingFromCoverState.SetValue(texture)
			})
		}
		return signals.Continue
	})
}

func PlayingFrom() schwifty.Button {
	return Button().
		BindSensitive(playingFromCanNavigateState).
		ConnectClicked(func(b gtk.Button) {
			router.Navigate(playingFromNavTargetState.Value())
		}).
		Child(
			HStack(
				VStack(
					Label(gettext.Get("Playing From")).HAlign(gtk.AlignStartValue).WithCSSClass("caption-heading").WithCSSClass("dimmed"),
					Label("").BindText(playingFromTitleState).HAlign(gtk.AlignStartValue).WithCSSClass("heading").Ellipsis(pango.EllipsizeEndValue),
				).VAlign(gtk.AlignCenterValue).MarginEnd(10),
				Image().
					BindPaintable(playingFromCoverState).
					Background("alpha(var(--view-fg-color), 0.1)").
					PixelSize(33).
					HAlign(gtk.AlignEndValue).HExpand(true).
					VAlign(gtk.AlignCenterValue).
					Overflow(gtk.OverflowHiddenValue).CornerRadius(5),
			).HAlign(gtk.AlignFillValue),
		).WithCSSClass("flat").CSS("button { margin-top: -10px; margin-bottom: -10px; margin-left: -10px; margin-right: -10px; padding-top: 10px; padding-bottom: 10px; }")
}
