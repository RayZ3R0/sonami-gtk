package sidebar

import (
	"log/slog"

	"codeberg.org/puregotk/puregotk/v4/gtk"
	"codeberg.org/puregotk/puregotk/v4/pango"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/resources"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	adwbindings "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/bindings/adw"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/utils/imgutil"
	"github.com/infinytum/injector"
)

func MiniPlayer() adwbindings.Bin {
	var miniPlayerLoadingIconSub *signals.Subscription
	var miniPlayerCanControl = state.NewStateful(false)

	var (
		coverState    = state.NewStateful[schwifty.Paintable](resources.MissingAlbum())
		trackTitle    = state.NewStateful[string]("")
		trackAlbum    = state.NewStateful[string]("")
		playPauseIcon = state.NewStateful("play-symbolic")
	)

	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		miniPlayerCanControl.SetValue(!ps.Loading)
		return signals.Continue
	})

	player.TrackChanged.On(func(trackInfo sonami.Track) bool {
		if trackInfo != nil {
			coverUrl := trackInfo.Cover(80)
			if coverUrl == "" {
				slog.Error("Failed to load cover URL")
				return signals.Continue
			}

			if texture, err := injector.MustInject[*imgutil.ImgUtil]().Load(coverUrl); err == nil {
				schwifty.OnMainThreadOncePure(func() {
					coverState.SetValue(texture)
				})
			}

			schwifty.OnMainThreadOncePure(func() {
				trackTitle.SetValue(sonami.FormatTitle(trackInfo))
				trackAlbum.SetValue(trackInfo.Album().Title())
			})
		}

		return signals.Continue
	})

	player.PlaybackStateChanged.On(func(state *player.PlaybackState) bool {
		schwifty.OnMainThreadOncePure(func() {
			switch state.Status {
			case player.PlaybackStatusPaused, player.PlaybackStatusStopped:
				playPauseIcon.SetValue("play-symbolic")
			case player.PlaybackStatusPlaying:
				playPauseIcon.SetValue("pause-symbolic")
			}
		})
		return signals.Continue
	})

	return Bin().Child(
		HStack(
			AspectFrame(
				Image().
					PixelSize(54).
					BindPaintable(coverState),
			).
				Overflow(gtk.OverflowHiddenValue).
				Background("alpha(var(--view-fg-color), 0.1)").
				CornerRadius(6),
			VStack(
				Label("").
					BindText(trackTitle).
					FontWeight(600).
					Ellipsis(pango.EllipsizeEndValue).
					HAlign(gtk.AlignStartValue),
				Label("").
					WithCSSClass("dimmed").
					Ellipsis(pango.EllipsizeEndValue).
					HAlign(gtk.AlignStartValue).
					BindText(trackAlbum),
			).
				Spacing(3).
				VAlign(gtk.AlignCenterValue).
				HAlign(gtk.AlignStartValue).
				HExpand(true),
			Spacer().VExpand(false),
			HStack(
				Button().
					TooltipText(gettext.Get("Play / Pause")).
					WithCSSClass("flat").
					BindIconName(playPauseIcon).
					ConnectClicked(func(b gtk.Button) {
						player.PlayPause()
					}).
					BindSensitive(miniPlayerCanControl).
					ConnectConstruct(func(b *gtk.Button) {
						ptr := b.GoPointer()
						miniPlayerLoadingIconSub = player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
							if ps.Loading {
								schwifty.OnMainThreadOncePure(func() {
									b := gtk.ButtonNewFromInternalPtr(ptr)
									child := Spinner().ToGTK()
									b.SetChild(child)
								})
							}
							return signals.Continue
						})
					}).
					ConnectDestroy(func(w gtk.Widget) {
						player.PlaybackStateChanged.Unsubscribe(miniPlayerLoadingIconSub)
					}),
				Button().
					TooltipText(gettext.Get("Next")).
					WithCSSClass("flat").
					IconName("skip-forward-large-symbolic").
					ActionName("win.player.next").
					BindSensitive(miniPlayerCanControl),
			).
				Spacing(7).
				HAlign(gtk.AlignEndValue).
				VAlign(gtk.AlignCenterValue),
		).
			Spacing(16).
			Padding(10).
			Margin(2).
			Background("alpha(var(--view-fg-color), 0.1)").
			CornerRadius(12),
	).Margin(12)
}
