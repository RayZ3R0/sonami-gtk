package tracklist_header

import (
	"fmt"

	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	appState "github.com/RayZ3R0/sonami-gtk/internal/state"
	favouritebutton "github.com/RayZ3R0/sonami-gtk/internal/ui/components/favourite_button"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/tracking"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
)

func secondaryControlsPlaylist(playlist sonami.Playlist, popover *gtk.PopoverMenu) schwifty.Box {
	var cache = appState.PlaylistsCache
	if playlist.IsMix() {
		cache = appState.MixesCache
	}
	return componentSecondaryControls(playlist, popover, favouritebutton.FavouriteButton(cache, playlist.ID()))
}

func NewPlaylist(playlist sonami.Playlist, playFunc func(), shuffleFunc func()) schwifty.Widget {
	coverUrl := playlist.Cover(154)
	title := playlist.Title()
	releaseDate := playlist.CreatedAt().Format("2006")
	creator := "TIDAL"
	if c := playlist.Creator(); c != nil {
		creator = c.Title()
	}
	var description string
	if playlist.IsMix() {
		description = gettext.Get("Personal Mix")
	} else {
		description = fmt.Sprintf("%d Track (%s)", playlist.Count(), tidalapi.FormatDuration(playlist.Duration()))
	}

	menu := gio.NewMenu()
	tracking.SetFinalizer("Menu", menu)

	queueAllItem := gio.NewMenuItem(gettext.Get("Add Playlist to Queue"), "win.player.queue")
	queueAllItem.SetActionAndTargetValue("win.player.queue", glib.NewVariantString(fmt.Sprintf("playlist/%s", playlist.ID())))
	menu.AppendItem(queueAllItem)
	tracking.SetFinalizer("MenuItem", queueAllItem)

	popover := gtk.NewPopoverMenuFromModel(&menu.MenuModel)
	tracking.SetFinalizer("Popover", popover)

	return template(coverUrl, title, releaseDate+" • "+creator, "\n"+description+"\n", componentControls(playFunc, shuffleFunc), secondaryControlsPlaylist(playlist, popover))
}
