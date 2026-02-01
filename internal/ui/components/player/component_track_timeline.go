package player

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var (
	durationState       = state.NewStateful("00:00")
	positionState       = state.NewStateful("00:00")
	timelineSliderState = state.NewStateful(0.0)

	playbackQualityText  = state.NewStateful("Max")
	playbackQualityClass = state.NewStateful("max")
)

func init() {
	player.TrackChanged.On(func(trackInfo *player.Track) bool {
		schwifty.OnMainThreadOncePure(func() {
			if trackInfo == nil {
				durationState.SetValue("00:00")
			} else {
				durationState.SetValue(tidalapi.FormatDuration(trackInfo.Duration))
			}
		})
		return signals.Continue
	})

	player.PlaybackStateChanged.On(func(state *player.PlaybackState) bool {
		schwifty.OnMainThreadOncePure(func() {
			positionState.SetValue(tidalapi.FormatDuration(state.Position))
			if state.Duration > 0 {
				durationState.SetValue(tidalapi.FormatDuration(state.Duration))
				timelineSliderState.SetValue(100.0 / float64(state.Duration) * float64(state.Position))
			} else {
				timelineSliderState.SetValue(0)
			}
		})
		return signals.Continue
	})

	player.PlaybackQualityChanged.On(func(quality v1.AudioQuality) bool {
		schwifty.OnMainThreadOncePure(func() {
			switch quality {
			case v1.AudioQualityLossy:
				playbackQualityText.SetValue(gettext.Get("Low (96 kbps)"))
				playbackQualityClass.SetValue("low")
			case v1.AudioQualityHighRes:
				playbackQualityText.SetValue(gettext.Get("Low (320 kbps)"))
				playbackQualityClass.SetValue("low")
			case v1.AudioQualityLossless:
				playbackQualityText.SetValue(gettext.Get("High"))
				playbackQualityClass.SetValue("high")
			case v1.AudioQualityHighResLossless:
				playbackQualityText.SetValue(gettext.Get("Max"))
				playbackQualityClass.SetValue("max")
			}
		})
		return signals.Continue
	})
}

func trackTimeline() schwifty.Widget {
	overlay := gtk.NewOverlay()
	overlay.SetChild(VStack(
		Scale(gtk.OrientationHorizontalValue).
			BindSensitive(isControllableState).
			BindValue(timelineSliderState).
			Increments(1, 1).
			HExpand(true).
			Range(0, 100).
			Background("transparent").
			HPadding(0).
			ConnectChangeValue(func(r gtk.Range, st gtk.ScrollType, f float64) bool {
				player.SeekToPercent(f)
				return false
			}),
		HStack(
			Label("").BindText(positionState),
			Spacer().VExpand(false),
			Label("").BindText(durationState),
		),
	).MarginBottom(2).ToGTK())
	overlay.AddOverlay(
		MenuButton().
			WithCSSClass("quality-selector").
			Child(
				Label("").
					WithCSSClass("caption-heading").
					BindText(playbackQualityText).
					BindCSSClass(playbackQualityClass).
					CornerRadius(10).
					HPadding(8).
					VPadding(4),
			).
			Popover(
				Popover(
					VStack(
						Button().
							Child(VStack(
								Label("Low (96 kbps)").WithCSSClass("heading"),
								Label("96 kbps AAC").WithCSSClass("caption"),
							)).
							ConnectClicked(func(b gtk.Button) {

							}).
							WithCSSClass("low"),
						Button().
							Child(VStack(
								Label("Low (320 kbps)").WithCSSClass("heading"),
								Label("320 kbps AAC").WithCSSClass("caption"),
							)).
							ConnectClicked(func(b gtk.Button) {

							}).
							WithCSSClass("low"),
						Button().
							Child(VStack(
								Label("High").WithCSSClass("heading"),
								Label("16-bit 44.1 kHz FLAC").WithCSSClass("caption"),
							)).
							ConnectClicked(func(b gtk.Button) {

							}).
							WithCSSClass("high"),
						Button().
							Child(VStack(
								Label("Max").WithCSSClass("heading"),
								Label("24-bit 48 kHz FLAC").WithCSSClass("caption"),
							)).
							ConnectClicked(func(b gtk.Button) {

							}).
							WithCSSClass("max"),
					).
						WithCSSClass("selector").
						Spacing(8),
				),
			).
			HExpand(false).
			VAlign(gtk.AlignEndValue).
			HAlign(gtk.AlignCenterValue).
			ToGTK(),
	)

	return ManagedWidget(&overlay.Widget)
}
