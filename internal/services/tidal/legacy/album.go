package legacy

import (
	"strconv"

	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
)

type Album struct {
	API   *tidalapi.TidalAPI
	Album v1.Album
}

func (a *Album) Artists(hints ...tonearm.FetchHint) (tonearm.Paginator[tonearm.Artist], error) {
	return tonearm.NewArrayPaginator([]tonearm.Artist{}), nil
}

func (a *Album) Cover(preferredSize int) (string, error) {
	return tidalapi.ImageURL(a.Album.Cover), nil
}

func (a *Album) ID() string {
	return strconv.Itoa(a.Album.ID)
}

func (a *Album) URL() string {
	return "https://tidal.com/album/" + a.ID()
}

func (a *Album) Route() string {
	return "album/" + a.ID()
}

func (a *Album) Title() string {
	return a.Album.Title
}

func (a *Album) Tracks(hints ...tonearm.FetchHint) (tonearm.Paginator[tonearm.Track], error) {
	// TODO: Implement track resolution with hints
	return tonearm.NewArrayPaginator([]tonearm.Track{}), nil
}

func NewAlbum(api *tidalapi.TidalAPI, album v1.Album) tonearm.Album {
	return &Album{
		API:   api,
		Album: album,
	}
}
