package syntax

import (
	gtkbindings "codeberg.org/dergs/tonearm/pkg/schwifty/bindings/gtk"
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func AspectFrame(child any) gtkbindings.AspectFrame {
	return managed("Scale", func() *gtk.AspectFrame {
		aspectFrame := gtk.NewAspectFrame(0.5, 0.5, 1.0, false)
		aspectFrame.SetChild(gtkbindings.ResolveWidget(child))
		return aspectFrame
	})
}

func Box(orientation gtk.Orientation, children ...any) gtkbindings.Box {
	return managed("Box", func() *gtk.Box {
		box := gtk.NewBox(orientation, 0)
		for _, child := range children {
			box.Append(gtkbindings.ResolveWidget(child))
		}
		return box
	})
}

func HStack(children ...any) gtkbindings.Box {
	return Box(gtk.OrientationHorizontalValue, children...)
}

func Spacer() gtkbindings.Box {
	return HStack().VExpand(true).HExpand(true)
}

func VStack(children ...any) gtkbindings.Box {
	return Box(gtk.OrientationVerticalValue, children...)
}

func Button() gtkbindings.Button {
	return managed("Button", func() *gtk.Button {
		btn := gtk.NewButton()
		btn.ConnectClicked(&callback.ButtonClickedCallback)
		return btn
	})
}

func CenterBox() gtkbindings.CenterBox {
	return managed("CenterBox", func() *gtk.CenterBox {
		return gtk.NewCenterBox()
	})
}

func Image() gtkbindings.Image {
	return managed("Image", func() *gtk.Image {
		return gtk.NewImage()
	})
}

func Label(text string) gtkbindings.Label {
	return managed("Label", func() *gtk.Label {
		return gtk.NewLabel(text)
	})
}

func MenuButton() gtkbindings.MenuButton {
	return managed("MenuButton", func() *gtk.MenuButton {
		return gtk.NewMenuButton()
	})
}

func Picture() gtkbindings.Picture {
	return managed("Picture", func() *gtk.Picture {
		return gtk.NewPicture()
	})
}

func Popover(child any) gtkbindings.Popover {
	return managed("Popover", func() *gtk.Popover {
		popover := gtk.NewPopover()
		popover.SetChild(gtkbindings.ResolveWidget(child))
		return popover
	})
}

func Scale(orientation gtk.Orientation) gtkbindings.Scale {
	return managed("Scale", func() *gtk.Scale {
		scale := gtk.NewScale(orientation, nil)
		scale.ConnectChangeValue(&callback.RangeChangeValueCallback)
		return scale
	})
}

func ScrolledWindow() gtkbindings.ScrolledWindow {
	return managed("ScrolledWindow", func() *gtk.ScrolledWindow {
		scrolledWindow := gtk.NewScrolledWindow()
		scrolledWindow.ConnectEdgeReached(&callback.ScrolledWindowEdgeReachedCallback)
		return scrolledWindow
	})
}

func SearchEntry() gtkbindings.SearchEntry {
	return managed("SearchEntry", func() *gtk.SearchEntry {
		searchEntry := gtk.NewSearchEntry()
		searchEntry.ConnectActivate(&callback.SearchEntryActivateCallback)
		searchEntry.ConnectSearchChanged(&callback.SearchChangedCallback)
		return searchEntry
	})
}

func Spinner() gtkbindings.Spinner {
	return managed("Spinner", func() *gtk.Spinner {
		spinner := gtk.NewSpinner()
		spinner.Start()
		return spinner
	})
}

func Widget(w *gtk.Widget) gtkbindings.Widget {
	return func() *gtkbindings.WrappedWidget {
		return &gtkbindings.WrappedWidget{Widget: *w}
	}
}

// WARN: Do not manage reference counting for a schwifty-managed widget. If you are not in control of the widget's lifecycle, use Widget() instead.
func ManagedWidget(w *gtk.Widget) gtkbindings.Widget {
	return managed("ManagedWidget", func() *gtkbindings.WrappedWidget {
		return &gtkbindings.WrappedWidget{Widget: *w}
	})
}
