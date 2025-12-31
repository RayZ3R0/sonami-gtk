package schwifty

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty/callback"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	"fmt"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)


type AspectFrame func() *gtk.AspectFrame

func (f AspectFrame) AddController(controller *gtk.EventController) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.AddController(controller)
		return widget
	}
}

func (f AspectFrame) ConnectConstruct(cb func(*gtk.AspectFrame)) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		cb(widget)
		return widget
	}
}

func (f AspectFrame) ConnectDestroy(cb func(gtk.Widget)) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		callback.HandleCallback(widget.Widget, "destroy", cb)
		return widget
	}
}

func (f AspectFrame) ConnectRealize(cb func(gtk.Widget)) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		callback.HandleCallback(widget.Widget, "realize", cb)
		return widget
	}
}

func (f AspectFrame) ConnectUnrealize(cb func(gtk.Widget)) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		callback.HandleCallback(widget.Widget, "unrealize", cb)
		return widget
	}
}

func (f AspectFrame) Focusable(focusable bool) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetFocusable(focusable)
		return widget
	}
}

func (f AspectFrame) FocusOnClick(focusOnClick bool) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetFocusOnClick(focusOnClick)
		return widget
	}
}

func (f AspectFrame) HAlign(align gtk.Align) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetHalign(align)
		return widget
	}
}

func (f AspectFrame) HExpand(expand bool) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetHexpand(expand)
		return widget
	}
}

func (f AspectFrame) HMargin(horizontal int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginEnd(horizontal)
		widget.SetMarginStart(horizontal)
		return widget
	}
}

func (f AspectFrame) Margin(margin int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginBottom(margin)
		widget.SetMarginEnd(margin)
		widget.SetMarginStart(margin)
		widget.SetMarginTop(margin)
		return widget
	}
}

func (f AspectFrame) MarginBottom(bottom int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginBottom(bottom)
		return widget
	}
}

func (f AspectFrame) MarginEnd(end int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginEnd(end)
		return widget
	}
}

func (f AspectFrame) MarginStart(start int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginStart(start)
		return widget
	}
}

func (f AspectFrame) MarginTop(top int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginTop(top)
		return widget
	}
}

func (f AspectFrame) Opacity(opacity float64) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetOpacity(opacity)
		return widget
	}
}

func (f AspectFrame) Overflow(overflow gtk.Overflow) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetOverflow(overflow)
		return widget
	}
}

func (f AspectFrame) SizeRequest(width, height int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetSizeRequest(width, height)
		return widget
	}
}

func (f AspectFrame) ToGTK() *gtk.Widget {
	val := f()
	return &val.Widget
}

func (f AspectFrame) VAlign(align gtk.Align) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetValign(align)
		return widget
	}
}

func (f AspectFrame) VExpand(expand bool) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetVexpand(expand)
		return widget
	}
}

func (f AspectFrame) Visible(visible bool) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetVisible(visible)
		return widget
	}
}

func (f AspectFrame) VMargin(vertical int) AspectFrame {
	return func() *gtk.AspectFrame {
		widget := f()
		widget.SetMarginTop(vertical)
		widget.SetMarginBottom(vertical)
		return widget
	}
}



func (f AspectFrame) Background(color string) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { background-color: %s; }", elementName, color)
		})()
	}
}

func (f AspectFrame) CornerRadius(radius int) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { border-radius: %dpx; }", elementName, radius)
		})()
	}
}

func (f AspectFrame) CSS(css string) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return css
		})()
	}
}

func (f AspectFrame) BindCSSClass(state *state.State[string]) AspectFrame {
	return func() *gtk.AspectFrame {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.AspectFrame) {
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

func (f AspectFrame) WithCSSClass(className string) AspectFrame {
	return func() *gtk.AspectFrame {
		w := f()
		w.GetStyleContext().AddClass(className)
		return w
	}
}

func (f AspectFrame) CSSWithCallback(cb func(elementName string) string) AspectFrame {
	return func() *gtk.AspectFrame {
		provider := gtk.NewCssProvider()
		return f.ConnectConstruct(func(t *gtk.AspectFrame) {
			provider.LoadFromString(cb(t.GetCssName()))
			t.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
		}).ConnectDestroy(func(w gtk.Widget) {
			w.GetStyleContext().RemoveProvider(provider)
			provider.Unref()
			provider = nil
		})()
	}
}

func (f AspectFrame) HPadding(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; padding-right: %dpx; }", elementName, padding, padding)
		})()
	}
}

func (f AspectFrame) MinHeight(minHeight int) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-height: %dpx; }", elementName, minHeight)
		})()
	}
}

func (f AspectFrame) MinWidth(minWidth int) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { min-width: %dpx; }", elementName, minWidth)
		})()
	}
}

func (f AspectFrame) Padding(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding: %dpx; }", elementName, padding)
		})()
	}
}

func (f AspectFrame) PaddingBottom(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; }", elementName, padding)
		})()
	}
}

func (f AspectFrame) PaddingEnd(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-right: %dpx; }", elementName, padding)
		})()
	}
}

func (f AspectFrame) PaddingStart(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-left: %dpx; }", elementName, padding)
		})()
	}
}

func (f AspectFrame) PaddingTop(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-top: %dpx; }", elementName, padding)
		})()
	}
}

func (f AspectFrame) VPadding(padding int) AspectFrame {
	return func() *gtk.AspectFrame {
		return f.CSSWithCallback(func(elementName string) string {
			return fmt.Sprintf("%s { padding-bottom: %dpx; padding-top: %dpx; }", elementName, padding, padding)
		})()
	}
}



func (f AspectFrame) BindVisible(state *state.State[bool]) AspectFrame {
	return func() *gtk.AspectFrame {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.AspectFrame) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetVisible(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f AspectFrame) BindHMargin(state *state.State[int]) AspectFrame {
	return func() *gtk.AspectFrame {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.AspectFrame) {
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

func (f AspectFrame) BindMargin(state *state.State[int]) AspectFrame {
	return func() *gtk.AspectFrame {
		var callbackId string
		return f.ConnectConstruct(func(widget *gtk.AspectFrame) {
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

func (f AspectFrame) BindMarginBottom(state *state.State[int]) AspectFrame {
	return func() *gtk.AspectFrame {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.AspectFrame) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginBottom(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f AspectFrame) BindMarginEnd(state *state.State[int]) AspectFrame {
	return func() *gtk.AspectFrame {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.AspectFrame) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginEnd(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f AspectFrame) BindMarginStart(state *state.State[int]) AspectFrame {
	return func() *gtk.AspectFrame {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.AspectFrame) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginStart(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f AspectFrame) BindMarginTop(state *state.State[int]) AspectFrame {
	return func() *gtk.AspectFrame {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.AspectFrame) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue int) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetMarginTop(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}

func (f AspectFrame) BindSensitive(state *state.State[bool]) AspectFrame {
	return func() *gtk.AspectFrame {
		var callbackId string
		return f.ConnectConstruct(func(w *gtk.AspectFrame) {
			widgetPtr := w.GoPointer()
			callbackId = state.AddCallback(func(newValue bool) {
				gtk.WidgetNewFromInternalPtr(widgetPtr).SetSensitive(newValue)
			})
		}).ConnectDestroy(func(w gtk.Widget) {
			state.RemoveCallback(callbackId)
		})()
	}
}
