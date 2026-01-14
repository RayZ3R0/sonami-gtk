package player

import (
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var trackID = ""

var (
	isControllable = state.NewStateful(false)
)

func init() {
	player.ControllableStateChanged.On(func(cs player.ControllableState) bool {
		isControllable.SetValue(cs.CanControl())
		return signals.Continue
	})

	player.TrackChanged.On(func(trackInfo *player.Track) bool {
		if trackInfo == nil {
			trackID = ""
		} else {
			trackID = trackInfo.ID
		}
		return signals.Continue
	})
}

func NewPlayer() schwifty.CenterBox {
	return CenterBox().WithCSSClass("player").
		CenterWidget(
			VStack(
				trackCover(),
				trackTitle().MarginTop(30),
				trackArtists(),
				controlsButtonRow().MarginTop(20),
				trackTimeline().MarginTop(20),
				HStack(
					Button().
						IconName("playlist-shuffle-symbolic").
						MinHeight(34).
						MinWidth(34).
						WithCSSClass("transparent").
						ActionName("win.player.shuffle"),
					Button().
						BindSensitive(isControllable).
						IconName("seek-backward-symbolic").
						MinHeight(34).
						MinWidth(34).
						WithCSSClass("transparent").
						ActionName("win.player.previous"),
					controlsPlayPause(),
					Button().
						BindSensitive(isControllable).
						IconName("seek-forward-symbolic").
						MinHeight(34).
						MinWidth(34).
						WithCSSClass("transparent").
						ActionName("win.player.next"),
					Button().
						MinHeight(34).
						MinWidth(34).
						WithCSSClass("transparent").
						ActionName("win.player.repeat").
						ConnectConstruct(func(b *gtk.Button) {
							ptr := b.GoPointer()
							player.RepeatModeChanged.On(func(state player.RepeatMode) bool {
								schwifty.OnMainThreadOnce(func(ptr uintptr) {
									b := gtk.ButtonNewFromInternalPtr(ptr)
									switch state {
									case player.RepeatModeNone:
										b.RemoveCssClass("color-accent")
										b.SetIconName("playlist-repeat-symbolic")
									case player.RepeatModeQueue:
										b.AddCssClass("color-accent")
										b.SetIconName("playlist-repeat-symbolic")
									case player.RepeatModeTrack:
										b.AddCssClass("color-accent")
										b.SetIconName("playlist-repeat-song-symbolic")
									}
								}, ptr)
								return signals.Continue
							})
						}),
				).
					Spacing(7).
					HAlign(gtk.AlignCenterValue).
					MarginTop(42),
			).VAlign(gtk.AlignCenterValue),
		).
		Margin(20)
}
