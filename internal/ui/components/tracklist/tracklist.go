package tracklist

import (
	"sync"

	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/factory"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type ColumnFunc[TrackType comparable] func(track TrackType, grid *gtk.Grid, position int, column int) int

type TrackList[TrackType comparable] struct {
	schwifty.Widget

	columnFuncs []ColumnFunc[TrackType]
	lock        sync.Mutex
	store       *gio.ListStore
	trackList   []TrackType
}

func (t *TrackList[TrackType]) AddTrack(track TrackType) {
	t.trackList = append(t.trackList, track)
	t.store.Append(&gtk.NewStringObject("").Object)
}

func (t *TrackList[TrackType]) BindTracks(state *state.State[[]TrackType]) {
	state.AddCallback(func(newValue []TrackType) {
		t.lock.Lock()
		defer t.lock.Unlock()

		for i, newTrack := range newValue {
			if i >= len(t.trackList) {
				t.trackList = append(t.trackList, newTrack)
				t.store.Append(&gtk.NewStringObject("").Object)
			} else if t.trackList[i] != newTrack {
				t.trackList[i] = newTrack
				t.store.Splice(uint(i), 1, []gobject.Object{gtk.NewStringObject("").Object}, 1)
			}
		}
		for i := len(t.trackList) - 1; i+1 > len(newValue); i-- {
			t.trackList = t.trackList[:i]
			t.store.Remove(uint(i))
		}
	})
}

func (t *TrackList[TrackType]) Clear() {
	t.store.RemoveAll()
	t.trackList = make([]TrackType, 0)
}

func (t *TrackList[TrackType]) onBind(_ gtk.SignalListItemFactory, listItem *gtk.ListItem) {
	track := t.trackList[listItem.GetPosition()]
	grid := gtk.GridNewFromInternalPtr(listItem.GetChild().GoPointer())
	defer grid.Unref()

	width := 0
	for _, columnFunc := range t.columnFuncs {
		width += columnFunc(track, grid, int(listItem.GetPosition()), width)
	}
}

func (t *TrackList[TrackType]) onSetup(_ gtk.SignalListItemFactory, listItem *gtk.ListItem) {
	grid := gtk.NewGrid()
	defer grid.Unref()
	grid.SetColumnHomogeneous(true)

	listItem.SetChild(&grid.Widget)
	listItem.SetActivatable(false)
}

func NewTrackList[TrackType comparable](columnFuncs ...ColumnFunc[TrackType]) *TrackList[TrackType] {
	t := &TrackList[TrackType]{
		columnFuncs: columnFuncs,
		trackList:   make([]TrackType, 0),
		store:       gio.NewListStore(gtk.StringObjectGLibType()),
	}

	factory := factory.NewSignalListItemFactory().
		ConnectSetup(t.onSetup).
		ConnectBind(t.onBind)()

	listView := gtk.NewListView(gtk.NewNoSelection(t.store), &factory.ListItemFactory)
	listView.SetSingleClickActivate(false)
	listView.SetOrientation(gtk.OrientationVerticalValue)
	t.Widget = ManagedWidget(&listView.Widget).Background("transparent")
	return t
}
