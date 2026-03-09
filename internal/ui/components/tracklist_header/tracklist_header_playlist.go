package tracklist_header

import (
	"fmt"

	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	appState "github.com/RayZ3R0/sonami-gtk/internal/state"
	// favouritebutton "github.com/RayZ3R0/sonami-gtk/internal/ui/components/favourite_button" // deferred: local favourites
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/tracking"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

func secondaryControlsPlaylist(playlist sonami.Playlist, popover *gtk.PopoverMenu) schwifty.Box {
	// Favourite button deferred — see hifi/deferred_features.md
	// Keeping cache refs for future local-DB implementation:
	_ = appState.MixesCache
	_ = appState.PlaylistsCache
	return componentSecondaryControls(playlist, popover, nil)
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
