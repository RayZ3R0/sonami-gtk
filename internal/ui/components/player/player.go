package player

import (
	"fmt"
	"log/slog"

	"codeberg.org/dergs/tidalwave/internal/notifications"
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	v1 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v1"

	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var trackID = state.NewStateful("")

var (
	playbackQualityText  = state.NewStateful("High")
	playbackQualityClass = state.NewStateful("high")
)

func init() {
	player.OnTrackChanged.On(func(trackInfo player.TrackInformation) bool {
		slog.Debug(fmt.Sprintf("Quality: %s", trackInfo.Quality))
		trackID.SetValue(trackInfo.ID)
		switch trackInfo.Quality {
		case v1.AudioQualityLossy:
			playbackQualityText.SetValue("Low (96 kbps)")
			playbackQualityClass.SetValue("low")
		case v1.AudioQualityHighRes:
			playbackQualityText.SetValue("Low (320 kbps)")
			playbackQualityClass.SetValue("low")
		case v1.AudioQualityLossless:
			playbackQualityText.SetValue("High")
			playbackQualityClass.SetValue("high")
		case v1.AudioQualityHighResLossless:
			playbackQualityText.SetValue("Max")
			playbackQualityClass.SetValue("max")
		}

		return signals.Continue
	})
}

func NewPlayer() schwifty.Box {
	return VStack(
		Spacer().MinHeight(24),
		trackCover(),
		trackTitle(),
		trackArtists(),
		HStack(
			MenuButton().
				Popover(controlsVolumeSlider()).
				IconName("audio-speakers-symbolic").
				WithCSSClass("transparent"),
			Button().
				ActionName("unimplemented").
				IconName("heart-outline-thick-symbolic").
				WithCSSClass("transparent"),
			Button().
				ActionName("unimplemented").
				IconName("compass2-symbolic").
				WithCSSClass("transparent"),
			Button().
				ActionName("unimplemented").
				IconName("library-symbolic").
				WithCSSClass("transparent"),
			Button().
				IconName("folder-publicshare-symbolic").
				WithCSSClass("transparent").
				ConnectClicked(func(gtk.Button) {
					id := trackID.Value()
					if id == "" {
						notifications.OnToast.Notify("No track is currently playing.")
						return
					}

					display := gdk.DisplayGetDefault()
					defer display.Unref()
					clipboard := display.GetClipboard()
					defer clipboard.Unref()

					clipboard.SetText(fmt.Sprintf("https://tidal.com/track/%s?u", id))
					notifications.OnToast.Notify("Copied track URL to clipboard.")
				}),
		).
			HAlign(gtk.AlignCenterValue).
			Spacing(7).
			MarginTop(27),
		trackTimeline(),
		Label("Max").
			WithCSSClass("quality-label").
			BindText(playbackQualityText).
			FontSize(10).
			FontWeight(700).
			BindCSSClass(playbackQualityClass).
			CornerRadius(10).
			HPadding(8).
			VPadding(4).
			HExpand(false).
			HAlign(gtk.AlignCenterValue),
		HStack(
			Button().
				IconName("media-playlist-shuffle-symbolic").
				MinHeight(34).
				MinWidth(34).
				WithCSSClass("transparent").
				ActionName("win.player.shuffle"),
			Button().
				IconName("media-seek-backward-symbolic").
				MinHeight(34).
				MinWidth(34).
				WithCSSClass("transparent").
				ActionName("win.player.back"),
			controlsPlayPause(),
			Button().
				IconName("media-seek-forward-symbolic").
				MinHeight(34).
				MinWidth(34).
				WithCSSClass("transparent").
				ActionName("win.player.next"),
			Button().
				IconName("media-playlist-repeat-song-symbolic").
				MinHeight(34).
				MinWidth(34).
				WithCSSClass("transparent").
				ActionName("win.player.repeat"),
		).
			Spacing(7).
			HAlign(gtk.AlignCenterValue).
			MarginTop(42).
			MarginBottom(42),
		Spacer(),
	).
		WithCSSClass("player")
}
