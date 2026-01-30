package player2

import (
	"context"
	"fmt"
	"strconv"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var (
	artistButtonSingle = Button().
				TooltipText(gettext.Get("Navigate to Artist")).
				IconName("music-artist2-symbolic").
				BindSensitive(isTrackLoadedState).
				WithCSSClass("transparent")
	artistButtonMultiple = MenuButton().
				TooltipText(gettext.Get("Navigate to Artist")).
				IconName("music-artist2-symbolic").
				BindSensitive(isTrackLoadedState).
				WithCSSClass("transparent")

	isTrackLoadedState = state.NewStateful(false)
	artistButtonState  = state.NewStateful[any](artistButtonSingle)
)

func init() {
	player.TrackChanged.OnLazy(func(t *player.Track) bool {
		isTrackLoadedState.SetValue(t != nil)

		if t != nil {
			if len(t.Artists) > 1 {
				menu := gio.NewMenu()
				defer menu.Unref()
				for _, artist := range t.Artists {
					menu.AppendItem(gio.NewMenuItem(artist.Attributes.Name, "win.route.artist::"+artist.ID))
				}
				artistButtonState.SetValue(artistButtonMultiple.MenuModel(&menu.MenuModel))
			} else {
				artistButtonState.SetValue(artistButtonSingle.ActionName("win.route.artist").ActionTargetValue(glib.NewVariantString(t.Artists[0].ID)))
			}
		}

		return signals.Continue
	})
}

func actionRow() schwifty.Box {
	return HStack(
		MenuButton().
			TooltipText(gettext.Get("Volume")).
			Popover(controlsVolumeSlider()).
			IconName("speakers-symbolic").
			WithCSSClass("transparent"),
		Button().
			TooltipText(gettext.Get("Add to Collection")).
			ActionName("unimplemented").
			IconName("heart-outline-thick-symbolic").
			// BindSensitive(isTrackLoaded). // Turn on once we implemented this feature
			WithCSSClass("transparent"),
		Button().
			TooltipText(gettext.Get("Navigate to Album")).
			IconName("cd-symbolic").
			BindSensitive(isTrackLoadedState).
			WithCSSClass("transparent").
			ConnectClicked(func(b gtk.Button) {
				track := player.TrackChanged.CurrentValue()
				if track == nil {
					return
				}

				var albumID string

				for _, album := range track.Albums {
					if album.Data.ID != "" {
						albumID = album.Data.ID
						break
					}
				}

				if albumID == "" {
					return
				}

				router.Navigate("album/" + albumID)
			}),
		Bin().BindChild(artistButtonState),
		Button().
			TooltipText(gettext.Get("Navigate to Track Mix")).
			IconName("compass2-symbolic").
			BindSensitive(isTrackLoadedState).
			WithCSSClass("transparent").
			ConnectClicked(func(b gtk.Button) {
				track := player.TrackChanged.CurrentValue()
				if track == nil {
					return
				}

				id, _ := strconv.Atoi(track.ID)

				tidalapi := injector.MustInject[*tidalapi.TidalAPI]()
				mix, err := tidalapi.V1.Tracks.Mix(context.Background(), id)

				if err != nil {
					if err.Error() == "unexpected status code: 404" {
						notifications.OnToast.Notify(gettext.Get("No mix found for the current track"))
						return
					}

					return
				}

				router.Navigate("playlist/" + mix.ID)
			}),
		Button().
			TooltipText(gettext.Get("Copy Track URL")).
			IconName("share-alt-symbolic").
			BindSensitive(isTrackLoadedState).
			WithCSSClass("transparent").
			ConnectClicked(func(gtk.Button) {
				display := gdk.DisplayGetDefault()
				defer display.Unref()
				clipboard := display.GetClipboard()
				defer clipboard.Unref()

				clipboard.SetText(fmt.Sprintf("https://tidal.com/track/%s?u", player.TrackChanged.CurrentValue().ID))
				notifications.OnToast.Notify(gettext.Get("Copied track URL to clipboard."))
			}),
	).HAlign(gtk.AlignCenterValue).Spacing(15).CSS("box { margin-top: -8px; margin-bottom: -8px; }")
}
