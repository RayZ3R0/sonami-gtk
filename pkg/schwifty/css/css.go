package css

import "github.com/jwijenbergh/puregotk/v4/gtk"

func Apply(widget *gtk.Widget, css string) {
	widget.Ref()
	defer widget.Unref()

	provider := gtk.NewCssProvider()
	provider.LoadFromString(css)
	widget.GetStyleContext().AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))
	provider.Unref()
}
