package schwifty

import "github.com/jwijenbergh/puregotk/v4/adw"

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen WindowTitle *adw.WindowTitle

func (f WindowTitle) Title(title string) WindowTitle {
	return func() *adw.WindowTitle {
		wt := f()
		wt.SetTitle(title)
		return wt
	}
}

func (f WindowTitle) SubTitle(subtitle string) WindowTitle {
	return func() *adw.WindowTitle {
		wt := f()
		wt.SetSubtitle(subtitle)
		return wt
	}
}
