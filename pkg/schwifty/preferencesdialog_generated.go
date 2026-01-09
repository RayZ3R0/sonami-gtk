package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type PreferencesDialog func() *adw.PreferencesDialog

func (f PreferencesDialog) AddController(controller *gtk.EventController) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f PreferencesDialog) ConnectConstruct(cb func(*adw.PreferencesDialog)) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f PreferencesDialog) ConnectDestroy(cb func(gtk.Widget)) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		callback.HandleCallback(widget.Object, "destroy", cb)
		return widget
	}
}

func (f PreferencesDialog) ConnectRealize(cb func(gtk.Widget)) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		callback.HandleCallback(widget.Object, "realize", cb)
		return widget
	}
}

func (f PreferencesDialog) ConnectUnrealize(cb func(gtk.Widget)) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		callback.HandleCallback(widget.Object, "unrealize", cb)
		return widget
	}
}

func (f PreferencesDialog) Focusable(focusable bool) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f PreferencesDialog) FocusOnClick(focusOnClick bool) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f PreferencesDialog) HAlign(align gtk.Align) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f PreferencesDialog) HExpand(expand bool) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f PreferencesDialog) HMargin(horizontal int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f PreferencesDialog) Margin(margin int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f PreferencesDialog) MarginBottom(bottom int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f PreferencesDialog) MarginEnd(end int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f PreferencesDialog) MarginStart(start int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f PreferencesDialog) MarginTop(top int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f PreferencesDialog) Opacity(opacity float64) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f PreferencesDialog) Overflow(overflow gtk.Overflow) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f PreferencesDialog) Sensitive(sensitive bool) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetSensitive(sensitive)
		return widget
	}
}

func (f PreferencesDialog) SizeRequest(width, height int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f PreferencesDialog) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f PreferencesDialog) VAlign(align gtk.Align) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f PreferencesDialog) VExpand(expand bool) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f PreferencesDialog) Visible(visible bool) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f PreferencesDialog) VMargin(vertical int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f PreferencesDialog) Background(color string) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f PreferencesDialog) CornerRadius(radius int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f PreferencesDialog) CSS(css string) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f PreferencesDialog) BindCSSClass(state *state.State[string]) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesDialog) {
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

func (f PreferencesDialog) WithCSSClass(className string) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		w := f()
		styleContext := w.GetStyleContext()
		defer styleContext.Unref()

		styleContext.AddClass(className)
		return w
	}
}

func (f PreferencesDialog) CSSWithCallback(cb func(elementName string) string) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *adw.PreferencesDialog) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f PreferencesDialog) HPadding(padding int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f PreferencesDialog) MinHeight(minHeight int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f PreferencesDialog) MinWidth(minWidth int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f PreferencesDialog) Padding(padding int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f PreferencesDialog) PaddingBottom(padding int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f PreferencesDialog) PaddingEnd(padding int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f PreferencesDialog) PaddingStart(padding int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f PreferencesDialog) PaddingTop(padding int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f PreferencesDialog) VPadding(padding int) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f PreferencesDialog) BindVisible(state *state.State[bool]) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesDialog) {
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

func (f PreferencesDialog) BindHMargin(state *state.State[int]) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesDialog) {
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

func (f PreferencesDialog) BindMargin(state *state.State[int]) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		var callbackId string
		return f.ConnectConstruct(func(widget *adw.PreferencesDialog) {
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

func (f PreferencesDialog) BindMarginBottom(state *state.State[int]) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesDialog) {
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

func (f PreferencesDialog) BindMarginEnd(state *state.State[int]) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesDialog) {
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

func (f PreferencesDialog) BindMarginStart(state *state.State[int]) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesDialog) {
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

func (f PreferencesDialog) BindMarginTop(state *state.State[int]) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesDialog) {
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

func (f PreferencesDialog) BindSensitive(state *state.State[bool]) PreferencesDialog {
	return func() *adw.PreferencesDialog {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.PreferencesDialog) {
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
