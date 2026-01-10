package media_card

import (
	"strconv"

	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

func NewArtistGeneric(id string, name string, coverUrl string) schwifty.Button {
	return Card(
		name,
		HStack(),
		coverUrl,
	).ActionName("win.route.artist").ActionTargetValue(glib.NewVariantString(id))
}

func NewLegacyArtist(artist *v2.ArtistItemData) schwifty.Button {
	return NewArtistGeneric(strconv.Itoa(artist.Id), artist.Name, tidalapi.ImageURL(artist.Picture))
}

func NewArtist(artist *openapi.Artist) schwifty.Button {
	coverUrl := ""
	for _, artwork := range artist.Included.PlainArtworks(artist.Data.Relationships.ProfileArt.Data...) {
		coverUrl = artwork.Attributes.Files.AtLeast(160).Href
		break
	}
	return NewArtistGeneric(artist.Data.ID, artist.Data.Attributes.Name, coverUrl)
}
