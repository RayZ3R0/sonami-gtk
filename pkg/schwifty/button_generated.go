package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type Button func() *gtk.Button

func (f Button) AddController(controller *gtk.EventController) Button {
	return func() *gtk.Button {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f Button) ConnectConstruct(cb func(*gtk.Button)) Button {
	return func() *gtk.Button {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f Button) ConnectDestroy(cb func(gtk.Widget)) Button {
	return func() *gtk.Button {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f Button) ConnectRealize(cb func(gtk.Widget)) Button {
	return func() *gtk.Button {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f Button) ConnectUnrealize(cb func(gtk.Widget)) Button {
	return func() *gtk.Button {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f Button) Focusable(focusable bool) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Button) FocusOnClick(focusOnClick bool) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Button) HAlign(align gtk.Align) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Button) HExpand(expand bool) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Button) HMargin(horizontal int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Button) Margin(margin int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Button) MarginBottom(bottom int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Button) MarginEnd(end int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Button) MarginStart(start int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Button) MarginTop(top int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Button) Opacity(opacity float64) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f Button) Overflow(overflow gtk.Overflow) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f Button) SizeRequest(width, height int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f Button) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Button) VAlign(align gtk.Align) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f Button) VExpand(expand bool) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f Button) Visible(visible bool) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f Button) VMargin(vertical int) Button {
	return func() *gtk.Button {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f Button) Background(color string) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f Button) CornerRadius(radius int) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f Button) CSS(css string) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f Button) BindCSSClass(state *state.State[string]) Button {
	return func() *gtk.Button {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Button) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue string) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).GetStyleContext().RemoveClass(state.Value())
					gtk.WidgetNewFromInternalPtr(u).GetStyleContext().AddClass(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) WithCSSClass(className string) Button {
	return func() *gtk.Button {
		w := f()
		w.GetStyleContext().AddClass(className)
		return w
	}
}

func (f Button) CSSWithCallback(cb func(elementName string) string) Button {
	return func() *gtk.Button {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *gtk.Button) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f Button) HPadding(padding int) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f Button) MinHeight(minHeight int) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f Button) MinWidth(minWidth int) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f Button) Padding(padding int) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f Button) PaddingBottom(padding int) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f Button) PaddingEnd(padding int) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f Button) PaddingStart(padding int) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f Button) PaddingTop(padding int) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f Button) VPadding(padding int) Button {
	return func() *gtk.Button {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f Button) BindVisible(state *state.State[bool]) Button {
	return func() *gtk.Button {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Button) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetVisible(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) BindHMargin(state *state.State[int]) Button {
	return func() *gtk.Button {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Button) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginEnd(newValue)
					gtk.WidgetNewFromInternalPtr(u).SetMarginStart(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) BindMargin(state *state.State[int]) Button {
	return func() *gtk.Button {
		var callbackId string
		return f.ConnectConstruct(func(widget *gtk.Button) {
			widgetPtr := widget.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginBottom(newValue)
					gtk.WidgetNewFromInternalPtr(u).SetMarginEnd(newValue)
					gtk.WidgetNewFromInternalPtr(u).SetMarginStart(newValue)
					gtk.WidgetNewFromInternalPtr(u).SetMarginTop(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) BindMarginBottom(state *state.State[int]) Button {
	return func() *gtk.Button {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Button) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginBottom(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) BindMarginEnd(state *state.State[int]) Button {
	return func() *gtk.Button {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Button) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginEnd(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) BindMarginStart(state *state.State[int]) Button {
	return func() *gtk.Button {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Button) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginStart(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) BindMarginTop(state *state.State[int]) Button {
	return func() *gtk.Button {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Button) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetMarginTop(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f Button) BindSensitive(state *state.State[bool]) Button {
	return func() *gtk.Button {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Button) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).SetSensitive(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}
