package tracklist

import (
	"slices"
	"strconv"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type ColumnFunc func(track *openapi.Track, grid *gtk.Grid, row int, column int) int
type LegacyColumnFunc func(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int

type TrackList struct {
	schwifty.Box

	container        *gtk.Grid
	titleVisibility  *state.State[bool]
	titleText        *state.State[string]
	routeButtonState *state.State[any]

	// Will be called per-track to generate their columns
	columnFuncs       []ColumnFunc
	legacyColumnFuncs []LegacyColumnFunc

	// rowMap maps track IDs to their row indices in the grid.
	rowMap       []string
	tracksSignal *signals.StatelessSignal[[]string]
}

func (t *TrackList) AddTrack(track *openapi.Track) {
	row := len(t.rowMap)
	t.AddTrackAt(track, row)
	t.tracksSignal.Notify(t.rowMap)
}

func (t *TrackList) AddTrackAt(track *openapi.Track, row int) {
	t.container.InsertRow(row)
	t.rowMap = append(t.rowMap[:row], append([]string{track.Data.ID}, t.rowMap[row:]...)...)

	width := 0
	for _, columnFunc := range t.columnFuncs {
		width += columnFunc(track, t.container, row, width)
	}
	t.tracksSignal.Notify(t.rowMap)
}

func (t *TrackList) AddLegacyTrack(track *v2.TrackItemData) {
	row := len(t.rowMap)
	t.rowMap = append(t.rowMap, strconv.Itoa(track.ID))

	width := 0
	for _, columnFunc := range t.legacyColumnFuncs {
		width += columnFunc(track, t.container, row, width)
	}

	t.tracksSignal.Notify(t.rowMap)
}

func (t *TrackList) BindTracks(state *state.State[[]*openapi.Track]) {
	state.AddCallback(func(newValue []*openapi.Track) {
		if len(newValue) > 0 {
			trackIds := map[string]bool{}
			for _, track := range newValue {
				trackIds[track.Data.ID] = true
			}

			newRowMap := []string{}
			for row, trackId := range t.rowMap {
				if _, ok := trackIds[trackId]; !ok {
					t.container.RemoveRow(row)
				} else {
					newRowMap = append(newRowMap, trackId)
				}
			}
			t.rowMap = newRowMap
		} else {
			t.Clear()
		}

		for expectedRow, track := range newValue {
			if slices.Contains(t.rowMap, track.Data.ID) {
				continue
			}
			t.AddTrackAt(track, expectedRow)
		}

		t.tracksSignal.Notify(t.rowMap)
	})
}

func (t *TrackList) Clear() {
	for i := len(t.rowMap) - 1; i >= 0; i-- {
		t.container.RemoveRow(i)
	}
	t.rowMap = []string{}
	t.tracksSignal.Notify(t.rowMap)
}

func (t *TrackList) SetTitle(title string) *TrackList {
	if title == "" {
		t.titleVisibility.SetValue(false)
		return t
	}
	t.titleText.SetValue(title)
	t.titleVisibility.SetValue(true)
	return t
}

func (t *TrackList) SetViewAllRoute(path string) *TrackList {
	t.routeButtonState.SetValue(Button().Child(
		Label("View All").FontSize(12),
	).
		MinHeight(10).
		MinWidth(10).
		HPadding(10).
		VAlign(gtk.AlignCenterValue).
		ConnectClicked(func(b gtk.Button) {
			router.Navigate(path)
		}))
	return t
}

func newTrackList(title string) *TrackList {
	titleText := state.New[string](title)
	titleVisibility := state.New[bool](title != "")
	routeButtonState := state.NewStateful[any](nil)
	container := gtk.NewGrid()
	tracksSignal := signals.NewStatelessSignal[[]string]()

	return &TrackList{
		titleVisibility:  titleVisibility,
		titleText:        titleText,
		routeButtonState: routeButtonState,
		container:        container,
		rowMap:           []string{},
		tracksSignal:     tracksSignal,
		Box: VStack(
			HStack(
				Label(title).
					BindText(titleText).
					BindVisible(titleVisibility).
					VAlign(gtk.AlignCenterValue).
					MarginStart(10).
					MarginBottom(10).
					FontWeight(600).
					FontSize(20).
					Visible(title != ""),
				Spacer().VExpand(false),
				CenterBox().BindCenterWidget(routeButtonState).HExpand(false).VExpand(false),
			),
			ManagedWidget(&container.Widget),
		).
			VAlign(gtk.AlignStartValue).
			ConnectConstruct(func(b *gtk.Box) {
				ptr := b.GoPointer()
				tracksSignal.On(func(newVal []string) bool {
					b := gtk.BoxNewFromInternalPtr(ptr)

					if len(newVal) == 0 {
						b.Hide()
					} else {
						b.Show()
					}

					return signals.Continue
				})
			}),
	}
}

func NewTrackList(title string, columns ...ColumnFunc) *TrackList {
	trackList := newTrackList(title)
	trackList.columnFuncs = columns
	return trackList
}

func NewLegacyTrackList(title string, columns ...LegacyColumnFunc) *TrackList {
	trackList := newTrackList(title)
	trackList.legacyColumnFuncs = columns
	return trackList
}
