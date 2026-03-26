package components

import (
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
)

func AuthRequired(description string) schwifty.StatusPage {
	return StatusPage().
		IconName("avatar-default-symbolic").
		Title(gettext.Get("Authentication Required")).
		Description(description).
		Child(
			Button().
				Label(gettext.Get("Sign In…")).
				WithCSSClass("pill").
				WithCSSClass("suggested-action").
				ActionName("win.sign-in").
				HAlign(gtk.AlignCenterValue),
		)
}
