package horizontal_list

import (
	"strconv"
	"strings"
	"time"

	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
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
