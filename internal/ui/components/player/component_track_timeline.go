package player

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/settings"
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

func hideCheckmarkHook(quality v1.AudioQuality) func(*gtk.Label) {
	return func(l *gtk.Label) {
		settings.Player().ConnectAudioQualityChanged(
			func(aq v1.AudioQuality) bool {
				if aq == quality {
					l.Show()
				} else {
					l.Hide()
				}

				return signals.Continue
			},
		)

		if settings.Player().GetAudioQuality() == quality {
			l.Show()
		} else {
			l.Hide()
		}
	}
}

func makeQualitySelectEntry(quality v1.AudioQuality, css, label, details string, popover **gtk.Popover) schwifty.Button {
	return Button().
		Child(
			HStack(
				VStack(
					Label(label).WithCSSClass("heading").HAlign(gtk.AlignStartValue),
					Label(details).WithCSSClass("caption").HAlign(gtk.AlignStartValue),
				),
				Spacer(),
				Label("✓").ConnectConstruct(hideCheckmarkHook(quality)),
			).Spacing(10),
		).
		ConnectClicked(func(b gtk.Button) {
			settings.Player().SetAudioQuality(quality)
			(*popover).Hide()
		}).
		WithCSSClass(css)
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

	var popover *gtk.Popover
	popover = Popover(
		VStack(
			makeQualitySelectEntry(v1.AudioQualityLossy, "low", "Low (96 kbps)", "96 kbps AAC", &popover),
			makeQualitySelectEntry(v1.AudioQualityHighRes, "low", "Low (320 kbps)", "320 kbps AAC", &popover),
			makeQualitySelectEntry(v1.AudioQualityLossless, "high", "High", "16-bit 44.1 kHz FLAC", &popover),
			makeQualitySelectEntry(v1.AudioQualityHighResLossless, "max", "Max", "24-bit 48 kHz FLAC", &popover),
		).
			WithCSSClass("selector").
			Spacing(8),
	)()

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
			Popover(popover).
			HExpand(false).
			VAlign(gtk.AlignEndValue).
			HAlign(gtk.AlignCenterValue).
			ToGTK(),
	)

	return ManagedWidget(&overlay.Widget)
}
