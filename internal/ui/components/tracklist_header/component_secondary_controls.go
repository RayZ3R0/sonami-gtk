package tracklist_header

import (
	"codeberg.org/puregotk/puregotk/v4/gdk"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/notifications"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

func componentSecondaryControls(shareable sonami.Shareable, popover *gtk.PopoverMenu, buttons ...any) schwifty.Box {
	buttons = append(buttons, Button().
		TooltipText(gettext.Get("Copy URL")).
		IconName("share-alt-symbolic").
		WithCSSClass("flat").
		ConnectClicked(func(gtk.Button) {
			display := gdk.DisplayGetDefault()
			defer display.Unref()
			clipboard := display.GetClipboard()
			defer clipboard.Unref()

			clipboard.SetText(shareable.URL())
			notifications.OnToast.Notify(gettext.Get("Copied URL to clipboard"))
		}),
		MenuButton().
			TooltipText(gettext.Get("More…")).
			Popover(popover).
			WithCSSClass("flat").
			WithCSSClass("circular").
			IconName("view-more-symbolic"))
	return HStack(
		buttons...,
	).Spacing(12).HAlign(gtk.AlignEndValue).HExpand(true).VAlign(gtk.AlignCenterValue)
}
