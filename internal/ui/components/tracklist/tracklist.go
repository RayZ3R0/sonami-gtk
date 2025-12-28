package tracklist

import (
	"strconv"

	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type ColumnFunc func(track *openapi.Track, grid *gtk.Grid, row int, column int) int
type LegacyColumnFunc func(track *v2.TrackItemData, grid *gtk.Grid, row int, column int) int

type TrackList struct {
	schwifty.Box

	container *gtk.Grid
	title     *gtk.Label

	// Will be called per-track to generate their columns
	columnFuncs       []ColumnFunc
	legacyColumnFuncs []LegacyColumnFunc

	// rowMap maps track IDs to their row indices in the grid.
	rowMap map[string]int
}

func (t *TrackList) AddTrack(track *openapi.Track) {
	row := len(t.rowMap)
	t.rowMap[track.Data.ID] = row

	width := 0
	for _, columnFunc := range t.columnFuncs {
		width += columnFunc(track, t.container, row, width)
	}
}

func (t *TrackList) AddLegacyTrack(track *v2.TrackItemData) {
	row := len(t.rowMap)
	t.rowMap[strconv.Itoa(track.ID)] = row

	width := 0
	for _, columnFunc := range t.legacyColumnFuncs {
		width += columnFunc(track, t.container, row, width)
	}
}

func (t *TrackList) SetTitle(title string) *TrackList {
	if title == "" {
		t.title.SetVisible(false)
		return t
	}
	t.title.SetText(title)
	t.title.SetVisible(true)
	return t
}

func newTrackList(trackList *TrackList, title string) *TrackList {
	trackList.container = gtk.NewGrid()

	trackList.title = Label(title).
		VAlign(gtk.AlignCenterValue).
		MarginStart(10).
		MarginBottom(10).
		FontWeight(600).
		FontSize(20).
		Visible(title != "")()

	trackList.Box = VStack(
		HStack(
			ManagedWidget(&trackList.title.Widget),
			Spacer().VExpand(false),
		),
		ManagedWidget(&trackList.container.Widget),
	).VAlign(gtk.AlignStartValue)
	return trackList
}

func NewTrackList(title string, columns ...ColumnFunc) *TrackList {
	trackList := &TrackList{
		columnFuncs: columns,
		rowMap:      make(map[string]int),
	}
	return newTrackList(trackList, title)
}

func NewLegacyTrackList(title string, columns ...LegacyColumnFunc) *TrackList {
	trackList := &TrackList{
		legacyColumnFuncs: columns,
		rowMap:            make(map[string]int),
	}
	return newTrackList(trackList, title)
}
