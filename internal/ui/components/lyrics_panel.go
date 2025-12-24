package components

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/ui/signals"
	"codeberg.org/dergs/tidalwave/pkg/gui"
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
   color: alpha(var(--view-fg-color), 0.8);
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

	// const lyrics = "Oh oh oh oh oh\nOh oh oh oh oh\n\nI remember the night in the heat of July, we were strangers\nAfter the thunder cracked the lightning was kind of amazing\nI got lost in your glow, felt so high even though I was sober\nYou had a reckless style, your smile was wearing a gold dust\nThe way we hold on it's a miracle\nThis is a love song forever in making\n\nOh, the years go by but you still light up, light up my heart\nOh, you know time flies but you still light up, light up my dark\nIt's been the wildest dream, a masterpiece, let's ride around the sun\nDarling, after all this time, you're still the one\n\nOh oh oh oh oh\nOh oh oh oh oh\n\nPeople come, people go, nothing lasts anymore\nIt's a strange world\nBut dreaming side by side, I wake up believing in angels\nYes, you know the way sparks still fly\nHow you read my mind, it's electric (it's electric)\nAnd when I leave this life\nI'll know loving you was always the best bit\nThe way we hold on it's a miracle\nThis is a love song forever in making\n\nOh, years go by but you still light up, light up my heart\nOh, you know time flies but you still light up, light up my dark\nIt's been the wildest dream, a masterpiece, let's ride around the sun\nDarling, after all this time, you're still the one\n\nStill the one (light up, light up my heart)\n\nIt's been the wildest dream, a masterpiece, let's ride around the sun\nDarling, after all this time\nYou know I need you by my side, I tell ya\nAfter all this time, you're still the one"
	const lyrics = "[00:01.50]Oh oh oh oh oh\n[00:04.82]Oh oh oh oh oh\n\n[00:08.33]I remember the night in the heat of July, we were strangers\n[00:15.48]After the thunder cracked the lightning was kind of amazing\n[00:23.25]I got lost in your glow, felt so high even though I was sober\n[00:30.14]You had a reckless style, your smile was wearing a gold dust\n[00:37.60]The way we hold on it's a miracle\n[00:41.39]This is a love song forever in making\n\n[00:45.52]Oh, the years go by but you still light up, light up my heart\n[00:52.88]Oh, you know time flies but you still light up, light up my dark\n[00:59.65]It's been the wildest dream, a masterpiece, let's ride around the sun\n[01:06.62]Darling, after all this time, you're still the one\n\n[01:11.22]Oh oh oh oh oh\n[01:12.62]Oh oh oh oh oh\n\n[01:14.69]People come, people go, nothing lasts anymore\n[01:21.27]It's a strange world\n[01:24.75]But dreaming side by side, I wake up believing in angels\n[01:31.14]Yes, you know the way sparks still fly\n[01:34.11]How you read my mind, it's electric (it's electric)\n[01:39.55]And when I leave this life\n[01:41.20]I'll know loving you was always the best bit\n[01:46.68]The way we hold on it's a miracle\n[01:50.36]This is a love song forever in making\n\n[01:54.68]Oh, years go by but you still light up, light up my heart\n[02:01.54]Oh, you know time flies but you still light up, light up my dark\n[02:08.38]It's been the wildest dream, a masterpiece, let's ride around the sun\n[02:15.52]Darling, after all this time, you're still the one\n[02:22.43]\n[02:26.35]Still the one (light up, light up my heart)\n\n[02:33.86]It's been the wildest dream, a masterpiece, let's ride around the sun\n[02:40.91]Darling, after all this time\n[02:44.15]You know I need you by my side, I tell ya\n[02:49.15]After all this time, you're still the one\n[02:53.33]"

	player.OnTrackChanged.On(func(track player.TrackInformation) bool {
		imgutil.AsyncGET(
			injector.MustInject[context.Context](),
			track.CoverURL,
			imgutil.ImageSetterFromImage(trackImage),
		)

		trackArtists.Text(track.ArtistNames())
		trackTitle.Text(track.Title)

		if lyrics == "" {
			lyricsView.SetChild(gui.Text("No lyrics available"))
		} else if ok, _ := regexp.MatchString(`\d{2}:\d{2}\.\d{2}\]`, lyrics); ok {
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
