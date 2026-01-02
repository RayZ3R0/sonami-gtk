package components

import (
	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/signals"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type RouteButton struct {
	*gtk.Button
	label *gtk.Label
	icon  *gtk.Image
}

func (r *RouteButton) SetActive(active bool) {
	if active {
		r.AddCssClass("active")
	} else {
		r.RemoveCssClass("active")
	}
}

func (r *RouteButton) SetTitle(title string) {
	r.label.SetText(title)
	r.label.SetVisible(title != "")
	if title == "" {
		r.RemoveCssClass("title")
	} else {
		r.AddCssClass("title")
	}
}

func (r *RouteButton) SetIcon(iconName string) {
	r.icon.SetFromIconName(iconName)
}

func NewRouteButton(path string) *RouteButton {
	routeButton := &RouteButton{
		icon:  Image().FromIconName("image-missing-symbolic")(),
		label: Label("").PaddingStart(7).PaddingEnd(7).Visible(false)(),
	}
	routeButton.Button = Button().
		PaddingStart(0).
		PaddingEnd(0).
		MinHeight(24).
		ConnectClicked(func(b gtk.Button) {
			router.Navigate(path)
		}).
		Child(
			HStack(
				Clamp().MaximumSize(16).Child(&routeButton.icon.Widget),
				routeButton.label,
			).
				Spacing(7).
				HMargin(9).
				VMargin(2),
		).
		WithCSSClass("transparent")()

	router.OnNavigate.On(func(newPath string) bool {
		routeButton.SetActive(path == newPath)
		return signals.Continue
	})

	return routeButton
}
