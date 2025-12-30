package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type SearchEntry func() *gtk.SearchEntry

func (f SearchEntry) AddController(controller *gtk.EventController) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f SearchEntry) ConnectConstruct(cb func(*gtk.SearchEntry)) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f SearchEntry) ConnectDestroy(cb func(gtk.Widget)) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f SearchEntry) ConnectRealize(cb func(gtk.Widget)) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f SearchEntry) ConnectUnrealize(cb func(gtk.Widget)) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f SearchEntry) Focusable(focusable bool) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f SearchEntry) FocusOnClick(focusOnClick bool) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f SearchEntry) HAlign(align gtk.Align) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f SearchEntry) HExpand(expand bool) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f SearchEntry) HMargin(horizontal int) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f SearchEntry) Margin(margin int) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f SearchEntry) MarginBottom(bottom int) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f SearchEntry) MarginEnd(end int) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f SearchEntry) MarginStart(start int) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f SearchEntry) MarginTop(top int) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f SearchEntry) Opacity(opacity float64) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f SearchEntry) Overflow(overflow gtk.Overflow) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f SearchEntry) SizeRequest(width, height int) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f SearchEntry) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f SearchEntry) VAlign(align gtk.Align) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f SearchEntry) VExpand(expand bool) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f SearchEntry) Visible(visible bool) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f SearchEntry) VMargin(vertical int) SearchEntry {
	return func() *gtk.SearchEntry {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f SearchEntry) Background(color string) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f SearchEntry) CornerRadius(radius int) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f SearchEntry) CSS(css string) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f SearchEntry) BindCSSClass(state *state.State[string]) SearchEntry {
	return func() *gtk.SearchEntry {
		var callbackId string
		return f.ConnectRealize(func(w gtk.Widget) {
			callbackId = state.AddCallback(func(newValue string) {
				w.GetStyleContext().RemoveClass(state.Value())
				w.GetStyleContext().AddClass(newValue)
			})
		}).ConnectUnrealize(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f SearchEntry) WithCSSClass(className string) SearchEntry {
	return func() *gtk.SearchEntry {
		w := f()
		w.GetStyleContext().AddClass(className)
		return w
	}
}

func (f SearchEntry) CSSWithCallback(cb func(elementName string) string) SearchEntry {
	return func() *gtk.SearchEntry {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *gtk.SearchEntry) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f SearchEntry) HPadding(padding int) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f SearchEntry) MinHeight(minHeight int) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f SearchEntry) MinWidth(minWidth int) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f SearchEntry) Padding(padding int) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f SearchEntry) PaddingBottom(padding int) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f SearchEntry) PaddingEnd(padding int) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f SearchEntry) PaddingStart(padding int) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f SearchEntry) PaddingTop(padding int) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f SearchEntry) VPadding(padding int) SearchEntry {
	return func() *gtk.SearchEntry {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f SearchEntry) BindVisible(state *state.State[bool]) SearchEntry {
	return func() *gtk.SearchEntry {
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

func (f SearchEntry) BindHMargin(state *state.State[int]) SearchEntry {
	return func() *gtk.SearchEntry {
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

func (f SearchEntry) BindMargin(state *state.State[int]) SearchEntry {
	return func() *gtk.SearchEntry {
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

func (f SearchEntry) BindMarginBottom(state *state.State[int]) SearchEntry {
	return func() *gtk.SearchEntry {
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

func (f SearchEntry) BindMarginEnd(state *state.State[int]) SearchEntry {
	return func() *gtk.SearchEntry {
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

func (f SearchEntry) BindMarginStart(state *state.State[int]) SearchEntry {
	return func() *gtk.SearchEntry {
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

func (f SearchEntry) BindMarginTop(state *state.State[int]) SearchEntry {
	return func() *gtk.SearchEntry {
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

func (f SearchEntry) BindSensitive(state *state.State[bool]) SearchEntry {
	return func() *gtk.SearchEntry {
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
