package factory

import (
	"codeberg.org/dergs/tonearm/pkg/schwifty/callback"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var (
	FactorySetupCallback = func(factory gtk.SignalListItemFactory, itemPtr uintptr) {
		callback.CallbackHandler[any](factory.Object, "setup", factory, gtk.ListItemNewFromInternalPtr(itemPtr))
	}
	FactoryBindCallback = func(factory gtk.SignalListItemFactory, itemPtr uintptr) {
		callback.CallbackHandler[any](factory.Object, "bind", factory, gtk.ListItemNewFromInternalPtr(itemPtr))
	}
)

type SignalListItemFactory func() *gtk.SignalListItemFactory

func (f SignalListItemFactory) ConnectBind(cb func(gtk.SignalListItemFactory, *gtk.ListItem)) SignalListItemFactory {
	return func() *gtk.SignalListItemFactory {
		factory := f()
		callback.HandleCallback(factory.Object, "bind", cb)
		return factory
	}
}

func (f SignalListItemFactory) ConnectSetup(cb func(gtk.SignalListItemFactory, *gtk.ListItem)) SignalListItemFactory {
	return func() *gtk.SignalListItemFactory {
		factory := f()
		callback.HandleCallback(factory.Object, "setup", cb)
		return factory
	}
}

func NewSignalListItemFactory() SignalListItemFactory {
	return func() *gtk.SignalListItemFactory {
		factory := gtk.NewSignalListItemFactory()
		factory.ConnectSetup(&FactorySetupCallback)
		factory.ConnectBind(&FactoryBindCallback)
		return factory
	}
}
