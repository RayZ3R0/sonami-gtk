package lyrics

import (
	"context"
	"log/slog"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/graphene"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

var (
	coverState   = state.NewStateful[schwifty.Paintable](resources.MissingAlbum())
	trackTitle   = state.NewStateful[string]("")
	trackArtists = state.NewStateful[string]("")
	loadingState = state.NewStateful(false)

	logger = slog.With("module", "lyrics")
)

type loadingWidget struct {
	*state.State[any]
}

func NewLoadingWidget(defaultState any) *loadingWidget {
	lw := &loadingWidget{
		State: state.NewStateful[any](defaultState),
	}

	return lw
}

func (lw *loadingWidget) SetValue(v any) {
	if v == nil {
		lw.State.SetValue(
			HStack(
				Clamp().
					MaximumSize(50).
					Child(Spinner()).
					HExpand(true).
					VExpand(true),
			).
				HExpand(true).
				VExpand(true),
		)
	} else {
		lw.State.SetValue(v)
	}
}

var lyricsList = NewLoadingWidget(
	HStack(
		Label("No lyrics available").
			HAlign(gtk.AlignCenterValue).
			VAlign(gtk.AlignCenterValue).
			HExpand(true).
			VExpand(true),
	),
)

var (
	userManuallyScrolled = state.NewStateful(false)
	scrollIsProgrammatic bool
)

var lyricsView = g.Lazy(func() (w *gtk.ScrolledWindow) {
	w = ScrolledWindow().
		BindChild(lyricsList.State).
		Policy(gtk.PolicyNeverValue, gtk.PolicyExternalValue)()

	cb := func(adj gtk.Adjustment) {
		if !scrollIsProgrammatic {
			userManuallyScrolled.SetValue(true)
		}
	}

	adj := w.GetVadjustment()
	defer adj.Unref()
	adj.ConnectValueChanged(&cb)

	return
})

func parseTimestamp(timestamp string) (resTime time.Duration, ok bool) {
	timestampRegex := regexp.MustCompile(`(\d{2}):(\d{2})\.(\d{2})`)
	parts := timestampRegex.FindStringSubmatch(timestamp)

	if len(parts) != 4 {
		ok = false
		return
	}

	ok = true

	minutes, _ := strconv.Atoi(parts[1])
	seconds, _ := strconv.Atoi(parts[2])
	centiseconds, _ := strconv.Atoi(parts[3])

	resTime = time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(centiseconds)*time.Millisecond*10

	return
}

var activeLyricIndex = state.NewStateful[uintptr](0)
var activeIndexChangeOnPlayerUpdate *signals.Subscription

type highlightTiming struct {
	Start, End time.Duration
	Ref        *tracking.WeakRef
}

func scrollToLyric(w *gtk.Button) {
	scrollIsProgrammatic = true
	parentWidget := w.GetParent()
	if parentWidget == nil {
		return
	}

	defer parentWidget.Unref()

	var bounds graphene.Rect
	w.ComputeBounds(parentWidget, &bounds)
	vadj := lyricsView().GetVadjustment()
	defer vadj.Unref()
	scrollViewHeight := lyricsView().GetHeight()

	// Calculate the position to center the active lyric
	widgetCenter := float64(bounds.GetY() + bounds.GetHeight()/2)
	scrollCenter := float64(scrollViewHeight / 2)
	targetPosition := widgetCenter - scrollCenter

	// Clamp the target position within valid bounds
	if targetPosition < vadj.GetLower() {
		targetPosition = vadj.GetLower()
	} else if targetPosition > vadj.GetUpper()-vadj.GetPageSize() {
		targetPosition = vadj.GetUpper() - vadj.GetPageSize()
	}

	vadj.SetValue(targetPosition)
	scrollIsProgrammatic = false
}

func setNewIndex(timing highlightTiming) {
	object := timing.Ref.Get()
	if object == nil {
		return
	}

	ptr := object.GoPointer()

	if activeLyricIndex.Value() != ptr {
		activeLyricIndex.SetValue(ptr)
	}

	if !userManuallyScrolled.Value() {
		schwifty.OnMainThreadOnce(func(uintptr) {
			w := gtk.ButtonNewFromInternalPtr(ptr)
			scrollToLyric(w)
			object.Unref()
		}, 0)
	} else {
		object.Unref()
	}
}

func parseLRCLyrics(lyrics string, trackInfo tonearm.Track) (lines []any) {
	// Handle lyrics with timings
	// Remove timing tags and split into lines
	timingRegex := regexp.MustCompile(`\[(\d{2}:\d{2}\.\d{2})\](.*)`)
	splitLyrics := strings.Split(lyrics, "\n")
	timings := []highlightTiming{}

	for i, line := range splitLyrics {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			lines = append(lines, Box(gtk.OrientationVerticalValue))
			continue
		}

		matches := timingRegex.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		timestampStart := matches[1]
		timeStart, _ := parseTimestamp(timestampStart)

		var timeEnd time.Duration = trackInfo.Duration()

		if i+1 < len(splitLyrics) {
			offset := 1
			nextLineMatches := timingRegex.FindStringSubmatch(splitLyrics[i+offset])

			for len(nextLineMatches) < 2 && i+offset+1 < len(splitLyrics) {
				offset++
				nextLineMatches = timingRegex.FindStringSubmatch(splitLyrics[i+offset])
			}

			if len(nextLineMatches) >= 2 {
				timestampEnd := nextLineMatches[1]
				timeEnd, _ = parseTimestamp(timestampEnd)
			}
		}

		if matches[2] == "" {
			timings = append(timings, highlightTiming{
				Start: timeStart,
				End:   timeEnd,
				Ref:   new(tracking.WeakRef),
			})

			continue
		}

		lyricText := line
		if len(matches) >= 3 {
			// Extract just the lyric text, removing the timing
			lyricText = strings.TrimSpace(matches[2])
		}

		boxWidget := lyricLine(lyricText, &lyricTiming{
			timeStart: timeStart,
			timeEnd:   timeEnd,
		})()

		lines = append(lines, boxWidget)

		timings = append(timings, highlightTiming{
			Start: timeStart,
			End:   timeEnd,
			Ref:   tracking.NewWeakRef(boxWidget),
		})
	}

	activeIndexChangeOnPlayerUpdate = player.PlaybackStateChanged.On(func(state *player.PlaybackState) (next bool) {
		next = signals.Continue
		if state.Status != player.PlaybackStatusPlaying {
			return
		}

		hasActive := false

		for _, timing := range timings {
			if state.Position > timing.End {
				continue
			}

			if timing.Start <= state.Position {
				timing.Ref.Use(func(obj *gobject.Object) {
					if activeLyricIndex.Value() != obj.Ptr {
						setNewIndex(timing)
					}
				})

				hasActive = true
				continue
			}

			if timing.Start <= state.Position+player.UpdateInterval {
				logger.Debug("next lyric line scheduled", "timing", timing.Start-state.Position)
				timing.Ref.Use(func(obj *gobject.Object) {
					time.AfterFunc(timing.Start-state.Position, func() {
						if activeLyricIndex.Value() != obj.Ptr {
							setNewIndex(timing)
						}
					})
				})

				continue
			}
		}

		if !hasActive && activeLyricIndex.Value() != 0 {
			schwifty.OnMainThreadOncePure(func() {
				activeLyricIndex.SetValue(0)
			})
		}

		return
	})

	// Disallow user scrolling
	schwifty.OnMainThreadOncePure(func() {
		lyricsView().SetPolicy(gtk.PolicyNeverValue, gtk.PolicyExternalValue)
	})

	return
}

