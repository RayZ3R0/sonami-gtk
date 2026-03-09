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
)

func secondaryControlsArtist(artist sonami.Artist, popover *gtk.PopoverMenu) schwifty.Box {
	return componentSecondaryControls(artist, popover, favouritebutton.FavouriteButton(appState.ArtistsCache, artist.ID()))
}

func NewArtist(artist sonami.Artist, playFunc func(), shuffleFunc func()) schwifty.Widget {
	coverUrl := artist.Cover(154)
	title := artist.Title()
	fans := gettext.GetN("%d Fan", "%d Fans", artist.FollowerCount(), artist.FollowerCount())

	menu := gio.NewMenu()
	tracking.SetFinalizer("Menu", menu)

	queueAllItem := gio.NewMenuItem(gettext.Get("Add Top Tracks to Queue"), "win.player.queue")
	queueAllItem.SetActionAndTargetValue("win.player.queue", glib.NewVariantString(fmt.Sprintf("artist/%s", artist.ID())))
	menu.AppendItem(queueAllItem)
	tracking.SetFinalizer("MenuItem", queueAllItem)

	popover := gtk.NewPopoverMenuFromModel(&menu.MenuModel)
	tracking.SetFinalizer("Popover", popover)

	return template(coverUrl, title, fans, artist.Description(), componentControls(playFunc, shuffleFunc), secondaryControlsArtist(artist, popover))
}
