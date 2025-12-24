package components

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/ui/signals"
	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
	"github.com/diamondburned/gotkit/gtkutil/imgutil"
	"github.com/infinytum/injector"
)

type LyricsPanel struct {
	*gui.BoxImpl

	lyrics     string
	lyricLines []LyricLine
}

type LyricLine struct {
	Timestamp  string
	TextWidget *gui.TextImpl
	BoxWidget  *gui.BoxImpl
}

var LyricsLineCSS = cssutil.Applier("lyrics-panel", `
box.lyric {
   background-color: alpha(var(--view-fg-color), 0.05);
   color: alpha(var(--view-fg-color), 0.5);
   transition-property: background, color;
   transition-duration: 300ms;
   transition-timing-function: ease-in-out;
}

box.lyric.active {
   background-color: alpha(var(--view-fg-color), 0.1);
   color: var(--view-fg-color);
}
`)

func NewLyricsPanel() *LyricsPanel {
	trackImage := gtk.NewImage()
	trackImage.SetPixelSize(54)
	trackImage.SetFromResource("/org/codeberg/dergs/tidalwave/icons/scalable/state/missing-album.svg")

	trackTitle := gui.Text("Track Name")
	trackArtists := gui.Text("Artist Name")

	lyricsView := gtk.NewScrolledWindow()
	lyricsView.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	lyricsView.SetVExpand(true)

	lyricsPanelWidget := &LyricsPanel{
		gui.VStack(
			gui.HStack(
				gui.AspectFrame(trackImage).
					Background("alpha(var(--view-fg-color), 0.1)").
					CornerRadius(6),
				gui.VStack(
					trackTitle.
						FontWeight(600).
						Ellipsis(pango.EllipsizeEnd).
						HAlign(gtk.AlignStart),
					trackArtists.
						Ellipsis(pango.EllipsizeEnd).
						HAlign(gtk.AlignStart),
				).
					VAlign(gtk.AlignCenter),
			).
				Spacing(16).
				Padding(12).
				MarginBottom(12).
				Background("alpha(var(--view-fg-color), 0.1)").
				CornerRadius(12),

			lyricsView,
		).
			Spacing(7).
			PaddingStart(16).
			PaddingEnd(16).
			PaddingTop(12).
			PaddingBottom(12),
		"",
		nil,
	}

	player.OnTrackChanged.On(func(trackInfo player.TrackInformation) bool {
		imgutil.AsyncGET(
			injector.MustInject[context.Context](),
			trackInfo.CoverURL,
			imgutil.ImageSetterFromImage(trackImage),
		)

		trackArtists.Text(trackInfo.ArtistNames())
		trackTitle.Text(trackInfo.Title)

		tidal := injector.MustInject[*tidalapi.TidalAPI]()
		track, err := tidal.OpenAPI.V2.Tracks.Track(context.Background(), trackInfo.ID, "lyrics")
		if err != nil {
			fmt.Errorf("", err)
			return signals.Continue
		}

		lyrics := ""
		isTimestamped := false

		for _, item := range track.Included {
			if lyricsAttribute := item.Attributes.Lyrics; lyricsAttribute != nil {
				if lyricsAttribute.LRCText != "" {
					isTimestamped = true
					lyrics = item.Attributes.Lyrics.LRCText
				} else if lyricsAttribute.Text != "" {
					isTimestamped = false
					lyrics = item.Attributes.Lyrics.Text
				} else {
					continue
				}
				break
			}
		}

		if lyrics == "" {
			lyricsView.SetChild(gui.Text("No lyrics available"))
		} else if isTimestamped {
			// Handle lyrics with timings
			// Remove timing tags and split into lines
			timingRegex := regexp.MustCompile(`\[(\d{2}:\d{2}\.\d{2})\](.*)`)
			splitLyrics := strings.Split(lyrics, "\n")

			lyricsPanelWidget.lyrics = lyrics

			// Clear any existing lyric lines
			lyricsPanelWidget.lyricLines = nil

			lines := []gtk.Widgetter{}

			for _, line := range splitLyrics {
				if strings.TrimSpace(line) == "" {
					continue
				}

				if ok, _ := regexp.MatchString(`^\[\d{2}:\d{2}\.\d{2}\]$`, line); ok {
					// Handle having none selected
					continue
				}

				matches := timingRegex.FindStringSubmatch(line)
				lyricText := line
				if len(matches) >= 3 {
					// Extract just the lyric text, removing the timing
					lyricText = strings.TrimSpace(matches[2])
				}

				textWidget := gui.Text(lyricText).
					HAlign(gtk.AlignCenter).
					VAlign(gtk.AlignCenter).
					FontSize(20).
					FontWeight(600).
					Wrap(true).
					Justify(gtk.JustifyCenter)

				boxWidget := gui.Box(textWidget).
					HExpand(true).
					PaddingTop(24).
					PaddingBottom(24).
					PaddingStart(16).
					PaddingEnd(16).
					CornerRadius(12).
					AddCSSClass("lyric")

				// Store the timestamp and text widget for later reference
				lyricsPanelWidget.lyricLines = append(lyricsPanelWidget.lyricLines, LyricLine{
					Timestamp:  matches[1],
					TextWidget: textWidget,
					BoxWidget:  boxWidget,
				})

				lines = append(lines, boxWidget)
			}

			newLyricsLines := gui.VStack(lines...).Spacing(12)
			lyricsView.SetChild(newLyricsLines)
			lyricsView.SetPolicy(gtk.PolicyNever, gtk.PolicyExternal)
		} else {
			// Handle lyrics without timings
			splitLyrics := strings.Split(lyrics, "\n")

			lyricsPanelWidget.lyrics = lyrics

			lines := []gtk.Widgetter{}

			for _, line := range splitLyrics {
				lines = append(lines, gui.Text(line))
			}

			newLyricsLines := gui.VStack(lines...)
			lyricsView.SetChild(newLyricsLines)
			lyricsView.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
		}

		return signals.Continue
	})

	var lastActiveLyricIndex = -1

	player.OnStateChanged.On(func(state player.State) bool {
		// Update lyrics highlighting based on current playback time
		if len(lyricsPanelWidget.lyricLines) > 0 {
			currentTime := state.Position
			var activeLyricIndex = -1

			for i, lyricLine := range lyricsPanelWidget.lyricLines {
				// Parse timestamp to seconds
				timeStr := lyricLine.Timestamp
				var minutes, seconds, centiseconds int
				if n := regexp.MustCompile(`(\d{2}):(\d{2})\.(\d{2})`).FindStringSubmatch(timeStr); len(n) == 4 {
					minutes = int(n[1][0]-'0')*10 + int(n[1][1]-'0')
					seconds = int(n[2][0]-'0')*10 + int(n[2][1]-'0')
					centiseconds = int(n[3][0]-'0')*10 + int(n[3][1]-'0')
				}
				lyricTime := minutes*60 + seconds + centiseconds/100.0

				// Determine end time (next lyric's timestamp or track end)
				var endTime int = currentTime + 1000 // Default to very far in future
				if i+1 < len(lyricsPanelWidget.lyricLines) {
					nextTimeStr := lyricsPanelWidget.lyricLines[i+1].Timestamp
					var nextMinutes, nextSeconds, nextCentiseconds int
					if n := regexp.MustCompile(`(\d{2}):(\d{2})\.(\d{2})`).FindStringSubmatch(nextTimeStr); len(n) == 4 {
						nextMinutes = int(n[1][0]-'0')*10 + int(n[1][1]-'0')
						nextSeconds = int(n[2][0]-'0')*10 + int(n[2][1]-'0')
						nextCentiseconds = int(n[3][0]-'0')*10 + int(n[3][1]-'0')
					}
					endTime = nextMinutes*60 + nextSeconds + nextCentiseconds/100.0
				}

				// Check if current time is within this lyric's time range
				if currentTime >= lyricTime && currentTime < endTime {
					lyricLine.BoxWidget.AddCSSClass("active")
					activeLyricIndex = i
				} else {
					// Dim non-current lyrics
					lyricLine.BoxWidget.RemoveCSSClass("active")
				}
			}

			// Auto-scroll to the active lyric only if the active index changed
			if activeLyricIndex >= 0 && activeLyricIndex != lastActiveLyricIndex {
				activeWidget := lyricsPanelWidget.lyricLines[activeLyricIndex].BoxWidget
				vadj := lyricsView.VAdjustment()

				// Get widget allocation and scroll window allocation
				widgetAlloc, _ := activeWidget.GTKWidget().ComputeBounds(activeWidget.GTKWidget().Parent())
				scollViewHeight := lyricsView.Height()
				currentScrollPosition := vadj.Value()

				// Check if the active lyric is outside the visible area
				widgetTop := float64(widgetAlloc.Y())
				widgetBottom := float64(widgetAlloc.Y() + widgetAlloc.Height())
				visibleTop := currentScrollPosition
				visibleBottom := currentScrollPosition + float64(scollViewHeight)

				fmt.Println(visibleTop, visibleBottom, widgetTop, widgetBottom)

				// Only scroll if the widget is outside the visible area
				if widgetTop < visibleTop || widgetBottom > visibleBottom {
					// Calculate the position to center the active lyric
					widgetCenter := float64(widgetAlloc.Y() + widgetAlloc.Height()/2)
					scrollCenter := float64(scollViewHeight / 2)
					targetPosition := widgetCenter - scrollCenter

					// Clamp the target position within valid bounds
					if targetPosition < vadj.Lower() {
						targetPosition = vadj.Lower()
					} else if targetPosition > vadj.Upper()-vadj.PageSize() {
						targetPosition = vadj.Upper() - vadj.PageSize()
					}

					vadj.SetValue(targetPosition)
				}

				lastActiveLyricIndex = activeLyricIndex
			}
		}
		return signals.Continue
	})

	LyricsLineCSS(lyricsPanelWidget)

	return lyricsPanelWidget
}
