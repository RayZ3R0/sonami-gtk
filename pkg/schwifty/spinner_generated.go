package schwifty

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type Spinner func() *gtk.Spinner

func (f Spinner) AddController(controller *gtk.EventController) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f Spinner) ConnectConstruct(cb func(*gtk.Spinner)) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f Spinner) ConnectDestroy(cb func(gtk.Widget)) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		callback.HandleCallback(widget.Object, "destroy", cb)
		return widget
	}
}

func (f Spinner) ConnectMap(cb func(gtk.Widget)) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		callback.HandleCallback(widget.Object, "map", cb)
		return widget
	}
}

func (f Spinner) ConnectRealize(cb func(gtk.Widget)) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		callback.HandleCallback(widget.Object, "realize", cb)
		return widget
	}
}

func (f Spinner) ConnectUnmap(cb func(gtk.Widget)) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		callback.HandleCallback(widget.Object, "unmap", cb)
		return widget
	}
}

func (f Spinner) ConnectUnrealize(cb func(gtk.Widget)) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		callback.HandleCallback(widget.Object, "unrealize", cb)
		return widget
	}
}

func (f Spinner) Focusable(focusable bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f Spinner) FocusOnClick(focusOnClick bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f Spinner) HAlign(align gtk.Align) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f Spinner) HExpand(expand bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f Spinner) HMargin(horizontal int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f Spinner) Margin(margin int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f Spinner) MarginBottom(bottom int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f Spinner) MarginEnd(end int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f Spinner) MarginStart(start int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f Spinner) MarginTop(top int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f Spinner) Opacity(opacity float64) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f Spinner) Overflow(overflow gtk.Overflow) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f Spinner) Sensitive(sensitive bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetSensitive(sensitive)
		return widget
	}
}

func (f Spinner) SizeRequest(width, height int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f Spinner) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f Spinner) VAlign(align gtk.Align) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f Spinner) VExpand(expand bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f Spinner) Visible(visible bool) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f Spinner) VMargin(vertical int) Spinner {
	return func() *gtk.Spinner {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f Spinner) Background(color string) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f Spinner) CornerRadius(radius int) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f Spinner) CSS(css string) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f Spinner) BindCSSClass(state *state.State[string]) Spinner {
	return func() *gtk.Spinner {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Spinner) {
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

func (f Spinner) WithCSSClass(className string) Spinner {
	return func() *gtk.Spinner {
		w := f()
		styleContext := w.GetStyleContext()
		defer styleContext.Unref()

		styleContext.AddClass(className)
		return w
	}
}

func (f Spinner) CSSWithCallback(cb func(elementName string) string) Spinner {
	return func() *gtk.Spinner {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *gtk.Spinner) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f Spinner) HPadding(padding int) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f Spinner) MinHeight(minHeight int) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f Spinner) MinWidth(minWidth int) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f Spinner) Padding(padding int) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f Spinner) PaddingBottom(padding int) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f Spinner) PaddingEnd(padding int) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f Spinner) PaddingStart(padding int) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f Spinner) PaddingTop(padding int) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f Spinner) VPadding(padding int) Spinner {
	return func() *gtk.Spinner {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f Spinner) BindVisible(state *state.State[bool]) Spinner {
	return func() *gtk.Spinner {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Spinner) {
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

func (f Spinner) BindHMargin(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Spinner) {
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

func (f Spinner) BindMargin(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
		var callbackId string
		return f.ConnectConstruct(func(widget *gtk.Spinner) {
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

func (f Spinner) BindMarginBottom(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Spinner) {
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

func (f Spinner) BindMarginEnd(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Spinner) {
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

func (f Spinner) BindMarginStart(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Spinner) {
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

func (f Spinner) BindMarginTop(state *state.State[int]) Spinner {
	return func() *gtk.Spinner {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Spinner) {
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

func (f Spinner) BindSensitive(state *state.State[bool]) Spinner {
	return func() *gtk.Spinner {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.Spinner) {
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
