package lyrics

import (
	"log/slog"
	"runtime"
	"time"

	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/player"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	"github.com/RayZ3R0/sonami-gtk/internal/ui/components/sidebar"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/utils/weak"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

var (
	lyricsPanel = state.NewStateful[any](lyricsStatusNoSongPlaying)

	logger = slog.With("module", "ui").WithGroup("ui").With("component", "lyrics")

	activeLyricIndex                = state.NewStateful[uintptr](0)
	activeIndexChangeOnPlayerUpdate *signals.Subscription
)

type highlightTiming struct {
	Start, End time.Duration
	Ref        weak.WidgetRef
}

type lyricsStatus int

const (
	lyricsStatusNoSongPlaying lyricsStatus = iota
	lyricsStatusError
	lyricsStatusLoading
	lyricsStatusNoLyrics
	lyricsStatusLoaded
)

func (s lyricsStatus) ToGTK() *gtk.Widget {
	sp := StatusPage().VExpand(true)

	switch s {
	case lyricsStatusNoSongPlaying:
		sp = sp.Title(gettext.Get("No Song Currently Playing")).IconName("no-track-symbolic")
	case lyricsStatusError:
		sp = sp.Title(gettext.Get("Error Fetching Lyrics")).IconName("dialog-error-symbolic")
	case lyricsStatusLoading:
		sp = sp.Title(gettext.Get("Loading Lyrics")).Loading()
	case lyricsStatusNoLyrics:
		sp = sp.Title(gettext.Get("No Lyrics Available")).IconName("no-lyrics-symbolic")
	case lyricsStatusLoaded:
		return &lyricsOverlay().Widget
	}

	return sp.ToGTK()
}

func init() {
	player.TrackChanged.OnLazy(func(trackInfo sonami.Track) bool {
		adj := lyricsView().GetVadjustment()
		defer adj.Unref()
		adj.SetValue(0)
		userManuallyScrolled.SetValue(false)

		lyricsPanel.SetValue(lyricsStatusLoading)
		lyricsList.SetValue(nil)
		defer runtime.GC()
		activeLyricIndex.SetValue(0)
		player.PlaybackStateChanged.Unsubscribe(activeIndexChangeOnPlayerUpdate)
		activeIndexChangeOnPlayerUpdate = nil

		if trackInfo == nil {
			lyricsPanel.SetValue(lyricsStatusNoSongPlaying)

			return signals.Continue
		}

		lyrics, isTimestamped, err := getLyrics(trackInfo.ID())
		if err != nil {
			logger.Error("Error while fetching lyrics", "error", err)
			lyricsPanel.SetValue(lyricsStatusError)

			return signals.Continue
		}

		if lyrics == "" {
			lyricsPanel.SetValue(lyricsStatusNoLyrics)

			return signals.Continue
		}

		setLyrics(isTimestamped, lyrics, trackInfo.Duration())
		lyricsPanel.SetValue(lyricsStatusLoaded)

		return signals.Continue
	})
}

func NewLyricsPanel() schwifty.Box {
	trackLoaded := state.NewStateful(false)

	player.TrackChanged.On(func(t sonami.Track) bool {
		trackLoaded.SetValue(t != nil)
		return signals.Continue
	})

	return VStack(
		sidebar.MiniPlayer().BindVisible(trackLoaded),
		Bin().BindChild(lyricsPanel),
	).
		Spacing(7).
		WithCSSClass("lyrics-panel")
}
