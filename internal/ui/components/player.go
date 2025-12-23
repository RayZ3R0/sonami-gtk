package components

import (
	"context"

	"codeberg.org/dergs/tidalwave/pkg/gui"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
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
	trackImage.SetOverflow(gtk.OverflowHidden)

	frame := gtk.NewAspectFrame(0.5, 0, 1, false)
	frame.SetCSSClasses([]string{"player-media-frame"})
	frame.SetChild(trackImage)
	frame.SetHExpand(true)
	frame.SetOverflow(gtk.OverflowHidden)

	title := gtk.NewLabel("[Track Title]")
	title.SetCSSClasses([]string{"player-track-title"})

	artistName := gtk.NewLabel("[Artist Name]")
	artistName.SetCSSClasses([]string{"player-track-artist"})

	slider := gtk.NewScale(gtk.OrientationHorizontal, nil)
	slider.SetHExpand(true)
	slider.SetRange(0, 100)
	slider.SetValue(50)
	guiSlider := gui.Wrapper(slider)

	player := &Player{
		gui.VStack(
			frame,
			title,
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
				gui.Text("00:00"),
				gui.Spacer().VExpand(false),
				gui.Text("99:99"),
			).
				MarginLeft(24).
				MarginRight(24),
			gui.Text("Max").
				CSS("label { background-color: #DAC100; border-radius: 11px; padding: 2px 6px; }").
				HExpand(false).
				HAlign(gtk.AlignCenter),
			gui.Spacer(),
			gui.HStack(
				gui.Wrapper(gtk.NewButtonFromIconName("media-playlist-shuffle-symbolic")).
					CSS(`button { min-width: 34px; min-height: 34px; } button:not(:hover) { background-color: transparent; }`),
				gui.Wrapper(gtk.NewButtonFromIconName("media-seek-backward-symbolic")).
					CSS(`button { min-width: 34px; min-height: 34px; } button:not(:hover) { background-color: transparent; }`),
				gui.Wrapper(gtk.NewButtonFromIconName("media-playback-start-symbolic")).
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

	playerCSS(player)

	return player
}
