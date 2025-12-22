package dynamic

import (
	"context"
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotkit/gtkutil/cssutil"
	"github.com/diamondburned/gotkit/gtkutil/imgutil"
	"github.com/infinytum/injector"
	. "github.com/nilathedragon/glinear/pkg/tidalapi/gui"
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
	box    *BoxImpl
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
	cover := gtk.NewImageFromResource("/com/github/nilathedragon/tidal-wave/placeholder/placeholder.jpg")
	cover.SetPixelSize(54)

	addToCollection := gtk.NewButtonFromIconName("heart-outline-thick-symbolic")
	addToCollection.SetHAlign(gtk.AlignCenter)
	addToCollection.SetVAlign(gtk.AlignCenter)
	cssutil.Apply(addToCollection, `button:not(:hover) { background-color: transparent; }`)

	addToQueue := gtk.NewButtonFromIconName("plus-symbolic")
	addToQueue.SetHAlign(gtk.AlignCenter)
	addToQueue.SetVAlign(gtk.AlignCenter)
	cssutil.Apply(addToQueue, `button:not(:hover) { background-color: transparent; }`)

	box := HStack(
		AspectFrame(cover).
			CornerRadius(10).
			HAlign(gtk.AlignStart).
			Overflow(gtk.OverflowHidden),
		Spacer(),
		HStack(
			addToCollection,
			addToQueue,
		),
	).HExpand(true)
	box.GTKWidget().SetFocusable(true)
	box.GTKWidget().SetFocusOnClick(true)
	trackListEntryCSS(box)

	album := Text("Album")
	album.CSS(`label { color: #939393; }`)
	artist := Text("Artist")
	time := Text("00:00")
	title := Text("Title")
	box.GTKWidget().ConnectDestroy(func() {
		fmt.Println("Destroy")
	})

	entry := &TrackListEntry{
		box:    box,
		album:  album,
		artist: artist,
		cover:  cover,
		time:   time,
		title:  title,

		addToCollection: addToCollection,
		addToQueue:      addToQueue,
	}

	box.GTKWidget().ConnectRealize(entry.onRealize)
	box.GTKWidget().ConnectUnrealize(entry.onUnrealize)

	ctrl := gtk.NewGestureClick()
	ctrl.ConnectPressed(entry.onClickStart)
	ctrl.ConnectReleased(func(nPress int, x, y float64) {
		// TODO: Implement "Play Track"
	})
	box.GTKWidget().AddController(ctrl)

	return entry
}

func (t *TrackListEntry) onRealize() {

}

func (t *TrackListEntry) onUnrealize() {

}

func (t *TrackListEntry) onClickStart(nPress int, x, y float64) {
	if t.addToCollection.StateFlags()&gtk.StateFlagActive != 0 || t.addToQueue.StateFlags()&gtk.StateFlagActive != 0 {
		return
	}
	t.box.GTKWidget().GrabFocus()
	t.box.GTKWidget().SetStateFlags(gtk.StateFlagActive, false)
}
