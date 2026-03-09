package tracklist_header

import (
	"codeberg.org/puregotk/puregotk/v4/gio"
	"codeberg.org/puregotk/puregotk/v4/glib"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/tracking"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

type shareablePlaybackSource interface {
	sonami.PlaybackSource
	sonami.Shareable
}

func secondaryControlsCollection(playbackSource shareablePlaybackSource, popover *gtk.PopoverMenu) schwifty.Box {
	return componentSecondaryControls(playbackSource, popover)
}

func NewCollection(playbackSource shareablePlaybackSource, playFunc func(), shuffleFunc func()) schwifty.Widget {
	coverUrl := playbackSource.Cover(154)
	title := playbackSource.Title()
	subtitle := gettext.Get("My Collection")

	menu := gio.NewMenu()
	tracking.SetFinalizer("Menu", menu)

	queueAllItem := gio.NewMenuItem(gettext.Get("Add My Tracks to Queue"), "win.player.queue")
	queueAllItem.SetActionAndTargetValue("win.player.queue", glib.NewVariantString("my_collection/tracks"))
	menu.AppendItem(queueAllItem)
	tracking.SetFinalizer("MenuItem", queueAllItem)

	popover := gtk.NewPopoverMenuFromModel(&menu.MenuModel)
	tracking.SetFinalizer("Popover", popover)

	return template(coverUrl, title, subtitle, "", componentControls(playFunc, shuffleFunc), secondaryControlsCollection(playbackSource, popover))
}

// componentSecondaryControlsNoShare is like componentSecondaryControls but without the share button.
func componentSecondaryControlsNoShare(popover *gtk.PopoverMenu) schwifty.Box {
	return HStack(
		MenuButton().
			TooltipText(gettext.Get("More…")).
			Popover(popover).
			WithCSSClass("flat").
			WithCSSClass("circular").
			IconName("view-more-symbolic"),
	).Spacing(12).HAlign(gtk.AlignEndValue).HExpand(true).VAlign(gtk.AlignCenterValue)
}

func NewLocalPlaylist(playlistID, playlistName, coverUrl string, trackCount int, playFunc func(), shuffleFunc func()) schwifty.Widget {
	subtitle := gettext.GetN("%d Track", "%d Tracks", trackCount, trackCount)

	menu := gio.NewMenu()
	tracking.SetFinalizer("Menu", menu)

	queueAllItem := gio.NewMenuItem(gettext.Get("Add Playlist to Queue"), "win.player.queue")
	queueAllItem.SetActionAndTargetValue("win.player.queue", glib.NewVariantString("local_playlist/"+playlistID))
	menu.AppendItem(queueAllItem)
	tracking.SetFinalizer("MenuItem", queueAllItem)

	popover := gtk.NewPopoverMenuFromModel(&menu.MenuModel)
	tracking.SetFinalizer("Popover", popover)

	return template(coverUrl, playlistName, subtitle, "", componentControls(playFunc, shuffleFunc), componentSecondaryControlsNoShare(popover))
}
