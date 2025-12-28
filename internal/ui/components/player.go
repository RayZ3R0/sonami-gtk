package components

import (
	"math"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

type Player struct {
	schwifty.Box

	cover      *gtk.Image
	title      *gtk.Label
	artistName *gtk.Label
}

func (m *Player) LoadCover(url string) {
	injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(url, m.cover)
}

func NewPlayer() *Player {
	trackImage := gtk.NewImage()
	trackImage.SetPixelSize(319)
	trackImage.SetOverflow(gtk.OverflowHiddenValue)
	trackImage.SetFromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg")

	title := Label("No Track").
		FontSize(24).
		FontWeight(800).
		LineHeight(1.2).
		HMargin(32).
		HAlign(gtk.AlignCenterValue).
		Ellipsis(pango.EllipsizeEndValue)()

	artistName := Label("").
		FontSize(16).
		FontWeight(700).
		LineHeight(1.2).
		Color("#1C71D8").
		TextDecoration("underline")()

	slider := Scale(gtk.OrientationHorizontalValue).
		HExpand(true).
		Range(0, 100).
		Value(50).
		ConnectChangeValue(func(r gtk.Range, st gtk.ScrollType, value float64) bool {
			player.Scrub(value)
			return false
		})()

	position := Label("00:00")()
	duration := Label("00:00")()

	playButton := Button().
		IconName("media-playback-start-symbolic").
		ConnectClicked(func(b gtk.Button) {
			player.PlayPause()
		})()

	volumeSlider := Scale(gtk.OrientationVerticalValue).
		Inverted(true).
		Range(0, 1).
		Value(0.5).
		ConnectChangeValue(func(r gtk.Range, st gtk.ScrollType, value float64) bool {
			// Cube the value to account for the logarithmic nature of human volume perception
			player.SetVolume(math.Pow(value, 3))
			return false
		})()

	player.OnVolumeChanged.On(func(volume float64) bool {
		// Cube root the volume to account for the logarithmic nature of human volume perception
		volumeSlider.SetValue(math.Pow(volume, 1.0/3.0))

		return signals.Continue
	})

	playerWidget := &Player{
		VStack(
			Spacer().MinHeight(24),
			AspectFrame(trackImage).
				CornerRadius(10).
				Overflow(gtk.OverflowHiddenValue).
				HAlign(gtk.AlignCenterValue).
				Background("alpha(var(--view-fg-color), 0.1)"),
			ManagedWidget(&title.Widget).MarginTop(35),
			ManagedWidget(&artistName.Widget),
			HStack(
				MenuButton().
					Popover(
						Popover(
							ManagedWidget(&volumeSlider.Widget).
								MinHeight(100).
								CSS(`scale:active { background-color: transparent; }`).
								HExpand(true),
						),
					).
					IconName("audio-speakers-symbolic").
					CSS("button:not(:hover) { background-color: transparent; }"),
				Button().
					IconName("heart-outline-thick-symbolic").
					CSS(`button:not(:hover) { background-color: transparent; }`),
				Button().
					IconName("compass2-symbolic").
					CSS(`button:not(:hover) { background-color: transparent; }`),
				Button().
					IconName("library-symbolic").
					CSS(`button:not(:hover) { background-color: transparent; }`),
				Button().
					IconName("folder-publicshare-symbolic").
					CSS(`button:not(:hover) { background-color: transparent; }`),
			).
				HAlign(gtk.AlignCenterValue).
				Spacing(7).
				MarginTop(27),
			ManagedWidget(&slider.Widget).
				Background("transparent").
				MarginTop(24).
				HMargin(24).
				HPadding(0),
			HStack(
				position,
				Spacer().VExpand(false),
				duration,
			).HMargin(24),
			Label("Max").
				Background("#DAC100").
				CornerRadius(10).
				HPadding(6).
				VPadding(2).
				HExpand(false).
				HAlign(gtk.AlignCenterValue),
			HStack(
				Button().
					IconName("media-playlist-shuffle-symbolic").
					MinHeight(34).
					MinWidth(34).
					CSS(`button:not(:hover) { background-color: transparent; }`),
				Button().
					IconName("media-seek-backward-symbolic").
					MinHeight(34).
					MinWidth(34).
					CSS(`button:not(:hover) { background-color: transparent; }`),
				ManagedWidget(&playButton.Widget).
					CornerRadius(21).
					HPadding(32).
					VPadding(9).
					CSS(`button:not(:hover) { background-color: var(--accent-bg-color); }`),
				Button().
					IconName("media-seek-forward-symbolic").
					MinHeight(34).
					MinWidth(34).
					CSS(`button:not(:hover) { background-color: transparent; }`),
				Button().
					IconName("media-playlist-repeat-song-symbolic").
					MinHeight(34).
					MinWidth(34).
					CSS(`button:not(:hover) { background-color: transparent; }`),
			).
				Spacing(7).
				HAlign(gtk.AlignCenterValue).
				MarginTop(42).
				MarginBottom(42),
			Spacer(),
		),
		trackImage,
		title,
		artistName,
	}

	player.OnTrackChanged.On(func(trackInfo player.TrackInformation) bool {
		if trackInfo.CoverURL != "" {
			playerWidget.LoadCover(trackInfo.CoverURL)
		} else {
			trackImage.SetFromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg")
		}
		duration.SetText(tidalapi.FormatDuration(int(trackInfo.Duration.Seconds())))
		playerWidget.artistName.SetText(trackInfo.ArtistNames())
		playerWidget.title.SetText(trackInfo.Title)
		return signals.Continue
	})

	player.OnStateChanged.On(func(state player.State) bool {
		position.SetText(tidalapi.FormatDuration(state.Position))
		if state.Duration > 0 {
			duration.SetText(tidalapi.FormatDuration(int(state.Duration)))
			slider.SetValue(100.0 / float64(state.Duration) * float64(state.Position))
		} else {
			slider.SetValue(0)
		}

		switch state.Status {
		case player.StatusPlaying:
			playButton.SetSensitive(true)
			playButton.SetIconName("media-playback-pause-symbolic")
		case player.StatusPaused, player.StatusStopped:
			playButton.SetSensitive(true)
			playButton.SetIconName("media-playback-start-symbolic")
		default:
			playButton.SetIconName("media-playback-start-symbolic")
			playButton.SetSensitive(false)
		}

		return signals.Continue
	})

	return playerWidget
}
