package gtk

import "github.com/jwijenbergh/puregotk/v4/gtk"

//go:generate go run codeberg.org/dergs/tonearm/pkg/schwifty/gen Grid *gtk.Grid gtk

func (f Grid) Attach(child any, column int, row int, width int, height int) Grid {
	return func() *gtk.Grid {
		grid := f()
		grid.Attach(ResolveWidget(child), column, row, width, height)
		return grid
	}
}

func (f Grid) ColumnHomogeneous(homogeneous bool) Grid {
	return func() *gtk.Grid {
		grid := f()
		grid.SetColumnHomogeneous(homogeneous)
		return grid
	}
}
