package components

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func MainContent(body schwifty.BaseWidgetable) schwifty.BaseWidgetable {
	return Clamp().Orientation(gtk.OrientationHorizontalValue).MaximumSize(2000).TighteningThreshold(2000).Child(body)
}
