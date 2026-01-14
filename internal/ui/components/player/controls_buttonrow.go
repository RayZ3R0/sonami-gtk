package player

import (
	"fmt"

	"codeberg.org/dergs/tonearm/internal/notifications"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func controlsButtonRow() schwifty.Box {
	return HStack(
		MenuButton().
			Popover(controlsVolumeSlider()).
			IconName("speakers-symbolic").
			WithCSSClass("transparent"),
		Button().
			ActionName("unimplemented").
			IconName("heart-outline-thick-symbolic").
			WithCSSClass("transparent"),
		Button().
			ActionName("unimplemented").
			IconName("compass2-symbolic").
			WithCSSClass("transparent"),
		Button().
			ActionName("unimplemented").
			IconName("library-symbolic").
			WithCSSClass("transparent"),
		Button().
			IconName("share-alt-symbolic").
			WithCSSClass("transparent").
			ConnectClicked(func(gtk.Button) {
				if trackID == "" {
					notifications.OnToast.Notify("No track is currently playing.")
					return
				}

				display := gdk.DisplayGetDefault()
				defer display.Unref()
				clipboard := display.GetClipboard()
				defer clipboard.Unref()

				clipboard.SetText(fmt.Sprintf("https://tidal.com/track/%s?u", trackID))
				notifications.OnToast.Notify("Copied track URL to clipboard.")
			}),
	).
		HAlign(gtk.AlignCenterValue).
		Spacing(7)
}
