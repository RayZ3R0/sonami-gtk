package my_collection

import (
	"context"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/internal/secrets"
	"codeberg.org/dergs/tidalwave/internal/ui/components/media_card"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func Playlists(params router.Params) *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: "My Collection",
			View:      Label("Please log in to view your collection"),
		}
	}

	userCollection, err := tidal.OpenAPI.V2.UserCollections.Playlists(context.Background(), userId, "", "playlists.coverArt", "playlists.ownerProfiles")
	if err != nil {
		return &router.Response{
			PageTitle: "My Collection",
			Error:     err,
		}
	}

	children := make([]any, 0)
	for _, playlist := range userCollection.Included.Playlists(userCollection.Data...) {
		children = append(children, CenterBox().CenterWidget(media_card.NewPlaylist(&playlist)))
	}

	return &router.Response{
		PageTitle: "My Collection",
		View: ScrolledWindow().
			Child(
				WrapBox(children...).VMargin(20).VAlign(gtk.AlignStartValue).Justify(adw.JustifyFillValue).JustifyLastLine(true),
			).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
