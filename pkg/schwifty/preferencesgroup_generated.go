package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type PreferencesGroup func() *adw.PreferencesGroup

func (f PreferencesGroup) AddController(controller *gtk.EventController) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f PreferencesGroup) ConnectConstruct(cb func(*adw.PreferencesGroup)) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f PreferencesGroup) ConnectDestroy(cb func(gtk.Widget)) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f PreferencesGroup) ConnectRealize(cb func(gtk.Widget)) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f PreferencesGroup) ConnectUnrealize(cb func(gtk.Widget)) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f PreferencesGroup) Focusable(focusable bool) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f PreferencesGroup) FocusOnClick(focusOnClick bool) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f PreferencesGroup) HAlign(align gtk.Align) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f PreferencesGroup) HExpand(expand bool) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f PreferencesGroup) HMargin(horizontal int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f PreferencesGroup) Margin(margin int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f PreferencesGroup) MarginBottom(bottom int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f PreferencesGroup) MarginEnd(end int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f PreferencesGroup) MarginStart(start int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f PreferencesGroup) MarginTop(top int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f PreferencesGroup) Opacity(opacity float64) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f PreferencesGroup) Overflow(overflow gtk.Overflow) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f PreferencesGroup) Sensitive(sensitive bool) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetSensitive(sensitive)
		return widget
	}
}

func (f PreferencesGroup) SizeRequest(width, height int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f PreferencesGroup) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f PreferencesGroup) VAlign(align gtk.Align) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f PreferencesGroup) VExpand(expand bool) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f PreferencesGroup) Visible(visible bool) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f PreferencesGroup) VMargin(vertical int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f PreferencesGroup) Background(color string) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f PreferencesGroup) CornerRadius(radius int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f PreferencesGroup) CSS(css string) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f PreferencesGroup) BindCSSClass(state *state.State[string]) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesGroup) {
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

func (f PreferencesGroup) WithCSSClass(className string) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		w := f()
		styleContext := w.GetStyleContext()
		defer styleContext.Unref()

		styleContext.AddClass(className)
		return w
	}
}

func (f PreferencesGroup) CSSWithCallback(cb func(elementName string) string) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *adw.PreferencesGroup) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f PreferencesGroup) HPadding(padding int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f PreferencesGroup) MinHeight(minHeight int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f PreferencesGroup) MinWidth(minWidth int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f PreferencesGroup) Padding(padding int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f PreferencesGroup) PaddingBottom(padding int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f PreferencesGroup) PaddingEnd(padding int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f PreferencesGroup) PaddingStart(padding int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f PreferencesGroup) PaddingTop(padding int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f PreferencesGroup) VPadding(padding int) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f PreferencesGroup) BindVisible(state *state.State[bool]) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesGroup) {
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

func (f PreferencesGroup) BindHMargin(state *state.State[int]) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesGroup) {
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

func (f PreferencesGroup) BindMargin(state *state.State[int]) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		var callbackId string
		return f.ConnectConstruct(func(widget *adw.PreferencesGroup) {
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

func (f PreferencesGroup) BindMarginBottom(state *state.State[int]) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesGroup) {
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

func (f PreferencesGroup) BindMarginEnd(state *state.State[int]) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesGroup) {
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

func (f PreferencesGroup) BindMarginStart(state *state.State[int]) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesGroup) {
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

func (f PreferencesGroup) BindMarginTop(state *state.State[int]) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesGroup) {
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

func (f PreferencesGroup) BindSensitive(state *state.State[bool]) PreferencesGroup {
	return func() *adw.PreferencesGroup {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesGroup) {
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
