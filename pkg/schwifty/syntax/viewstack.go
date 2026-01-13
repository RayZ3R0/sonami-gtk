package syntax

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

func ViewStack(children ...any) schwifty.ViewStack {
	return managed("ViewStack", func() *adw.ViewStack {
		viewStack := adw.NewViewStack()
		for _, child := range children {
			viewStack.Add(schwifty.ResolveWidget(child))
		}
		return viewStack
	})
}
