package tracklist

import (
	"math"
	"slices"
	"sync"
	"unsafe"

	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/factory"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/schwifty/tracking"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gdk"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/gobject"
	"github.com/jwijenbergh/puregotk/v4/graphene"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type TrackWithID interface {
	comparable
	GetID() string
}

type ColumnFunc[TrackType TrackWithID] func(track TrackType, grid *gtk.Grid, position int, column int) int

type TrackList[TrackType TrackWithID] struct {
	schwifty.Widget

	columnFuncs []ColumnFunc[TrackType]
	lock        sync.Mutex
	store       *gio.ListStore
	trackList   []TrackType

	trackBeingMoved   *TrackType
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
	container := adw.BinNewFromInternalPtr(listItem.GetChild().GoPointer())
	defer container.Unref()

	var subscription *signals.Subscription
	var ref *tracking.WeakRef
	grid := Grid().ConnectConstruct(func(g *gtk.Grid) {
		ref = tracking.NewWeakRef(g)
		subscription = player.TrackChanged.On(func(t *player.Track) bool {
			if obj := ref.Get(); obj != nil {
				grid := gtk.GridNewFromInternalPtr(obj.Ptr)
				if t != nil && track.GetID() == t.ID {
					grid.AddCssClass("playing")
				} else {
					grid.RemoveCssClass("playing")
				}
			}
			return signals.Continue
		})
	}).ConnectDestroy(func(w gtk.Widget) {
		player.TrackChanged.Unsubscribe(subscription)
	})()
	grid.SetColumnHomogeneous(true)

	width := 0
	for _, columnFunc := range t.columnFuncs {
		width += columnFunc(track, grid, int(listItem.GetPosition()), width)
	}

	container.SetChild(&grid.Widget)
}

func (t *TrackList[TrackType]) onUnbind(_ gtk.SignalListItemFactory, listItem *gtk.ListItem) {
	container := adw.BinNewFromInternalPtr(listItem.GetChild().GoPointer())
	defer container.Unref()
	container.SetChild(nil)
}

func (t *TrackList[TrackType]) onSetup(_ gtk.SignalListItemFactory, listItem *gtk.ListItem) {
	container := Bin()()
	defer container.Unref()

	var dragSource *gtk.DragSource
	t.reorderable.AddCallback(func(newValue bool) {
		if newValue && dragSource == nil {
			dragSource = DragSource().
				Actions(gdk.ActionMoveValue).
				ConnectPrepare(func(dragSource gtk.DragSource, x, y float64) gdk.ContentProvider {
					t.movingSourceIndex = int(listItem.GetPosition())
					t.movingTargetIndex = int(listItem.GetPosition())

					track := t.trackList[t.movingSourceIndex]
					t.trackBeingMoved = &track
					return *gdk.NewContentProviderTyped(gobject.TypePointerVal, &track)
				}).
				ConnectDragBegin(func(dragSource gtk.DragSource, drag gdk.Drag) {
					drag.SetData("t", uintptr(unsafe.Pointer(t)))

					snapshot := gtk.NewSnapshot()
					p := container.Widget.GetParent()
					defer p.Unref()
					p.SnapshotChild(&container.Widget, snapshot)
					size := graphene.Size{}
					size.Init(
						float32(container.Widget.GetAllocatedWidth()),
						float32(container.Widget.GetAllocatedHeight()),
					)
					paintable := snapshot.ToPaintable(&size)

					dragSource.SetIcon(paintable, paintable.GetIntrinsicWidth()/2, paintable.GetIntrinsicHeight()/2)
				}).
				ConnectDragCancel(func(dragSource gtk.DragSource, drag gdk.Drag, reason gdk.DragCancelReason) bool {
					content := drag.GetContent()
					defer content.Unref()

					value := gobject.Value{}
					value.Init(gobject.TypePointerVal)
					content.GetValue(&value)

					t.trackList = append(t.trackList[:t.movingTargetIndex], t.trackList[t.movingTargetIndex+1:]...)
					t.store.Remove(uint(t.movingTargetIndex))

					t.trackList = append(t.trackList[:t.movingSourceIndex], append([]TrackType{*t.trackBeingMoved}, t.trackList[t.movingSourceIndex:]...)...)
					t.store.Insert(uint(t.movingSourceIndex), &gtk.NewStringObject("").Object)

					t.movingSourceIndex = -1
					t.movingTargetIndex = -1
					t.trackBeingMoved = nil

					return true
				})()
			container.AddController(&dragSource.EventController)
		} else if dragSource != nil {
			container.RemoveController(&dragSource.EventController)
			dragSource = nil
		}
	})

	listItem.SetChild(&container.Widget)
	listItem.SetActivatable(false)
}

func NewTrackList[TrackType TrackWithID](columnFuncs ...ColumnFunc[TrackType]) *TrackList[TrackType] {
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
		ConnectBind(t.onBind).
		ConnectUnbind(t.onUnbind)()

	listView := gtk.NewListView(gtk.NewNoSelection(t.store), &factory.ListItemFactory)
	listView.SetSingleClickActivate(false)
	listView.SetOrientation(gtk.OrientationVerticalValue)

	var dropTarget *gtk.DropTargetAsync
	t.reorderable.AddCallback(func(newValue bool) {
		if newValue && dropTarget == nil {
			dropTarget = DropTargetAsync(gdk.NewContentFormatsForGtype(gobject.TypePointerVal), gdk.ActionMoveValue).
				ConnectAccept(func(dropTarget gtk.DropTargetAsync, drop gdk.Drop) bool {
					drag := drop.GetDrag()
					defer drag.Unref()
					if drag.GetData("t") != uintptr(unsafe.Pointer(t)) {
						return false
					}

					return true
				}).
				ConnectDragMotion(func(dropTarget gtk.DropTargetAsync, drop gdk.Drop, x, y float64) gdk.DragAction {
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

						t.trackList = append(trackList[:i], append([]TrackType{*t.trackBeingMoved}, trackList[i:]...)...)
						t.store.Insert(uint(i), &gtk.NewStringObject("").Object)

						t.movingTargetIndex = i
					}

					return gdk.ActionMoveValue
				}).
				ConnectDrop(func(dropTarget gtk.DropTargetAsync, drop gdk.Drop, x, y float64) bool {
					drop.Finish(gdk.ActionMoveValue)

					t.reorderCallback(t.movingSourceIndex, t.movingTargetIndex, *t.trackBeingMoved)

					t.movingSourceIndex = -1
					t.movingTargetIndex = -1

					return true
				})()
			listView.AddController(&dropTarget.EventController)
		} else if dropTarget != nil {
			listView.RemoveController(&dropTarget.EventController)
			dropTarget = nil
		}
	})

	t.Widget = ManagedWidget(&listView.Widget).WithCSSClass("tracklist")
	return t
}
