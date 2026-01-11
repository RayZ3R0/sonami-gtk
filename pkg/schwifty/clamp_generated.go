package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type Clamp func() *adw.Clamp

func (f Clamp) AddController(controller *gtk.EventController) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f Clamp) ConnectConstruct(cb func(*adw.Clamp)) Clamp {
	return func() *adw.Clamp {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f Clamp) ConnectDestroy(cb func(gtk.Widget)) Clamp {
	return func() *adw.Clamp {
		widget := f()
		callback.HandleCallback(widget.Object, "destroy", cb)
		return widget
	}
}

func (f Clamp) ConnectMap(cb func(gtk.Widget)) Clamp {
	return func() *adw.Clamp {
		widget := f()
		callback.HandleCallback(widget.Object, "map", cb)
		return widget
	}
}

func (f Clamp) ConnectRealize(cb func(gtk.Widget)) Clamp {
	return func() *adw.Clamp {
		widget := f()
		callback.HandleCallback(widget.Object, "realize", cb)
		return widget
	}
}

func (f Clamp) ConnectUnmap(cb func(gtk.Widget)) Clamp {
	return func() *adw.Clamp {
		widget := f()
		callback.HandleCallback(widget.Object, "unmap", cb)
		return widget
	}
}

func (f Clamp) ConnectUnrealize(cb func(gtk.Widget)) Clamp {
	return func() *adw.Clamp {
		widget := f()
		callback.HandleCallback(widget.Object, "unrealize", cb)
		return widget
	}
}

func (f Clamp) Focusable(focusable bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Clamp) FocusOnClick(focusOnClick bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Clamp) HAlign(align gtk.Align) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Clamp) HExpand(expand bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Clamp) HMargin(horizontal int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Clamp) Margin(margin int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Clamp) MarginBottom(bottom int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Clamp) MarginEnd(end int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Clamp) MarginStart(start int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Clamp) MarginTop(top int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Clamp) Opacity(opacity float64) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f Clamp) Overflow(overflow gtk.Overflow) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f Clamp) Sensitive(sensitive bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetSensitive(sensitive)
		return widget
	}
}

func (f Clamp) SizeRequest(width, height int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f Clamp) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Clamp) VAlign(align gtk.Align) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f Clamp) VExpand(expand bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f Clamp) Visible(visible bool) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f Clamp) VMargin(vertical int) Clamp {
	return func() *adw.Clamp {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f Clamp) Background(color string) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f Clamp) CornerRadius(radius int) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f Clamp) CSS(css string) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f Clamp) BindCSSClass(state *state.State[string]) Clamp {
	return func() *adw.Clamp {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.Clamp) {
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

func (f Clamp) WithCSSClass(className string) Clamp {
	return func() *adw.Clamp {
		w := f()
		styleContext := w.GetStyleContext()
		defer styleContext.Unref()

		styleContext.AddClass(className)
		return w
	}
}

func (f Clamp) CSSWithCallback(cb func(elementName string) string) Clamp {
	return func() *adw.Clamp {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *adw.Clamp) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f Clamp) HPadding(padding int) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f Clamp) MinHeight(minHeight int) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f Clamp) MinWidth(minWidth int) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f Clamp) Padding(padding int) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f Clamp) PaddingBottom(padding int) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f Clamp) PaddingEnd(padding int) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f Clamp) PaddingStart(padding int) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f Clamp) PaddingTop(padding int) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f Clamp) VPadding(padding int) Clamp {
	return func() *adw.Clamp {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f Clamp) BindVisible(state *state.State[bool]) Clamp {
	return func() *adw.Clamp {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.Clamp) {
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

func (f Clamp) BindHMargin(state *state.State[int]) Clamp {
	return func() *adw.Clamp {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.Clamp) {
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

func (f Clamp) BindMargin(state *state.State[int]) Clamp {
	return func() *adw.Clamp {
		var callbackId string
		return f.ConnectConstruct(func(widget *adw.Clamp) {
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

func (f Clamp) BindMarginBottom(state *state.State[int]) Clamp {
	return func() *adw.Clamp {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.Clamp) {
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

func (f Clamp) BindMarginEnd(state *state.State[int]) Clamp {
	return func() *adw.Clamp {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.Clamp) {
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

func (f Clamp) BindMarginStart(state *state.State[int]) Clamp {
	return func() *adw.Clamp {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.Clamp) {
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

func (f Clamp) BindMarginTop(state *state.State[int]) Clamp {
	return func() *adw.Clamp {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.Clamp) {
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

func (f Clamp) BindSensitive(state *state.State[bool]) Clamp {
	return func() *adw.Clamp {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.Clamp) {
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
