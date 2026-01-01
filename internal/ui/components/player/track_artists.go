package player

import (
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var artistSeparator = Label(", ").
	FontSize(16).
	FontWeight(700).
	LineHeight(1.2).HAlign(gtk.AlignCenterValue)

var artistLabel = artistSeparator.
	Color("#1C71D8").
	TextDecoration("underline")

var linkDisabler = func(gtk.LinkButton) bool {
	return true
}

func trackArtists() schwifty.Widget {
	textView := gtk.NewTextView()
	textView.SetHexpand(true)
	textView.SetVexpand(false)
	textView.SetValign(gtk.AlignStartValue)
	textView.SetHalign(gtk.AlignFillValue)
	textView.SetJustification(gtk.JustifyCenterValue)
	textView.SetWrapMode(gtk.WrapWordValue)
	textView.SetEditable(false)
	textView.SetCursorVisible(false)
	textView.SetSizeRequest(380, -1)

	widgetPtr := textView.GoPointer()
	subscriptionID := player.OnTrackChanged.On(func(trackInfo player.TrackInformation) bool {
		textView := gtk.TextViewNewFromInternalPtr(widgetPtr)
		tag := gtk.NewTextTagTable()
		buffer := gtk.NewTextBuffer(tag)
		tag.Unref()
		textView.SetBuffer(buffer)
		buffer.Unref()

		var iter gtk.TextIter
		buffer.GetStartIter(&iter)

		for i, artist := range trackInfo.Artists {
			anchor := buffer.CreateChildAnchor(&iter)
			btn := gtk.NewLinkButton("https://tidal.com/artist/" + artist.ID)
			btn.SetActionName("win.route.artist")
			btn.SetActionTargetValue(glib.NewVariantString(artist.ID))
			btn.SetChild(artistLabel.Text(artist.Attributes.Name).ToGTK())
			btn.ConnectActivateLink(&linkDisabler)
			textView.AddChildAtAnchor(Widget(&btn.Widget).Padding(0).ToGTK(), anchor)
			btn.Unref()
			if i != len(trackInfo.Artists)-1 {
				anchor := buffer.CreateChildAnchor(&iter)
				textView.AddChildAtAnchor(artistSeparator.ToGTK(), anchor)
			}
		}

		return signals.Continue
	})

	return ManagedWidget(&textView.Widget).Background("transparent").ConnectDestroy(func(w gtk.Widget) {
		player.OnTrackChanged.Unsubscribe(subscriptionID)
	})
}
