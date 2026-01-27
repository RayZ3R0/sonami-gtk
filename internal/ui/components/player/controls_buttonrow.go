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
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("module", "ui", "component", "player")

var (
	isTrackLoaded = state.NewStateful(false)
)

func controlsButtonRow() schwifty.Box {
	player.TrackChanged.OnLazy(func(t *player.Track) bool {
		isTrackLoaded.SetValue(t != nil)
		return signals.Continue
	})

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
			WithCSSClass("transparent"),
		Button().
			TooltipText(gettext.Get("Navigate to Track Mix")).
			IconName("compass2-symbolic").
			BindSensitive(isTrackLoaded).
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

					logger.Error("error fetching mix for track", "error", err)
					return
				}

				router.Navigate("playlist/" + mix.ID)
			}),
		Button().
			TooltipText(gettext.Get("Navigate to Album")).
			IconName("library-symbolic").
			WithCSSClass("transparent").
			BindSensitive(isTrackLoaded).
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
		Button().
			TooltipText(gettext.Get("Share Track URL")).
			IconName("share-alt-symbolic").
			WithCSSClass("transparent").
			BindSensitive(isTrackLoaded).
			ConnectClicked(func(gtk.Button) {
				if trackID == "" {
					notifications.OnToast.Notify(gettext.Get("No track is currently playing."))
					return
				}

				display := gdk.DisplayGetDefault()
				defer display.Unref()
				clipboard := display.GetClipboard()
				defer clipboard.Unref()

				clipboard.SetText(fmt.Sprintf("https://tidal.com/track/%s?u", trackID))
				notifications.OnToast.Notify(gettext.Get("Copied track URL to clipboard."))
			}),
	).
		HAlign(gtk.AlignCenterValue).
		Spacing(7)
}
