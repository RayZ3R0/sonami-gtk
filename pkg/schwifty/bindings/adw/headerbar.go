package adw

import (
	gtkbindings "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/bindings/gtk"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/callback"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/state"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/utils/weak"
	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gtk"
)

//go:generate go run github.com/RayZ3R0/sonami-gtk/pkg/schwifty/gen HeaderBar *adw.HeaderBar adw

func (f HeaderBar) BindDecorationLayout(state *state.State[string]) HeaderBar {
	return func() *adw.HeaderBar {
		var callbackId string
		var ref weak.WidgetRef
		return f.ConnectConstruct(func(w *adw.HeaderBar) {
			ref = weak.NewWidgetRef(&w.Widget)
			callbackId = state.AddCallback(func(newValue string) {
				callback.OnMainThreadOncePure(func() {
					if obj := ref.Get(); obj != nil {
						defer obj.Unref()
						adw.HeaderBarNewFromInternalPtr(obj.Ptr).SetDecorationLayout(newValue)
					}
				})
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

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
			hb.PackEnd(gtkbindings.ResolveWidget(w))
		}
		return hb
	}
}

func (f HeaderBar) PackStart(widget ...any) HeaderBar {
	return func() *adw.HeaderBar {
		hb := f()
		for _, w := range widget {
			hb.PackStart(gtkbindings.ResolveWidget(w))
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
		hb.SetTitleWidget(gtkbindings.ResolveWidget(widget))
		return hb
	}
}
