package media_card

import (
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

func NewPlaylistGeneric(id string, title string, creator string, itemCount int, coverUrl string) schwifty.Button {
	return Card(
		title,
		VStack(
			SubTitle(creator),
			SubTitle(gettext.GetN("%d Track", "%d Tracks", itemCount, itemCount)),
		),
		coverUrl,
	).ActionName("win.route.playlist").ActionTargetValue(glib.NewVariantString(id))
}

func NewLegacyPlaylist(playlist *v2.PlaylistItemData) schwifty.Button {
	creator := "TIDAL"
	if playlist.Creator.Name != "" {
		creator = playlist.Creator.Name
	}
	return NewPlaylistGeneric(playlist.UUID, playlist.Title, creator, playlist.NumberOfTracks, tidalapi.ImageURL(playlist.SquareImage))
}

func NewPlaylist(playlist *openapi.Playlist) schwifty.Button {
	coverUrl := ""
	for _, artwork := range playlist.Included.PlainArtworks(playlist.Data.Relationships.CoverArt.Data...) {
		if artwork.Attributes.IsPicture() {
			coverUrl = artwork.Attributes.Files.AtLeast(160).Href
			break
		}
	}
	creator := "TIDAL"
	for _, profile := range playlist.Included.PlainArtists(playlist.Data.Relationships.OwnerProfiles.Data...) {
		creator = profile.Attributes.Name
		break
	}
	return NewPlaylistGeneric(playlist.Data.ID, playlist.Data.Attributes.Name, creator, playlist.Data.Attributes.NumberOfItems, coverUrl)
}
