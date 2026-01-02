package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type WrapBox func() *adw.WrapBox

func (f WrapBox) AddController(controller *gtk.EventController) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f WrapBox) ConnectConstruct(cb func(*adw.WrapBox)) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f WrapBox) ConnectDestroy(cb func(gtk.Widget)) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f WrapBox) ConnectRealize(cb func(gtk.Widget)) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f WrapBox) ConnectUnrealize(cb func(gtk.Widget)) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f WrapBox) Focusable(focusable bool) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f WrapBox) FocusOnClick(focusOnClick bool) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f WrapBox) HAlign(align gtk.Align) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f WrapBox) HExpand(expand bool) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f WrapBox) HMargin(horizontal int) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f WrapBox) Margin(margin int) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f WrapBox) MarginBottom(bottom int) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f WrapBox) MarginEnd(end int) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f WrapBox) MarginStart(start int) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f WrapBox) MarginTop(top int) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f WrapBox) Opacity(opacity float64) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f WrapBox) Overflow(overflow gtk.Overflow) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f WrapBox) SizeRequest(width, height int) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f WrapBox) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f WrapBox) VAlign(align gtk.Align) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f WrapBox) VExpand(expand bool) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f WrapBox) Visible(visible bool) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f WrapBox) VMargin(vertical int) WrapBox {
	return func() *adw.WrapBox {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f WrapBox) Background(color string) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f WrapBox) CornerRadius(radius int) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f WrapBox) CSS(css string) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f WrapBox) BindCSSClass(state *state.State[string]) WrapBox {
	return func() *adw.WrapBox {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.WrapBox) {
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

func (f WrapBox) WithCSSClass(className string) WrapBox {
	return func() *adw.WrapBox {
		w := f()
		styleContext := w.GetStyleContext()
		defer styleContext.Unref()

		styleContext.AddClass(className)
		return w
	}
}

func (f WrapBox) CSSWithCallback(cb func(elementName string) string) WrapBox {
	return func() *adw.WrapBox {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *adw.WrapBox) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f WrapBox) HPadding(padding int) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f WrapBox) MinHeight(minHeight int) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f WrapBox) MinWidth(minWidth int) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f WrapBox) Padding(padding int) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f WrapBox) PaddingBottom(padding int) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f WrapBox) PaddingEnd(padding int) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f WrapBox) PaddingStart(padding int) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f WrapBox) PaddingTop(padding int) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f WrapBox) VPadding(padding int) WrapBox {
	return func() *adw.WrapBox {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f WrapBox) BindVisible(state *state.State[bool]) WrapBox {
	return func() *adw.WrapBox {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.WrapBox) {
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

func (f WrapBox) BindHMargin(state *state.State[int]) WrapBox {
	return func() *adw.WrapBox {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.WrapBox) {
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

func (f WrapBox) BindMargin(state *state.State[int]) WrapBox {
	return func() *adw.WrapBox {
		var callbackId string
		return f.ConnectConstruct(func(widget *adw.WrapBox) {
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

func (f WrapBox) BindMarginBottom(state *state.State[int]) WrapBox {
	return func() *adw.WrapBox {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.WrapBox) {
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

func (f WrapBox) BindMarginEnd(state *state.State[int]) WrapBox {
	return func() *adw.WrapBox {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.WrapBox) {
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

func (f WrapBox) BindMarginStart(state *state.State[int]) WrapBox {
	return func() *adw.WrapBox {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.WrapBox) {
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

func (f WrapBox) BindMarginTop(state *state.State[int]) WrapBox {
	return func() *adw.WrapBox {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.WrapBox) {
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

func (f WrapBox) BindSensitive(state *state.State[bool]) WrapBox {
	return func() *adw.WrapBox {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.WrapBox) {
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
