package lyrics

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"codeberg.org/dergs/tidalwave/internal/g"
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/resources"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/graphene"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

var (
	coverState   = state.NewStateful[schwifty.Paintable](resources.MissingAlbum())
	trackTitle   = state.NewStateful[string]("")
	trackArtists = state.NewStateful[string]("")
)

var lyricsList = state.NewStateful[any](
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
		BindChild(lyricsList).
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

type highlightTiming struct {
	Start, End time.Duration
	Address    uintptr
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
	activeLyricIndex.SetValue(timing.Address)

	if !userManuallyScrolled.Value() {
		schwifty.OnMainThreadOnce(func(uintptr) {
			w := gtk.ButtonNewFromInternalPtr(timing.Address)
			scrollToLyric(w)
		}, 0)
	}
}

func init() {
	var activeIndexChangeOnPlayerUpdate *signals.Subscription

	player.OnTrackChanged.On(func(trackInfo player.TrackInformation) bool {
		defer runtime.GC()
		activeLyricIndex.SetValue(0)
		player.OnStateChanged.Unsubscribe(activeIndexChangeOnPlayerUpdate)
		activeIndexChangeOnPlayerUpdate = nil

		if trackInfo.ID == "" {
			coverState.SetValue(resources.MissingAlbum())
			trackTitle.SetValue("")
			trackArtists.SetValue("")

			lyricsList.SetValue(
				HStack(
					Label("No lyrics available").
						HAlign(gtk.AlignCenterValue).
						VAlign(gtk.AlignCenterValue).
						HExpand(true).
						VExpand(true),
				),
			)

			return signals.Continue
		}

		if texture, err := injector.MustInject[*imgutil.ImgUtil]().Load(trackInfo.CoverURL); err == nil {
			fmt.Println("setting texture:", texture)
			coverState.SetValue(texture)
			texture.Unref()
		}

		trackTitle.SetValue(trackInfo.Title)
		trackArtists.SetValue(trackInfo.ArtistNames())

		tidal := injector.MustInject[*tidalapi.TidalAPI]()
		track, err := tidal.OpenAPI.V2.Tracks.Track(context.Background(), trackInfo.ID, "lyrics")
		if err != nil {
			fmt.Errorf("", err)
			return signals.Continue
		}

		lyrics := ""
		isTimestamped := false

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

		if lyrics == "" {
			lyricsList.SetValue(
				HStack(
					Label("No lyrics available").
						HAlign(gtk.AlignCenterValue).
						VAlign(gtk.AlignCenterValue).
						HExpand(true).
						VExpand(true),
				),
			)

			return signals.Continue
		}

		lines := []any{}
		if isTimestamped {
			// Handle lyrics with timings
			// Remove timing tags and split into lines
			timingRegex := regexp.MustCompile(`\[(\d{2}:\d{2}\.\d{2})\](.*)`)
			splitLyrics := strings.Split(lyrics, "\n")
			timings := []highlightTiming{}

			for i, line := range splitLyrics {
				// Skip empty lines
				if strings.TrimSpace(line) == "" {
					continue
				}

				if ok, _ := regexp.MatchString(`^\[\d{2}:\d{2}\.\d{2}\]$`, line); ok {
					continue
				}

				matches := timingRegex.FindStringSubmatch(line)
				lyricText := line
				timestampStart := matches[1]
				if len(matches) >= 3 {
					// Extract just the lyric text, removing the timing
					lyricText = strings.TrimSpace(matches[2])
				}
				timeStart, _ := parseTimestamp(timestampStart)

				var timeEnd time.Duration = trackInfo.Duration
				if i+1 < len(splitLyrics) {
					nextLineMatches := timingRegex.FindStringSubmatch(splitLyrics[i+1])
					if len(nextLineMatches) >= 2 {
						timestampEnd := nextLineMatches[1]
						timeEnd, _ = parseTimestamp(timestampEnd)
					}
				}

				var classListener string

				boxWidget := lyricLine(lyricText, lyricTiming{
					timed:     true,
					timeStart: timeStart,
					timeEnd:   timeEnd,
				}).
					WithCSSClass("timed").
					ConnectConstruct(func(b *gtk.Button) {
						ptr := b.GoPointer()
						classListener = activeLyricIndex.AddCallback(func(newValue uintptr) {
							widget := gtk.ButtonNewFromInternalPtr(ptr)
							if newValue == ptr {
								widget.AddCssClass("active")
							} else {
								widget.RemoveCssClass("active")
							}
						})
					}).
					ConnectDestroy(func(w gtk.Widget) {
						activeLyricIndex.RemoveCallback(classListener)
					}).
					ConnectClicked(func(gtk.Button) {
						userManuallyScrolled.SetValue(false)
						player.SeekTo(timeStart)
					})()

				timings = append(timings, highlightTiming{
					Start:   timeStart,
					End:     timeEnd,
					Address: boxWidget.GoPointer(),
				})

				lines = append(lines, boxWidget)
			}

			activeIndexChangeOnPlayerUpdate = player.OnStateChanged.On(func(state player.State) (next bool) {
				next = signals.Continue
				if state.Status != player.StatusPlaying {
					return
				}

				hasActive := false

				for _, timing := range timings {
					if state.Position > timing.End {
						continue
					}

					if timing.Start <= state.Position {
						if activeLyricIndex.Value() != timing.Address {
							setNewIndex(timing)
						}

						hasActive = true
						continue
					}

					if timing.Start <= state.Position+player.UpdateInterval {
						fmt.Println("triggering in ", timing.Start-state.Position)
						time.AfterFunc(timing.Start-state.Position, func() {
							if activeLyricIndex.Value() != timing.Address {
								setNewIndex(timing)
							}
						})

						continue
					}
				}

				if !hasActive {
					activeLyricIndex.SetValue(0)
				}

				return
			})

			// Disallow user scrolling
			lyricsView().SetPolicy(gtk.PolicyNeverValue, gtk.PolicyExternalValue)
		} else {
			// Handle lyrics without timings
			splitLyrics := strings.Split(lyrics, "\n")

			for _, lyricText := range splitLyrics {
				if lyricText == "" {
					continue
				}

				boxWidget := lyricLine(lyricText, lyricTiming{timed: false})()
				lines = append(lines, boxWidget)
			}

			// Allow user to scroll
			lyricsView().SetPolicy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue)
		}

		lyricsList.SetValue(
			VStack(lines...).
				Spacing(12).
				HExpand(true).
				VExpand(true).
				ConnectDestroy(func(w gtk.Widget) {
					fmt.Println("SKIBIDI TOILET")
				}),
		)

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
				FromIconName("view-refresh-symbolic"),
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
