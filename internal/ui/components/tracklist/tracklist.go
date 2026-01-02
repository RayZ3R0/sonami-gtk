package tracklist

import (
	"slices"
	"strconv"

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

	container       *gtk.Grid
	titleVisibility *state.State[bool]
	titleText       *state.State[string]

	// Will be called per-track to generate their columns
	columnFuncs       []ColumnFunc
	legacyColumnFuncs []LegacyColumnFunc

	// rowMap maps track IDs to their row indices in the grid.
	rowMap []string
}

func (t *TrackList) AddTrack(track *openapi.Track) {
	row := len(t.rowMap)
	t.rowMap = append(t.rowMap, track.Data.ID)

	width := 0
	for _, columnFunc := range t.columnFuncs {
		width += columnFunc(track, t.container, row, width)
	}
}

func (t *TrackList) AddLegacyTrack(track *v2.TrackItemData) {
	row := len(t.rowMap)
	t.rowMap = append(t.rowMap, strconv.Itoa(track.ID))

	width := 0
	for _, columnFunc := range t.legacyColumnFuncs {
		width += columnFunc(track, t.container, row, width)
	}
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

		for _, track := range newValue {
			if slices.Contains(t.rowMap, track.Data.ID) {
				continue
			}
			t.AddTrack(track)
		}
	})
}

func (t *TrackList) Clear() {
	for i := len(t.rowMap) - 1; i >= 0; i-- {
		t.container.RemoveRow(i)
	}
	t.rowMap = []string{}
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

func newTrackList(title string) *TrackList {
	titleText := state.New[string](title)
	titleVisibility := state.New[bool](title != "")
	container := gtk.NewGrid()

	return &TrackList{
		titleVisibility: titleVisibility,
		titleText:       titleText,
		container:       container,
		rowMap:          []string{},
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
			),
			ManagedWidget(&container.Widget),
		).VAlign(gtk.AlignStartValue),
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
