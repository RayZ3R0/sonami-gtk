package player2

import (
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

var playingFromCoverState = state.NewStateful[schwifty.Paintable](resources.MissingAlbum())
var playingFromTitleState = state.NewStateful[string]("Nothing")
var playingFromNavTargetState = state.NewStateful[string]("")

func init() {
	player.SourceChanged.OnLazy(func(s *player.Source) bool {
		if s != nil {
			playingFromTitleState.SetValue(s.Title)
			playingFromNavTargetState.SetValue(s.Route)

			if s.CoverURL != "" {
				texture, err := injector.MustInject[*imgutil.ImgUtil]().LoadCropped(s.CoverURL)
				if err != nil {
					slog.Error("failed to load source cover", "error", err)
					return signals.Continue
				}
				schwifty.OnMainThreadOncePure(func() {
					playingFromCoverState.SetValue(texture)
					texture.Unref()
				})
			} else {
				schwifty.OnMainThreadOncePure(func() {
					playingFromCoverState.SetValue(resources.MissingAlbum())
				})
			}
		}
		return signals.Continue
	})
}

func PlayingFrom() schwifty.Box {
	return HStack(
		VStack(
			Label(gettext.Get("Playing From")).HAlign(gtk.AlignStartValue).WithCSSClass("caption-heading").WithCSSClass("dimmed"),
			Label("").BindText(playingFromTitleState).HAlign(gtk.AlignStartValue).WithCSSClass("heading").Ellipsis(pango.EllipsizeEndValue),
		).VAlign(gtk.AlignCenterValue).MarginEnd(10),
		Image().
			BindPaintable(playingFromCoverState).
			Background("alpha(var(--view-fg-color), 0.1)").
			PixelSize(30).
			HAlign(gtk.AlignEndValue).HExpand(true).
			VAlign(gtk.AlignCenterValue).
			Overflow(gtk.OverflowHiddenValue).CornerRadius(5),
	).HAlign(gtk.AlignFillValue)
}
