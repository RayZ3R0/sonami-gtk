package legacy

import (
	"fmt"
	"strconv"

	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

var artistLogger = logger.With("type", "legacy_artist")

type Artist struct {
	API    *tidalapi.TidalAPI
	Artist v1.Artist
}

func (a *Artist) ID() string {
	return strconv.Itoa(a.Artist.ID)
}

func (a *Artist) Name() string {
	return a.Artist.Name
}

func (a *Artist) ProfilePicture(preferredSize int) (string, error) {
	return "", nil
}

func (a *Artist) URL() string {
	return fmt.Sprintf("https://tidal.com/artist/%s", a.ID())
}

func NewArtist(api *tidalapi.TidalAPI, artist v1.Artist) tonearm.Artist {
	return &Artist{
		API:    api,
		Artist: artist,
	}
}
