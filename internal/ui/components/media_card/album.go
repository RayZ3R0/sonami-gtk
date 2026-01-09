package media_card

import (
	"strconv"
	"strings"
	"time"

	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

func newAlbum(id string, title string, artists string, year string, coverUrl string) schwifty.Button {
	return Card(
		title,
		VStack(
			SubTitle(artists),
			SubTitle(year),
		),
		coverUrl,
	).ActionName("win.route.album").ActionTargetValue(glib.NewVariantString(id))
}

func NewLegacyAlbum(album *v2.AlbumItemData) schwifty.Button {
	artists := make([]string, 0)
	for _, artist := range album.Artists {
		artists = append(artists, artist.Name)
	}
	releaseDate, _ := time.Parse(time.DateOnly, album.ReleaseDate)

	return newAlbum(strconv.Itoa(album.Id), album.Title, strings.Join(artists, ", "), releaseDate.Format("2006"), tidalapi.ImageURL(album.Cover))
}

func NewAlbum(album *openapi.Album) schwifty.Button {
	coverUrl := ""
	for _, artwork := range album.Included.PlainArtworks(album.Data.Relationships.CoverArt.Data...) {
		coverUrl = artwork.Attributes.Files.AtLeast(160).Href
		break
	}
	artists := make([]string, 0)
	for _, artist := range album.Included.PlainArtists(album.Data.Relationships.Artists.Data...) {
		artists = append(artists, artist.Attributes.Name)
	}
	return newAlbum(album.Data.ID, album.Data.Attributes.Title, strings.Join(artists, ", "), album.Data.Attributes.ReleaseDate.Format("2006"), coverUrl)
}
