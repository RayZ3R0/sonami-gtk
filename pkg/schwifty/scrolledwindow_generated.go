package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type ScrolledWindow func() *gtk.ScrolledWindow

func (f ScrolledWindow) AddController(controller *gtk.EventController) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f ScrolledWindow) ConnectConstruct(cb func(*gtk.ScrolledWindow)) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f ScrolledWindow) ConnectDestroy(cb func(gtk.Widget)) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f ScrolledWindow) ConnectRealize(cb func(gtk.Widget)) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f ScrolledWindow) ConnectUnrealize(cb func(gtk.Widget)) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f ScrolledWindow) Focusable(focusable bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f ScrolledWindow) FocusOnClick(focusOnClick bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f ScrolledWindow) HAlign(align gtk.Align) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f ScrolledWindow) HExpand(expand bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f ScrolledWindow) HMargin(horizontal int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f ScrolledWindow) Margin(margin int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f ScrolledWindow) MarginBottom(bottom int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f ScrolledWindow) MarginEnd(end int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f ScrolledWindow) MarginStart(start int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f ScrolledWindow) MarginTop(top int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f ScrolledWindow) Opacity(opacity float64) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f ScrolledWindow) Overflow(overflow gtk.Overflow) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f ScrolledWindow) SizeRequest(width, height int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f ScrolledWindow) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f ScrolledWindow) VAlign(align gtk.Align) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f ScrolledWindow) VExpand(expand bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f ScrolledWindow) Visible(visible bool) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f ScrolledWindow) VMargin(vertical int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f ScrolledWindow) Background(color string) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f ScrolledWindow) CornerRadius(radius int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f ScrolledWindow) CSS(css string) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f ScrolledWindow) BindCSSClass(state *state.State[string]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.ScrolledWindow) {
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

func (f ScrolledWindow) WithCSSClass(className string) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		w := f()
		w.GetStyleContext().AddClass(className)
		return w
	}
}

func (f ScrolledWindow) CSSWithCallback(cb func(elementName string) string) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *gtk.ScrolledWindow) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f ScrolledWindow) HPadding(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f ScrolledWindow) MinHeight(minHeight int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f ScrolledWindow) MinWidth(minWidth int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f ScrolledWindow) Padding(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f ScrolledWindow) PaddingBottom(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f ScrolledWindow) PaddingEnd(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f ScrolledWindow) PaddingStart(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f ScrolledWindow) PaddingTop(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f ScrolledWindow) VPadding(padding int) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f ScrolledWindow) BindVisible(state *state.State[bool]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.ScrolledWindow) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetVisible(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f ScrolledWindow) BindHMargin(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.ScrolledWindow) {
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

func (f ScrolledWindow) BindMargin(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectConstruct(func(widget *gtk.ScrolledWindow) {
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

func (f ScrolledWindow) BindMarginBottom(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.ScrolledWindow) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginBottom(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f ScrolledWindow) BindMarginEnd(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.ScrolledWindow) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginEnd(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f ScrolledWindow) BindMarginStart(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.ScrolledWindow) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginStart(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f ScrolledWindow) BindMarginTop(state *state.State[int]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.ScrolledWindow) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginTop(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f ScrolledWindow) BindSensitive(state *state.State[bool]) ScrolledWindow {
	return func() *gtk.ScrolledWindow {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.ScrolledWindow) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetSensitive(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}
