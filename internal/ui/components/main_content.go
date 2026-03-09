package components

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

func MainContent(body schwifty.BaseWidgetable) schwifty.BaseWidgetable {
	return Clamp().Orientation(gtk.OrientationHorizontalValue).MaximumSize(2000).TighteningThreshold(2000).Child(body)
}
