package tracklist_header

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type shareablePlaybackSource interface {
	tonearm.PlaybackSource
	tonearm.Shareable
}

func secondaryControlsCollection(playbackSource shareablePlaybackSource) schwifty.Box {
	return componentSecondaryControls(playbackSource)
}

func NewCollection(playbackSource shareablePlaybackSource, playFunc func(), shuffleFunc func()) schwifty.Box {
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

	return template(coverUrl, title, subtitle, "", componentControls(playFunc, shuffleFunc, popover), secondaryControlsCollection(playbackSource))
}
