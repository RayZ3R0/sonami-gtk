package my_collection

import (
	"context"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/ui/components/media_card"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Playlists() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View:      Label(gettext.Get("Please log in to view your collection")),
		}
	}

	userCollection, err := tidal.OpenAPI.V2.UserCollections.Playlists(context.Background(), userId, "", "playlists.coverArt", "playlists.ownerProfiles")
	if err != nil {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			Error:     err,
		}
	}

	children := make([]any, 0)
	for _, playlist := range userCollection.Included.Playlists(userCollection.Data...) {
		children = append(children, CenterBox().CenterWidget(media_card.NewPlaylist(&playlist)))
	}

	return &router.Response{
		PageTitle: gettext.Get("My Playlists"),
		View: ScrolledWindow().
			Child(
				WrapBox(children...).VMargin(20).VAlign(gtk.AlignStartValue).Justify(adw.JustifyFillValue).JustifyLastLine(true),
			).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
