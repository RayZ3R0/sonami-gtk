package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type Popover func() *gtk.Popover

func (f Popover) AddController(controller *gtk.EventController) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f Popover) ConnectConstruct(cb func(*gtk.Popover)) Popover {
	return func() *gtk.Popover {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f Popover) ConnectDestroy(cb func(gtk.Widget)) Popover {
	return func() *gtk.Popover {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f Popover) ConnectRealize(cb func(gtk.Widget)) Popover {
	return func() *gtk.Popover {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f Popover) ConnectUnrealize(cb func(gtk.Widget)) Popover {
	return func() *gtk.Popover {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f Popover) Focusable(focusable bool) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Popover) FocusOnClick(focusOnClick bool) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Popover) HAlign(align gtk.Align) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Popover) HExpand(expand bool) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Popover) HMargin(horizontal int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Popover) Margin(margin int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Popover) MarginBottom(bottom int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Popover) MarginEnd(end int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Popover) MarginStart(start int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Popover) MarginTop(top int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Popover) Opacity(opacity float64) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f Popover) Overflow(overflow gtk.Overflow) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f Popover) SizeRequest(width, height int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f Popover) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Popover) VAlign(align gtk.Align) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f Popover) VExpand(expand bool) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f Popover) Visible(visible bool) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f Popover) VMargin(vertical int) Popover {
	return func() *gtk.Popover {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f Popover) Background(color string) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f Popover) CornerRadius(radius int) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f Popover) CSS(css string) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f Popover) BindCSSClass(state *state.State[string]) Popover {
	return func() *gtk.Popover {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Popover) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue string) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).GetStyleContext().RemoveClass(state.Value())
				gtk.WidgetNewFromInternalPtr(widgetPtr).GetStyleContext().AddClass(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Popover) WithCSSClass(className string) Popover {
	return func() *gtk.Popover {
		w := f()
		w.GetStyleContext().AddClass(className)
		return w
	}
}

func (f Popover) CSSWithCallback(cb func(elementName string) string) Popover {
	return func() *gtk.Popover {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *gtk.Popover) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f Popover) HPadding(padding int) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f Popover) MinHeight(minHeight int) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f Popover) MinWidth(minWidth int) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f Popover) Padding(padding int) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f Popover) PaddingBottom(padding int) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f Popover) PaddingEnd(padding int) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f Popover) PaddingStart(padding int) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f Popover) PaddingTop(padding int) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f Popover) VPadding(padding int) Popover {
	return func() *gtk.Popover {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f Popover) BindVisible(state *state.State[bool]) Popover {
	return func() *gtk.Popover {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Popover) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetVisible(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Popover) BindHMargin(state *state.State[int]) Popover {
	return func() *gtk.Popover {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Popover) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginEnd(newValue)
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginStart(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Popover) BindMargin(state *state.State[int]) Popover {
	return func() *gtk.Popover {
		var callbackId string
		return f.ConnectConstruct(func(widget *gtk.Popover) {
			widgetPtr := widget.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginBottom(newValue)
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginEnd(newValue)
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginStart(newValue)
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginTop(newValue)
			})
		}).ConnectDestroy(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Popover) BindMarginBottom(state *state.State[int]) Popover {
	return func() *gtk.Popover {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Popover) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginBottom(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Popover) BindMarginEnd(state *state.State[int]) Popover {
	return func() *gtk.Popover {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Popover) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginEnd(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Popover) BindMarginStart(state *state.State[int]) Popover {
	return func() *gtk.Popover {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Popover) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginStart(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Popover) BindMarginTop(state *state.State[int]) Popover {
	return func() *gtk.Popover {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Popover) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginTop(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Popover) BindSensitive(state *state.State[bool]) Popover {
	return func() *gtk.Popover {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Popover) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetSensitive(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}
