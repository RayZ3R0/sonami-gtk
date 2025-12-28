package player

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

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
		trackTimeline(),
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
			controlsPlayPause(),
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
	)
}
