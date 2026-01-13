package media_card

import (
	"strconv"
	"strings"

	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
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
