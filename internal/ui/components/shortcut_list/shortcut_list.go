package shortcut_list

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/adw"
)

type ShortcutList struct {
	schwifty.Widget

	container *adw.WrapBox
}

func (f *ShortcutList) Append(widget schwifty.BaseWidgetable) {
	f.container.Append(widget.ToGTK())
}

func NewShortcutList() *ShortcutList {
	container := WrapBox().
		ChildSpacing(10).
		LineSpacing(10)()
	return &ShortcutList{
		Widget:    ManagedWidget(&container.Widget),
		container: container,
	}
}
