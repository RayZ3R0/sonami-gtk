package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type StatusPage func() *adw.StatusPage

func (f StatusPage) AddController(controller *gtk.EventController) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f StatusPage) ConnectConstruct(cb func(*adw.StatusPage)) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f StatusPage) ConnectDestroy(cb func(gtk.Widget)) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f StatusPage) ConnectRealize(cb func(gtk.Widget)) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f StatusPage) ConnectUnrealize(cb func(gtk.Widget)) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f StatusPage) Focusable(focusable bool) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f StatusPage) FocusOnClick(focusOnClick bool) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f StatusPage) HAlign(align gtk.Align) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f StatusPage) HExpand(expand bool) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f StatusPage) HMargin(horizontal int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f StatusPage) Margin(margin int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f StatusPage) MarginBottom(bottom int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f StatusPage) MarginEnd(end int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f StatusPage) MarginStart(start int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f StatusPage) MarginTop(top int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f StatusPage) Opacity(opacity float64) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f StatusPage) Overflow(overflow gtk.Overflow) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f StatusPage) SizeRequest(width, height int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f StatusPage) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f StatusPage) VAlign(align gtk.Align) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f StatusPage) VExpand(expand bool) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f StatusPage) Visible(visible bool) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f StatusPage) VMargin(vertical int) StatusPage {
	return func() *adw.StatusPage {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f StatusPage) Background(color string) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f StatusPage) CornerRadius(radius int) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f StatusPage) CSS(css string) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f StatusPage) BindCSSClass(state *state.State[string]) StatusPage {
	return func() *adw.StatusPage {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.StatusPage) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue string) {
				oldState := state.Value()
				callback.OnMainThreadOnce(func(u uintptr) {
					gtk.WidgetNewFromInternalPtr(u).GetStyleContext().RemoveClass(oldState)
					gtk.WidgetNewFromInternalPtr(u).GetStyleContext().AddClass(newValue)
				}, widgetPtr)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f StatusPage) WithCSSClass(className string) StatusPage {
	return func() *adw.StatusPage {
		w := f()
		w.GetStyleContext().AddClass(className)
		return w
	}
}

func (f StatusPage) CSSWithCallback(cb func(elementName string) string) StatusPage {
	return func() *adw.StatusPage {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *adw.StatusPage) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f StatusPage) HPadding(padding int) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f StatusPage) MinHeight(minHeight int) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f StatusPage) MinWidth(minWidth int) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f StatusPage) Padding(padding int) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f StatusPage) PaddingBottom(padding int) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f StatusPage) PaddingEnd(padding int) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f StatusPage) PaddingStart(padding int) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f StatusPage) PaddingTop(padding int) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f StatusPage) VPadding(padding int) StatusPage {
	return func() *adw.StatusPage {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f StatusPage) BindVisible(state *state.State[bool]) StatusPage {
	return func() *adw.StatusPage {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.StatusPage) {
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

func (f StatusPage) BindHMargin(state *state.State[int]) StatusPage {
	return func() *adw.StatusPage {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.StatusPage) {
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

func (f StatusPage) BindMargin(state *state.State[int]) StatusPage {
	return func() *adw.StatusPage {
		var callbackId string
		return f.ConnectConstruct(func(widget *adw.StatusPage) {
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

func (f StatusPage) BindMarginBottom(state *state.State[int]) StatusPage {
	return func() *adw.StatusPage {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.StatusPage) {
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

func (f StatusPage) BindMarginEnd(state *state.State[int]) StatusPage {
	return func() *adw.StatusPage {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.StatusPage) {
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

func (f StatusPage) BindMarginStart(state *state.State[int]) StatusPage {
	return func() *adw.StatusPage {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.StatusPage) {
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

func (f StatusPage) BindMarginTop(state *state.State[int]) StatusPage {
	return func() *adw.StatusPage {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.StatusPage) {
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

func (f StatusPage) BindSensitive(state *state.State[bool]) StatusPage {
	return func() *adw.StatusPage {
		var callbackId string
		return f.ConnectConstruct(func(w *adw.StatusPage) {
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
