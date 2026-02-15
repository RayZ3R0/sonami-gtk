package player

import (
	"context"
	"fmt"
	"log/slog"
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
	"codeberg.org/dergs/tonearm/pkg/tonearm"
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
				WithCSSClass("flat")
	artistButtonMultiple = MenuButton().
				TooltipText(gettext.Get("Navigate to Artist")).
				IconName("music-artist2-symbolic").
				BindSensitive(isTrackLoadedState).
				WithCSSClass("flat")

	isTrackLoadedState = state.NewStateful(false)
	artistButtonState  = state.NewStateful[any](artistButtonSingle)
)

func init() {
	player.TrackChanged.OnLazy(func(t tonearm.Track) bool {
		isTrackLoadedState.SetValue(t != nil)

		if t != nil {

			artistPaginator, err := t.Artists()
			if err != nil {
				slog.Error("Failed to load artists", "error", err)
				return signals.Continue
			}

			artists, err := artistPaginator.GetAll()
			if err != nil {
				slog.Error("Failed to load all artists", "error", err)
				return signals.Continue
			}

			if len(artists) > 1 {
				menu := gio.NewMenu()
				defer menu.Unref()
				for _, artist := range artists {
					menu.AppendItem(gio.NewMenuItem(artist.Name(), "win.route.artist::"+artist.ID()))
				}
				artistButtonState.SetValue(artistButtonMultiple.MenuModel(&menu.MenuModel))
			} else {
				artistButtonState.SetValue(artistButtonSingle.ActionName("win.route.artist").ActionTargetValue(glib.NewVariantString(artists[0].ID())))
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
			WithCSSClass("flat"),
		Button().
			TooltipText(gettext.Get("Add to Collection")).
			ActionName("unimplemented").
			IconName("heart-outline-thick-symbolic").
			// BindSensitive(isTrackLoaded). // Turn on once we implemented this feature
			WithCSSClass("flat"),
		Button().
			TooltipText(gettext.Get("Navigate to Album")).
			IconName("cd-symbolic").
			BindSensitive(isTrackLoadedState).
			WithCSSClass("flat").
			ConnectClicked(func(b gtk.Button) {
				track := player.TrackChanged.CurrentValue()
				if track == nil {
					return
				}

				album, err := track.Album()
				if err != nil {
					slog.Error("Failed to load album", "error", err)
					return
				}

				router.Navigate("album/" + album.ID())
			}),
		Bin().BindChild(artistButtonState),
		Button().
			TooltipText(gettext.Get("Navigate to Track Mix")).
			IconName("compass2-symbolic").
			BindSensitive(isTrackLoadedState).
			WithCSSClass("flat").
			ConnectClicked(func(b gtk.Button) {
				track := player.TrackChanged.CurrentValue()
				if track == nil {
					return
				}

				id, _ := strconv.Atoi(track.ID())

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
			WithCSSClass("flat").
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
