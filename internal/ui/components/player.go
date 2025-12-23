package components

import (
	"context"
	"strings"
	"time"

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

type Player struct {
	*gui.BoxImpl

	cover      *gtk.Image
	title      *gtk.Label
	artistName *gtk.Label
}

var playerCSS = cssutil.Applier("player", `
	.player-track-title {
		font-size: 24px;
		line-height: 1.2;
		font-weight: 800;
	}

	.player-track-artist {
		font-size: 16px;
		line-height: 1.2;
		font-weight: 700;
		text-decoration: underline;
		color: #1C71D8;
	}

	.player-button:not(:hover):not(:focus) {
		background-color: transparent;
	}
`)

func (m *Player) LoadCover(url string) {
	imgutil.AsyncGET(injector.MustInject[context.Context](), url, imgutil.ImageSetterFromImage(m.cover))
}

func NewPlayer() *Player {
	trackImage := gtk.NewImage()
	trackImage.SetPixelSize(319)

	title := gtk.NewLabel("[Track Title]")
	title.SetCSSClasses([]string{"player-track-title"})
	title.SetWrap(true)
	title.SetWrapMode(pango.WrapWordChar)
	title.SetHAlign(gtk.AlignCenter)

	artistName := gtk.NewLabel("[Artist Name]")
	artistName.SetCSSClasses([]string{"player-track-artist"})
	artistName.SetWrap(true)
	artistName.SetWrapMode(pango.WrapWordChar)

	slider := gtk.NewScale(gtk.OrientationHorizontal, nil)
	slider.SetHExpand(true)
	slider.SetRange(0, 100)
	slider.SetValue(50)
	slider.ConnectChangeValue(func(scroll gtk.ScrollType, value float64) (ok bool) {
		player.Scrub(value)
		return false
	})
	guiSlider := gui.Wrapper(slider)

	position := gui.Text("00:00")
	duration := gui.Text("00:00")

	playButton := gtk.NewButtonFromIconName("media-playback-start-symbolic")
	playButton.ConnectClicked(func() {
		player.PlayPause()
	})

	playerWidget := &Player{
		gui.VStack(
			gui.Spacer(),
			gui.AspectFrame(trackImage).CornerRadius(10).Overflow(gtk.OverflowHidden).HAlign(gtk.AlignCenter),
			gui.Wrapper(title).MarginTop(35),
			artistName,
			gui.HStack(
				gui.Wrapper(gtk.NewButtonFromIconName("audio-speakers-symbolic")).
					CSS(`button:not(:hover) { background-color: transparent; }`),
				gui.Wrapper(gtk.NewButtonFromIconName("heart-outline-thick-symbolic")).
					CSS(`button:not(:hover) { background-color: transparent; }`),
				gui.Wrapper(gtk.NewButtonFromIconName("compass2-symbolic")).
					CSS(`button:not(:hover) { background-color: transparent; }`),
				gui.Wrapper(gtk.NewButtonFromIconName("library-symbolic")).
					CSS(`button:not(:hover) { background-color: transparent; }`),
				gui.Wrapper(gtk.NewButtonFromIconName("folder-publicshare-symbolic")).
					CSS(`button:not(:hover) { background-color: transparent; }`),
			).
				HAlign(gtk.AlignCenter).
				Spacing(7).
				MarginTop(27),
			guiSlider.
				MarginTop(24).
				MarginLeft(24).
				MarginRight(24).
				CSS(`scale { background-color: transparent; }`),
			gui.HStack(
				position,
				gui.Spacer().VExpand(false),
				duration,
			).
				MarginLeft(24).
				MarginRight(24),
			gui.Text("Max").
				CSS("label { background-color: #DAC100; border-radius: 11px; padding: 2px 6px; }").
				HExpand(false).
				HAlign(gtk.AlignCenter),
			gui.HStack(
				gui.Wrapper(gtk.NewButtonFromIconName("media-playlist-shuffle-symbolic")).
					CSS(`button { min-width: 34px; min-height: 34px; } button:not(:hover) { background-color: transparent; }`),
				gui.Wrapper(gtk.NewButtonFromIconName("media-seek-backward-symbolic")).
					CSS(`button { min-width: 34px; min-height: 34px; } button:not(:hover) { background-color: transparent; }`),
				gui.Wrapper(playButton).
					CSS(`
						button {
							padding: 9px 32px;
							border-radius: 21px;
						}

						button:not(:hover) {
							background-color: var(--accent-bg-color);
						}
					`),
				gui.Wrapper(gtk.NewButtonFromIconName("media-seek-forward-symbolic")).
					CSS(`button { min-width: 34px; min-height: 34px; } button:not(:hover) { background-color: transparent; }`),
				gui.Wrapper(gtk.NewButtonFromIconName("media-playlist-repeat-song-symbolic")).
					CSS(`button { min-width: 34px; min-height: 34px; } button:not(:hover) { background-color: transparent; }`),
			).
				Spacing(7).
				HAlign(gtk.AlignCenter).MarginTop(36).MarginBottom(36),
			gui.Spacer(),
		),
		trackImage,
		title,
		artistName,
	}

	playerCSS(playerWidget)
	player.OnState.On(func(state player.State) bool {
		if state.Track != nil {
			playerWidget.LoadCover(tidalapi.ImageURL(state.Track.Album.Cover))
			artists := make([]string, 0)
			for _, artist := range state.Track.Artists {
				artists = append(artists, artist.Name)
			}
			playerWidget.artistName.SetText(strings.Join(artists, ", "))
			playerWidget.title.SetText(state.Track.Title)
		}

		if state.Duration > 0 {
			duration.GTKWidget().SetText((time.Duration(state.Duration) * time.Second).String())
			position.GTKWidget().SetText((time.Duration(state.Position) * time.Second).String())
			slider.SetValue(100.0 / float64(state.Duration) * float64(state.Position))
		} else {
			slider.SetValue(0)
		}

		switch state.Status {
		case player.StatusPlaying:
			playButton.SetSensitive(true)
			playButton.SetIconName("media-playback-pause-symbolic")
		case player.StatusPaused:
			playButton.SetSensitive(true)
			playButton.SetIconName("media-playback-start-symbolic")
		default:
			playButton.SetSensitive(false)
		}

		return signals.Continue
	})

	return playerWidget
}
