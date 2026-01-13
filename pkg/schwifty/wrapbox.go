package schwifty

import "github.com/jwijenbergh/puregotk/v4/adw"

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen WrapBox *adw.WrapBox

func (f WrapBox) Append(child any) WrapBox {
	return func() *adw.WrapBox {
		wrap := f()
		wrap.Append(ResolveWidget(child))
		return wrap
	}
}

func (f WrapBox) ChildSpacing(spacing int) WrapBox {
	return func() *adw.WrapBox {
		wrap := f()
		wrap.SetChildSpacing(spacing)
		return wrap
	}
}

func (f WrapBox) Justify(justify adw.JustifyMode) WrapBox {
	return func() *adw.WrapBox {
		wrap := f()
		wrap.SetJustify(justify)
		return wrap
	}
}

func (f WrapBox) JustifyLastLine(justify bool) WrapBox {
	return func() *adw.WrapBox {
		wrap := f()
		wrap.SetJustifyLastLine(justify)
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
