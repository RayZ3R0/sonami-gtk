package horizontal_list

import (
	"strconv"
	"strings"

	"codeberg.org/dergs/tidalwave/pkg/schwifty"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	v2 "codeberg.org/dergs/tidalwave/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/glib"
)

func newTrack(id string, title string, artists string, coverUrl string) schwifty.Button {
	return Card(
		title,
		SubTitle(artists),
		coverUrl,
	).ActionName("win.player.play-track").ActionTargetValue(glib.NewVariantString(id))
}

func NewLegacyTrack(track *v2.TrackItemData) schwifty.Button {
	artists := make([]string, 0)
	for _, artist := range track.Artists {
		artists = append(artists, artist.Name)
	}

	return newTrack(strconv.Itoa(track.ID), track.Title, strings.Join(artists, ", "), tidalapi.ImageURL(track.Album.Cover))
}
