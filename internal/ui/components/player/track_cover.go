package player

import (
	"log/slog"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdkpixbuf"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var coverState = state.New[*gdkpixbuf.Pixbuf](nil)

func init() {
	player.OnTrackChanged.On(func(trackInfo player.TrackInformation) bool {
		if trackInfo.CoverURL != "" {
			pixbuf, err := injector.MustInject[*imgutil.ImgUtil]().LoadPixbuf(trackInfo.CoverURL)
			if err != nil {
				slog.Error("failed to load track cover", "error", err)
				return signals.Continue
			}
			coverState.SetValue(pixbuf)
			pixbuf.Unref()
		} else {
			pixbuf, err := gdkpixbuf.NewPixbufFromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg")
			if err != nil {
				slog.Error("failed to load track cover from resource", "error", err)
				return signals.Continue
			}
			coverState.SetValue(pixbuf)
			pixbuf.Unref()
		}
		return signals.Continue
	})
}

func trackCover() schwifty.AspectFrame {
	return AspectFrame(
		Image().
			PixelSize(319).
			Overflow(gtk.OverflowHiddenValue).
			FromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg").
			BindPixbuf(coverState),
	).
		CornerRadius(10).
		Overflow(gtk.OverflowHiddenValue).
		HAlign(gtk.AlignCenterValue).
		Background("alpha(var(--view-fg-color), 0.1)")
}
