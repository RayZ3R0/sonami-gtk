package shortcut_list

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
	"codeberg.org/puregotk/puregotk/v4/adw"
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
		Widget:    Widget(&container.Widget),
		container: container,
	}
}
