package my_collection

import (
	"context"
	"log/slog"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
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

type itemizedUserPlaylistsCollection struct {
	*user_collections.UserCollections
}

func (i *itemizedUserPlaylistsCollection) Items(ctx context.Context, id, cursor string, include ...string) (*openapi.Response[[]openapi.Relationship], error) {
	return i.Playlists(ctx, id, cursor, include...)
}

func Playlists() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()
	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			View:      Label(gettext.Get("Please log in to view your collection")),
		}
	}

	itemizedCollection := &itemizedUserPlaylistsCollection{tidal.OpenAPI.V2.UserCollections}
	paginator := pagination.NewPaginator(itemizedCollection, userId, func(r *openapi.Response[[]openapi.Relationship]) []openapi.Playlist {
		return r.Included.Playlists(r.Data...)
	}, "playlists.coverArt", "playlists.ownerProfiles")

	userCollection, err := paginator.GetFirstPage()
	if err != nil {
		return &router.Response{
			PageTitle: gettext.Get("My Collection"),
			Error:     err,
		}
	}

	children := make([]any, 0)
	for _, playlist := range userCollection {
		children = append(children, CenterBox().CenterWidget(media_card.NewPlaylist(&playlist)))
	}

	list := WrapBox(children...).VMargin(20).VAlign(gtk.AlignStartValue).Justify(adw.JustifyFillValue).JustifyLastLine(true).ToGTK()

	return &router.Response{
		PageTitle: gettext.Get("My Playlists"),
		View: ScrolledWindow().
			Child(list).
			ConnectEdgeReached(func(sw gtk.ScrolledWindow, pt gtk.PositionType) {
				if pt == gtk.PosBottomValue {
					go func() {
						if !paginator.IsConsumed() {
							items, err := paginator.Next()
							if err != nil {
								return
							}

							schwifty.OnMainThreadOnce(func(u uintptr) {
								list := adw.WrapBoxNewFromInternalPtr(u)
								for _, playlist := range items {
									child := CenterBox().CenterWidget(media_card.NewPlaylist(&playlist)).ToGTK()
									list.Append(child)
								}
							}, list.GoPointer())
						} else {
							slog.Debug("No more playlists to fetch")
						}
					}()
				}
			}).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
