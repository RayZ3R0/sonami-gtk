package adw

import (
	"codeberg.org/puregotk/puregotk/v4/adw"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/bindings/gtk"
)

//go:generate go run github.com/RayZ3R0/sonami-gtk/pkg/schwifty/gen WrapBox *adw.WrapBox adw

func (f WrapBox) Append(child any) WrapBox {
	return func() *adw.WrapBox {
		wrap := f()
		wrap.Append(gtk.ResolveWidget(child))
		return wrap
	}
}

func (f WrapBox) ChildSpacing(spacing int32) WrapBox {
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

func (f WrapBox) LineSpacing(spacing int32) WrapBox {
	return func() *adw.WrapBox {
		wrap := f()
		wrap.SetLineSpacing(spacing)
		return wrap
	}
}
