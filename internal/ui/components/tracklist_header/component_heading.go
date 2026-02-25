package tracklist_header

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gtk"
	"github.com/jwijenbergh/puregotk/v4/pango"
)

func componentHeading(title string, subtitle string) schwifty.Box {
	return VStack(
		Label(title).WithCSSClass("title-2").Ellipsis(pango.EllipsizeEndValue).HAlign(gtk.AlignStartValue),
		Label(subtitle).Ellipsis(pango.EllipsizeEndValue).WithCSSClass("heading").WithCSSClass("dimmed").HAlign(gtk.AlignStartValue),
	).HAlign(gtk.AlignStartValue)
}
