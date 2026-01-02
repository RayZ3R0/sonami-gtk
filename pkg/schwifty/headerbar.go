package schwifty

import "github.com/jwijenbergh/puregotk/v4/adw"

//go:generate go run codeberg.org/dergs/tidalwave/pkg/schwifty/gen HeaderBar *adw.HeaderBar

func (f HeaderBar) CenteringPolicy(policy adw.CenteringPolicy) HeaderBar {
	return func() *adw.HeaderBar {
		hb := f()
		hb.SetCenteringPolicy(policy)
		return hb
	}
}

func (f HeaderBar) DecorationLayout(layout string) HeaderBar {
	return func() *adw.HeaderBar {
		hb := f()
		hb.SetDecorationLayout(layout)
		return hb
	}
}

func (f HeaderBar) PackEnd(widget ...any) HeaderBar {
	return func() *adw.HeaderBar {
		hb := f()
		for _, w := range widget {
			hb.PackEnd(ResolveWidget(w))
		}
		return hb
	}
}

func (f HeaderBar) PackStart(widget ...any) HeaderBar {
	return func() *adw.HeaderBar {
		hb := f()
		for _, w := range widget {
			hb.PackStart(ResolveWidget(w))
		}
		return hb
	}
}

func (f HeaderBar) ShowBackButton(show bool) HeaderBar {
	return func() *adw.HeaderBar {
		hb := f()
		hb.SetShowBackButton(show)
		return hb
	}
}

func (f HeaderBar) ShowEndTitleButtons(show bool) HeaderBar {
	return func() *adw.HeaderBar {
		hb := f()
		hb.SetShowEndTitleButtons(show)
		return hb
	}
}

func (f HeaderBar) ShowStartTitleButtons(show bool) HeaderBar {
	return func() *adw.HeaderBar {
		hb := f()
		hb.SetShowStartTitleButtons(show)
		return hb
	}
}

func (f HeaderBar) TitleWidget(widget any) HeaderBar {
	return func() *adw.HeaderBar {
		hb := f()
		hb.SetTitleWidget(ResolveWidget(widget))
		return hb
	}
}