func parseUntimedLyrics(lyrics string) (lines []any) {
	// Handle lyrics without timings
	splitLyrics := strings.SplitSeq(lyrics, "\n")

	for lyricText := range splitLyrics {
		if strings.TrimSpace(lyricText) == "" {
			lines = append(lines, Box(gtk.OrientationVerticalValue))
			continue
		}

		boxWidget := lyricLine(lyricText, nil)()
		lines = append(lines, boxWidget)
	}

	// Allow user to scroll
	schwifty.OnMainThreadOncePure(func() {
		lyricsView().SetPolicy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue)
	})

	return
}

func loadMiniplayerState(trackInfo tonearm.Track) {
	go func() {
		if texture, err := injector.MustInject[*imgutil.ImgUtil]().Load(trackInfo.Cover(80)); err == nil {
			schwifty.OnMainThreadOncePure(func() {
				coverState.SetValue(texture)
				texture.Unref()
			})
		}
	}()

	schwifty.OnMainThreadOncePure(func() {
		trackTitle.SetValue(tonearm.FormatTitle(trackInfo))
		trackArtists.SetValue(strings.Join(trackInfo.Artists().Names(), ", "))
	})
}

func getLyrics(ID string) (lyrics string, isTimestamped bool, err error) {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	var track *openapi.Track
	track, err = tidal.OpenAPI.V2.Tracks.Track(context.Background(), ID, "lyrics")
	if err != nil {
		return
	}

	for _, lyric := range track.Included.PlainLyrics(track.Data.Relationships.Lyrics.Data...) {
		if lyric.Attributes.LRCText != "" {
			isTimestamped = true
			lyrics = lyric.Attributes.LRCText
		} else if lyric.Attributes.Text != "" {
			isTimestamped = false
			lyrics = lyric.Attributes.Text
		}
		break
	}

	return
}

