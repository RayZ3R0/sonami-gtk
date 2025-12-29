package schwifty

import "github.com/jwijenbergh/puregotk/v4/adw"

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen WrapBox *adw.WrapBox

func (f WrapBox) ChildSpacing(spacing int) WrapBox {
	return func() *adw.WrapBox {
		wrap := f()
		wrap.SetChildSpacing(spacing)
		return wrap
	}
}

func (f WrapBox) LineSpacing(spacing int) WrapBox {
	return func() *adw.WrapBox {
		wrap := f()
		wrap.SetLineSpacing(spacing)
		return wrap
	}
}
