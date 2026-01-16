package tracklist

import (
	"math"
	"reflect"
	"slices"
	"sync"
	"unsafe"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/factory"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"github.com/go-gst/go-glib/glib"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/gobject/types"
	"github.com/jwijenbergh/puregotk/v4/graphene"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type ColumnFunc[TrackType comparable] func(track TrackType, grid *gtk.Grid, position int, column int) int

type TrackList[TrackType comparable] struct {
	schwifty.Widget

	columnFuncs []ColumnFunc[TrackType]
	lock        sync.Mutex
	store       *gio.ListStore
	trackList   []TrackType

	movingSourceIndex int
	movingTargetIndex int
	reorderable       *state.State[bool]
	reorderCallback   func(sourceIndex, targetIndex int, track TrackType)
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

func (t *TrackList[TrackType]) SetReorderCallback(cb func(sourceIndex, targetIndex int, track TrackType)) {
	t.reorderCallback = cb
	t.reorderable.SetValue(cb != nil)
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

	var dragSource *gtk.DragSource
	t.reorderable.AddCallback(func(newValue bool) {
		if newValue && dragSource == nil {
			dragSource = gtk.NewDragSource()
			dragSource.SetActions(gdk.ActionMoveValue)
			// FIX: Switch to Schwifty-managed Callbacks
			dragSource.ConnectPrepare(g.Ptr(func(_ gtk.DragSource, _, _ float64) gdk.ContentProvider {
				t.movingSourceIndex = int(listItem.GetPosition())
				t.movingTargetIndex = int(listItem.GetPosition())

				track := t.trackList[t.movingSourceIndex]
				return *gdk.NewContentProviderTyped(types.GType(glib.TYPE_POINTER), &track)
			}))
			// FIX: Switch to Schwifty-managed Callbacks
			dragSource.ConnectDragBegin(g.Ptr(func(source gtk.DragSource, dragPtr uintptr) {
				drag := gdk.DragNewFromInternalPtr(dragPtr)
				drag.SetData("t", uintptr(unsafe.Pointer(t)))

				snapshot := gtk.NewSnapshot()
				p := grid.Widget.GetParent()
				defer p.Unref()
				p.SnapshotChild(&grid.Widget, snapshot)
				size := graphene.Size{}
				size.Init(
					float32(grid.Widget.GetAllocatedWidth()),
					float32(grid.Widget.GetAllocatedHeight()),
				)
				paintable := snapshot.ToPaintable(&size)

				source.SetIcon(paintable, paintable.GetIntrinsicWidth()/2, paintable.GetIntrinsicHeight()/2)

				t.trackList = append(t.trackList[:t.movingSourceIndex], t.trackList[t.movingSourceIndex+1:]...)
				t.store.Remove(uint(t.movingSourceIndex))

				trackType := reflect.TypeFor[TrackType]()
				track := reflect.Zero(trackType).Interface().(TrackType)
				t.trackList = append(t.trackList[:t.movingTargetIndex], append([]TrackType{track}, t.trackList[t.movingTargetIndex:]...)...)
				t.store.Insert(uint(t.movingTargetIndex), &gtk.NewStringObject("").Object)
			}))
			// FIX: Switch to Schwifty-managed Callbacks
			dragSource.ConnectDragCancel(g.Ptr(func(source gtk.DragSource, dragPtr uintptr, reason gdk.DragCancelReason) bool {
				drag := gdk.DragNewFromInternalPtr(dragPtr)
				content := drag.GetContent()
				defer content.Unref()

				value := gobject.Value{}
				value.Init(gobject.TypePointerVal)
				content.GetValue(&value)
				track := *(*TrackType)(unsafe.Pointer(value.GetPointer()))

				t.trackList = append(t.trackList[:t.movingTargetIndex], t.trackList[t.movingTargetIndex+1:]...)
				t.store.Remove(uint(t.movingTargetIndex))

				t.trackList = append(t.trackList[:t.movingSourceIndex], append([]TrackType{track}, t.trackList[t.movingSourceIndex:]...)...)
				t.store.Insert(uint(t.movingSourceIndex), &gtk.NewStringObject("").Object)

				t.movingSourceIndex = -1
				t.movingTargetIndex = -1

				return true
			}))
			grid.AddController(&dragSource.EventController)
		} else if dragSource != nil {
			grid.RemoveController(&dragSource.EventController)
			dragSource.Unref()
			dragSource = nil
		}
	})

	listItem.SetChild(&grid.Widget)
	listItem.SetActivatable(false)
}

func NewTrackList[TrackType comparable](columnFuncs ...ColumnFunc[TrackType]) *TrackList[TrackType] {
	t := &TrackList[TrackType]{
		columnFuncs: columnFuncs,
		trackList:   make([]TrackType, 0),
		store:       gio.NewListStore(gtk.StringObjectGLibType()),

		movingSourceIndex: -1,
		movingTargetIndex: -1,
		reorderable:       state.NewStateful(false),
	}

	factory := factory.NewSignalListItemFactory().
		ConnectSetup(t.onSetup).
		ConnectBind(t.onBind)()

	listView := gtk.NewListView(gtk.NewNoSelection(t.store), &factory.ListItemFactory)
	listView.SetSingleClickActivate(false)
	listView.SetOrientation(gtk.OrientationVerticalValue)

	var dropTarget *gtk.DropTargetAsync
	t.reorderable.AddCallback(func(newValue bool) {
		if newValue && dropTarget == nil {
			dropTarget = gtk.NewDropTargetAsync(gdk.NewContentFormatsForGtype(types.GType(glib.TYPE_POINTER)), gdk.ActionMoveValue)
			// FIX: Switch to Schwifty-managed Callbacks
			dropTarget.ConnectAccept(g.Ptr(func(_ gtk.DropTargetAsync, dropPtr uintptr) bool {
				drop := gdk.DropNewFromInternalPtr(dropPtr)

				drag := drop.GetDrag()
				defer drag.Unref()
				if drag.GetData("t") != uintptr(unsafe.Pointer(t)) {
					return false
				}

				return true
			}))
			// FIX: Switch to Schwifty-managed Callbacks
			dropTarget.ConnectDragMotion(g.Ptr(func(target gtk.DropTargetAsync, dropPtr uintptr, x, y float64) gdk.DragAction {
				height := listView.GetFirstChild().GetAllocatedHeight()
				i := int(math.Floor(y / float64(height)))

				if i != t.movingTargetIndex {
					if i > len(t.trackList) {
						return gdk.ActionNoneValue
					}
					trackList := slices.Clone(t.trackList)
					if t.movingTargetIndex != -1 {
						t.store.Remove(uint(t.movingTargetIndex))
						trackList = append(t.trackList[:t.movingTargetIndex], t.trackList[t.movingTargetIndex+1:]...)
					}

					trackType := reflect.TypeFor[TrackType]()
					track := reflect.Zero(trackType).Interface().(TrackType)
					t.trackList = append(trackList[:i], append([]TrackType{track}, trackList[i:]...)...)
					t.store.Insert(uint(i), &gtk.NewStringObject("").Object)

					t.movingTargetIndex = i
				}

				return gdk.ActionMoveValue
			}))
			// FIX: Switch to Schwifty-managed Callbacks
			dropTarget.ConnectDrop(g.Ptr(func(target gtk.DropTargetAsync, dropPtr uintptr, x, y float64) bool {
				drop := gdk.DropNewFromInternalPtr(dropPtr)
				value := gobject.Value{}
				value.Init(gobject.TypePointerVal)
				drop.GetDrag().GetContent().GetValue(&value)
				track := *(*TrackType)(unsafe.Pointer(value.GetPointer()))
				drop.Finish(gdk.ActionMoveValue)

				t.reorderCallback(t.movingSourceIndex, t.movingTargetIndex, track)

				t.movingSourceIndex = -1
				t.movingTargetIndex = -1

				return true
			}))
			listView.AddController(&dropTarget.EventController)
			// FIX: Switch to Schwifty-managed Callbacks
			listView.ConnectDestroy(g.Ptr(func(w gtk.Widget) {
				w.RemoveController(&dropTarget.EventController)
			}))
		} else if dropTarget != nil {
			listView.RemoveController(&dropTarget.EventController)
			dropTarget.Unref()
			dropTarget = nil
		}
	})

	t.Widget = ManagedWidget(&listView.Widget).Background("transparent")
	return t
}