func setLyricsEmptyState(msg string) {
	schwifty.OnMainThreadOncePure(func() {
		lyricsList.SetValue(
			HStack(
				Label(msg).
					HAlign(gtk.AlignCenterValue).
					VAlign(gtk.AlignCenterValue).
					HExpand(true).
					VExpand(true),
			).
				VExpand(true),
		)
	})
}

func init() {
	player.TrackChanged.On(func(trackInfo tonearm.Track) bool {
		lyricsList.SetValue(nil)
		defer runtime.GC()
		activeLyricIndex.SetValue(0)
		player.PlaybackStateChanged.Unsubscribe(activeIndexChangeOnPlayerUpdate)
		activeIndexChangeOnPlayerUpdate = nil

		if trackInfo == nil {
			schwifty.OnMainThreadOncePure(func() {
				coverState.SetValue(resources.MissingAlbum())
				trackTitle.SetValue("")
				trackArtists.SetValue("")

				setLyricsEmptyState(gettext.Get("No song currently playing"))
			})

			return signals.Continue
		}

		loadMiniplayerState(trackInfo)
		lyrics, isTimestamped, err := getLyrics(trackInfo.ID())
		if err != nil {
			logger.Error("Error while fetching lyrics", "error", err)
			schwifty.OnMainThreadOncePure(func() {
				setLyricsEmptyState(gettext.Get("Error fetching lyrics"))
			})

			return signals.Continue
		}

		if lyrics == "" {
			setLyricsEmptyState(gettext.Get("No lyrics available"))

			return signals.Continue
		}

		lines := []any{}
		if isTimestamped {
			lines = parseLRCLyrics(lyrics, trackInfo)
		} else {
			lines = parseUntimedLyrics(lyrics)
		}

		schwifty.OnMainThreadOncePure(func() {
			lyricsList.SetValue(
				VStack(lines...).
					Spacing(12).
					HExpand(true).
					VExpand(true),
			)
		})

		return signals.Continue
	})

}

func NewLyricsPanel() schwifty.Box {
	overlay := gtk.NewOverlay()
	overlay.SetChild(&lyricsView().Widget)
	overlay.AddOverlay(&Button().
		HAlign(gtk.AlignEndValue).
		VAlign(gtk.AlignEndValue).
		Margin(7).
		TooltipText(gettext.Get("Sync with track")).
		BindVisible(userManuallyScrolled).
		ConnectClicked(func(b gtk.Button) {
			if activeLyricButtonPtr := activeLyricIndex.Value(); activeLyricButtonPtr != 0 {
				w := gtk.ButtonNewFromInternalPtr(activeLyricButtonPtr)
				scrollToLyric(w)
			}
			userManuallyScrolled.SetValue(false)
		}).
		Child(
			Image().
				FromIconName("arrow-circular-top-right-symbolic"),
		)().
		Widget)
	return VStack(
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
					BindText(trackArtists).
					Ellipsis(pango.EllipsizeEndValue).
					HAlign(gtk.AlignStartValue),
			).
				VAlign(gtk.AlignCenterValue),
		).
			Spacing(16).
			Padding(12).
			MarginBottom(12).
			Background("alpha(var(--view-fg-color), 0.1)").
			CornerRadius(12),
		overlay,
	).
		Spacing(7).
		PaddingStart(16).
		PaddingEnd(16).
		PaddingTop(12).
		PaddingBottom(12).
		WithCSSClass("lyrics-panel")
}
