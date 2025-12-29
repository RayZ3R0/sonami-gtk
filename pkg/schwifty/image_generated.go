package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/css"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type Image func() *gtk.Image

func (f Image) AddController(controller *gtk.EventController) Image {
	return func() *gtk.Image {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f Image) ConnectConstruct(cb func(*gtk.Image)) Image {
	return func() *gtk.Image {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f Image) ConnectDestroy(cb func(gtk.Widget)) Image {
	return func() *gtk.Image {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f Image) ConnectRealize(cb func(gtk.Widget)) Image {
	return func() *gtk.Image {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f Image) ConnectUnrealize(cb func(gtk.Widget)) Image {
	return func() *gtk.Image {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f Image) Focusable(focusable bool) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Image) FocusOnClick(focusOnClick bool) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Image) HAlign(align gtk.Align) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Image) HExpand(expand bool) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Image) HMargin(horizontal int) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Image) Margin(margin int) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Image) MarginBottom(bottom int) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Image) MarginEnd(end int) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Image) MarginStart(start int) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Image) MarginTop(top int) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Image) Opacity(opacity float64) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f Image) Overflow(overflow gtk.Overflow) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f Image) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Image) VAlign(align gtk.Align) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f Image) VExpand(expand bool) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f Image) Visible(visible bool) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f Image) VMargin(vertical int) Image {
	return func() *gtk.Image {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f Image) Background(color string) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { background-color: %s; }", widget.GetCssName(), color))
		return widget
	}
}

func (f Image) CornerRadius(radius int) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { border-radius: %dpx; }", widget.GetCssName(), radius))
		return widget
	}
}

func (f Image) CSS(css string) Image {
	return func() *gtk.Image {
		widget := f()
		widget.Ref()
		defer widget.Unref()

		provider := gtk.NewCssProvider()
		provider.LoadFromString(css)
		widget.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		provider.Unref()

		return widget
	}
}

func (f Image) HPadding(padding int) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}

func (f Image) MinHeight(minHeight int) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-height: %dpx; }", widget.GetCssName(), minHeight))
		return widget
	}
}

func (f Image) MinWidth(minWidth int) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { min-width: %dpx; }", widget.GetCssName(), minWidth))
		return widget
	}
}

func (f Image) Padding(padding int) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Image) PaddingBottom(padding int) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Image) PaddingEnd(padding int) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-right: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Image) PaddingStart(padding int) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-left: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Image) PaddingTop(padding int) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-top: %dpx; }", widget.GetCssName(), padding))
		return widget
	}
}

func (f Image) VPadding(padding int) Image {
	return func() *gtk.Image {
		widget := f()
		css.Apply(&widget.Widget, fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", widget.GetCssName(), padding, padding))
		return widget
	}
}



func (f Image) BindVisible(state *state.State[bool]) Image {
	return func() *gtk.Image {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue bool) {
				w.SetVisible(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Image) BindHMargin(state *state.State[int]) Image {
	return func() *gtk.Image {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				w.SetMarginEnd(newValue)
				w.SetMarginStart(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Image) BindMargin(state *state.State[int]) Image {
	return func() *gtk.Image {
		var callbackId string
		return f.ConnectRealize(func(widget gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				widget.SetMarginBottom(newValue)
				widget.SetMarginEnd(newValue)
				widget.SetMarginStart(newValue)
				widget.SetMarginTop(newValue)
			})
		}).ConnectUnrealize(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Image) BindMarginBottom(state *state.State[int]) Image {
	return func() *gtk.Image {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				w.SetMarginBottom(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Image) BindMarginEnd(state *state.State[int]) Image {
	return func() *gtk.Image {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				w.SetMarginEnd(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Image) BindMarginStart(state *state.State[int]) Image {
	return func() *gtk.Image {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				w.SetMarginStart(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Image) BindMarginTop(state *state.State[int]) Image {
	return func() *gtk.Image {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue int) {
				w.SetMarginTop(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Image) BindSensitive(state *state.State[bool]) Image {
	return func() *gtk.Image {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue bool) {
				w.SetSensitive(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}
