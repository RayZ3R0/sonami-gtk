package components

import (
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
)

type RowTitle struct {
	schwifty.Box

	titleVisibility  *state.State[bool]
	titleText        *state.State[string]
	routeButtonState *state.State[any]
}

func (t *RowTitle) SetTitle(title string) *RowTitle {
	if title == "" {
		t.titleVisibility.SetValue(false)
		return t
	}
	t.titleText.SetValue(title)
	t.titleVisibility.SetValue(true)
	return t
}

func (t *RowTitle) SetViewAllRoute(path string) *RowTitle {
	t.routeButtonState.SetValue(Button().Child(
		Label(gettext.Get("View All")).WithCSSClass("caption-heading"),
	).
		MinHeight(10).
		MinWidth(10).
		HPadding(10).
		VAlign(gtk.AlignCenterValue).
		ConnectClicked(func(b gtk.Button) {
			router.Navigate(path)
		}))
	return t
}

func NewRowTitle() *RowTitle {
	t := &RowTitle{
		titleVisibility:  state.NewStateful(false),
		titleText:        state.NewStateful(""),
		routeButtonState: state.NewStateful[any](nil),
	}

	t.Box = HStack(
		Label("").
			Visible(false).
			BindText(t.titleText).
			BindVisible(t.titleVisibility).
			VAlign(gtk.AlignCenterValue).
			MarginStart(10).
			MarginBottom(5).
			WithCSSClass("title-2"),
		Spacer().VExpand(false),
		CenterBox().BindCenterWidget(t.routeButtonState).HExpand(false).VExpand(false),
	)

	return t
}
