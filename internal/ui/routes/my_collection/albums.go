package my_collection

import (
	"context"
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components/media_card"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/openapi/v2/user_collections"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/pagination"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/adw"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

type itemizedUserAlbumsCollection struct {
	*user_collections.UserCollections
}

func (i *itemizedUserAlbumsCollection) Items(ctx context.Context, id, cursor string, include ...string) (*openapi.Response[[]openapi.Relationship], error) {
	return i.Albums(ctx, id, cursor, include...)
}

func Albums() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View:      Label(gettext.Get("Please log in to view your collection")),
		}
	}

	itemizedCollection := &itemizedUserAlbumsCollection{tidal.OpenAPI.V2.UserCollections}

	paginator := pagination.NewPaginator(itemizedCollection, userId, func(r *openapi.Response[[]openapi.Relationship]) []openapi.Album {
		return r.Included.Albums(r.Data...)
	}, "albums.coverArt", "albums.artists")

	userCollection, err := paginator.GetFirstPage()
	if err != nil {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			Error:     err,
		}
	}

	children := make([]any, 0)
	for _, album := range userCollection {
		children = append(children, CenterBox().CenterWidget(media_card.NewAlbum(&album)))
	}

	list := WrapBox(children...).VMargin(20).VAlign(gtk.AlignStartValue).Justify(adw.JustifyFillValue).JustifyLastLine(true).ToGTK()

	return &router.Response{
		PageTitle: gettext.Get("My Albums"),
		View: ScrolledWindow().
			Child(list).
			ConnectReachEdgeSoon(gtk.PosBottomValue, func() bool {
				if !paginator.IsConsumed() {
					items, err := paginator.Next()
					if err != nil {
						return signals.Continue
					}

					schwifty.OnMainThreadOnce(func(u uintptr) {
						list := adw.WrapBoxNewFromInternalPtr(u)
						for _, album := range items {
							child := CenterBox().CenterWidget(media_card.NewAlbum(&album)).ToGTK()
							list.Append(child)
						}
					}, list.GoPointer())
				} else {
					slog.Debug("No more albums to fetch")
					return signals.Unsubscribe
				}
				return signals.Continue
			}).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
