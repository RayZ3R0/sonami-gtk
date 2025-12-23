package dynamic

import (
	. "codeberg.org/dergs/tidalwave/pkg/gui"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type TrackList struct {
	*BoxImpl
	container *gtk.Grid
	title     *TextImpl
}

func (t *TrackList) Append(child *TrackListEntry, row int) *TrackList {
	child.AttachToGrid(t.container, row)
	return t
}

func (t *TrackList) SetTitle(title string) *TrackList {
	t.title.Text(title)
	return t
}

func NewTrackList() *TrackList {
	titleLabel := Text("Track List")
	container := gtk.NewGrid()
	return &TrackList{
		BoxImpl: VStack(
			HStack(
				titleLabel.
					VAlign(gtk.AlignCenter).
					MarginLeft(10).
					MarginBottom(10).
					FontWeight(600).
					FontSize(20),
				Spacer(),
			),
			container,
		).HExpand(true),
		container: container,
		title:     titleLabel,
	}
}
