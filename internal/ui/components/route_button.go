package components

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/signals"
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
)

var routeButtonCSS = cssutil.Applier("route-button", `
	.route-button {
		padding-left: 0px;
		padding-right: 0px;
		min-height: 24px;
	}

	.route-button.title {
		padding-left: 7px;
		padding-right: 7px;
	}

	.route-button:not(:hover):not(.active) {
		background-color: transparent;
	}
`)

type RouteButton struct {
	*gtk.Button
	box   *gtk.Box
	label *gtk.Label
	clamp *adw.Clamp
}

func (r *RouteButton) SetActive(active bool) {
	if active {
		r.AddCSSClass("active")
	} else {
		r.RemoveCSSClass("active")
	}
}

func (r *RouteButton) SetTitle(title string) {
	r.label.SetText(title)
	r.label.SetVisible(title != "")
	if title == "" {
		r.RemoveCSSClass("title")
	} else {
		r.AddCSSClass("title")
	}
}

func (r *RouteButton) SetIcon(iconName string) {
	r.clamp.SetChild(gtk.NewImageFromIconName(iconName))
}

func NewRouteButton(path string) *RouteButton {
	box := gtk.NewBox(gtk.OrientationHorizontal, 7)
	box.SetMarginBottom(2)
	box.SetMarginTop(2)
	box.SetMarginStart(9)
	box.SetMarginEnd(9)

	image := gtk.NewImageFromIconName("sidebar-show-symbolic")

	clamp := adw.NewClamp()
	clamp.SetMaximumSize(16)
	clamp.SetChild(image)

	label := gtk.NewLabel("")
	label.SetVisible(false)

	box.Append(clamp)
	box.Append(label)

	button := gtk.NewButton()
	button.SetChild(box)
	button.ConnectClicked(func() {
		router.Navigate(path, nil)
	})
	routeButtonCSS(button)

	routeButton := &RouteButton{button, box, label, clamp}

	router.OnNavigate.On(func(newPath string) bool {
		routeButton.SetActive(path == newPath)
		return signals.Continue
	})

	return routeButton
}
