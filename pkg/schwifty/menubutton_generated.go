package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type MenuButton func() *gtk.MenuButton

func (f MenuButton) AddController(controller *gtk.EventController) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f MenuButton) ConnectConstruct(cb func(*gtk.MenuButton)) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f MenuButton) ConnectDestroy(cb func(gtk.Widget)) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f MenuButton) ConnectRealize(cb func(gtk.Widget)) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f MenuButton) ConnectUnrealize(cb func(gtk.Widget)) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f MenuButton) Focusable(focusable bool) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f MenuButton) FocusOnClick(focusOnClick bool) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f MenuButton) HAlign(align gtk.Align) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f MenuButton) HExpand(expand bool) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f MenuButton) HMargin(horizontal int) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f MenuButton) Margin(margin int) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f MenuButton) MarginBottom(bottom int) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f MenuButton) MarginEnd(end int) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f MenuButton) MarginStart(start int) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f MenuButton) MarginTop(top int) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f MenuButton) Opacity(opacity float64) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f MenuButton) Overflow(overflow gtk.Overflow) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f MenuButton) Sensitive(sensitive bool) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetSensitive(sensitive)
		return widget
	}
}

func (f MenuButton) SizeRequest(width, height int) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f MenuButton) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f MenuButton) VAlign(align gtk.Align) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f MenuButton) VExpand(expand bool) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f MenuButton) Visible(visible bool) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f MenuButton) VMargin(vertical int) MenuButton {
	return func() *gtk.MenuButton {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f MenuButton) Background(color string) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f MenuButton) CornerRadius(radius int) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f MenuButton) CSS(css string) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f MenuButton) BindCSSClass(state *state.State[string]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.MenuButton) {
			ptr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue string) {
				oldValue := state.Value()
				callback.OnMainThreadOnce(func(u uintptr) {
					w := gtk.ButtonNewFromInternalPtr(u)
					styleContext := w.GetStyleContext()
					defer styleContext.Unref()

					styleContext.RemoveClass(oldValue)
					styleContext.AddClass(newValue)
				}, ptr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f MenuButton) WithCSSClass(className string) MenuButton {
	return func() *gtk.MenuButton {
		w := f()
		styleContext := w.GetStyleContext()
		defer styleContext.Unref()

		styleContext.AddClass(className)
		return w
	}
}

func (f MenuButton) CSSWithCallback(cb func(elementName string) string) MenuButton {
	return func() *gtk.MenuButton {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *gtk.MenuButton) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f MenuButton) HPadding(padding int) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f MenuButton) MinHeight(minHeight int) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f MenuButton) MinWidth(minWidth int) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f MenuButton) Padding(padding int) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f MenuButton) PaddingBottom(padding int) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f MenuButton) PaddingEnd(padding int) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f MenuButton) PaddingStart(padding int) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f MenuButton) PaddingTop(padding int) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f MenuButton) VPadding(padding int) MenuButton {
	return func() *gtk.MenuButton {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f MenuButton) BindVisible(state *state.State[bool]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.MenuButton) {
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

func (f MenuButton) BindHMargin(state *state.State[int]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.MenuButton) {
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

func (f MenuButton) BindMargin(state *state.State[int]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		return f.ConnectConstruct(func(widget *gtk.MenuButton) {
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

func (f MenuButton) BindMarginBottom(state *state.State[int]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.MenuButton) {
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

func (f MenuButton) BindMarginEnd(state *state.State[int]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.MenuButton) {
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

func (f MenuButton) BindMarginStart(state *state.State[int]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.MenuButton) {
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

func (f MenuButton) BindMarginTop(state *state.State[int]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.MenuButton) {
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

func (f MenuButton) BindSensitive(state *state.State[bool]) MenuButton {
	return func() *gtk.MenuButton {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.MenuButton) {
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
