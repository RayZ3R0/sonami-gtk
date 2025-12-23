package dynamic

import (
	"context"

	. "codeberg.org/dergs/tidalwave/pkg/gui"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
	"github.com/diamondburned/gotkit/gtkutil/imgutil"
	"github.com/infinytum/injector"
)

var trackListEntryCSS = cssutil.Applier("track-list-entry", `
	.track-list-entry {
		padding: 10px;
		border-radius: 10px;
		transition-duration: 0.2s;
		transition-property: background-color;
		transition-timing-function: cubic-bezier(0.25, 0.46, 0.45, 0.94);
	}

	.track-list-entry:hover {
		background-color: rgba(255,255,255,0.15);
	}

	.track-list-entry:focus:active {
		background-color: rgba(255,255,255,0.3);
	}
`)

type TrackListEntry struct {
	box             *BoxImpl
	boxGestureClick *gtk.GestureClick

	album  *TextImpl
	artist *TextImpl
	cover  *gtk.Image
	time   *TextImpl
	title  *TextImpl

	addToCollection *gtk.Button
	addToQueue      *gtk.Button
}

func (t *TrackListEntry) AttachToGrid(grid *gtk.Grid, row int) {
	titleAlbumStack := VStack(
		t.title.HAlign(gtk.AlignStart),
		t.album.HAlign(gtk.AlignStart),
	).Spacing(3).VAlign(gtk.AlignCenter)

	grid.Attach(titleAlbumStack.MarginLeft(84).HAlign(gtk.AlignStart), 0, row, 3, 1)
	grid.Attach(t.artist.HAlign(gtk.AlignStart), 3, row, 2, 1)
	grid.Attach(t.time.MarginRight(88).HAlign(gtk.AlignEnd), 5, row, 1, 1)
	grid.Attach(t.box, 0, row, 6, 1)
}

func (t *TrackListEntry) SetAlbum(album string) *TrackListEntry {
	t.album.Text(album)
	return t
}

func (t *TrackListEntry) SetArtist(artist string) *TrackListEntry {
	t.artist.Text(artist)
	return t
}

func (t *TrackListEntry) SetCoverFromURL(url string) *TrackListEntry {
	imgutil.AsyncGET(injector.MustInject[context.Context](), url, imgutil.ImageSetterFromImage(t.cover))
	return t
}

func (t *TrackListEntry) SetTitle(title string) *TrackListEntry {
	t.title.Text(title)
	return t
}

func (t *TrackListEntry) SetTime(time string) *TrackListEntry {
	t.time.Text(time)
	return t
}

func NewTrackListEntry() *TrackListEntry {
	trackListEntry := &TrackListEntry{
		boxGestureClick: gtk.NewGestureClick(),

		album:  Text("Album").CSS(`label { color: #939393; }`),
		artist: Text("Artist"),
		cover:  gtk.NewImage(),
		time:   Text("00:00"),
		title:  Text("Title"),

		addToCollection: gtk.NewButtonFromIconName("heart-outline-thick-symbolic"),
		addToQueue:      gtk.NewButtonFromIconName("plus-symbolic"),
	}
	trackListEntry.cover.SetPixelSize(54)
	trackListEntry.box = HStack(
		AspectFrame(trackListEntry.cover).
			CornerRadius(10).
			HAlign(gtk.AlignStart).
			Overflow(gtk.OverflowHidden),
		Spacer(),
		HStack(
			Wrapper(trackListEntry.addToCollection).
				HAlign(gtk.AlignCenter).
				VAlign(gtk.AlignCenter).
				CSS(`button:not(:hover) { background-color: transparent; }`),
			Wrapper(trackListEntry.addToQueue).
				HAlign(gtk.AlignCenter).
				VAlign(gtk.AlignCenter).
				CSS(`button:not(:hover) { background-color: transparent; }`),
		),
	).HExpand(true).Focusable(true).FocusOnClick(true)
	trackListEntryCSS(trackListEntry.box)

	ctrl := gtk.NewGestureClick()
	ctrl.ConnectPressed(func(nPress int, x, y float64) {
		box := Wrapper(ctrl.Widget()).GTKWidget()
		if child := box.FocusChild(); child != nil {
			if Wrapper(child).GTKWidget().StateFlags()&gtk.StateFlagActive != 0 {
				return
			}
		}
		box.GrabFocus()
		box.SetStateFlags(gtk.StateFlagActive, false)
	})
	ctrl.ConnectReleased(func(nPress int, x, y float64) {
		// TODO: Implement "Play Track"
	})
	trackListEntry.box.GTKWidget().AddController(ctrl)

	return trackListEntry
}
