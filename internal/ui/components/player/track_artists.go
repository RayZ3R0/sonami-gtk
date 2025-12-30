package player

import (
	"codeberg.org/dergs/tidalwave/internal/player"
	"codeberg.org/dergs/tidalwave/internal/signals"
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/schwifty/state"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var artistsState = state.NewStateful[any](Label(""))

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

func init() {
	player.OnTrackChanged.On(func(trackInfo player.TrackInformation) bool {
		schwifty.OnMainThreadOnce(func(u uintptr) {
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

			buffer := textView.GetBuffer()
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

			artistsState.SetValue(ManagedWidget(&textView.Widget).MinWidth(380).Background("transparent"))
		}, 0)
		return signals.Continue
	})
}

func trackArtists() schwifty.CenterBox {
	return CenterBox().BindCenterWidget(artistsState).HMargin(20)
}
