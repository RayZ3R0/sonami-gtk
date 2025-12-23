package components

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/ui/signals"
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type RouteButton struct {
	*gtk.Button
	box   *gtk.Box
	label *gtk.Label
	clamp *adw.Clamp
}

func (r *RouteButton) SetActive(active bool) {
	if active {
		r.RemoveCSSClass("image-button")
		r.box.SetMarginStart(7)
		r.box.SetMarginEnd(7)
	} else {
		r.SetCSSClasses([]string{"image-button"})
		r.box.SetMarginStart(12)
		r.box.SetMarginEnd(12)
	}
}

func (r *RouteButton) SetTitle(title string) {
	r.label.SetText(title)
	r.label.SetVisible(title != "")
}

func (r *RouteButton) SetIcon(iconName string) {
	r.clamp.SetChild(gtk.NewImageFromIconName(iconName))
}

func NewRouteButton(path string) *RouteButton {
	box := gtk.NewBox(gtk.OrientationHorizontal, 7)
	box.SetMarginBottom(2)
	box.SetMarginStart(12)
	box.SetMarginEnd(12)
	box.SetMarginTop(2)

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
	button.AddCSSClass("image-button")
	button.ConnectClicked(func() {
		router.NavigateTo(path, nil)
	})

	routeButton := &RouteButton{button, box, label, clamp}

	router.OnNavigate.On(func(newPath string) bool {
		routeButton.SetActive(path == newPath)
		return signals.Continue
	})

	return routeButton
}
