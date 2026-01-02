package v2

import (
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/internal"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/openapi/v2/albums"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/openapi/v2/playlists"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/openapi/v2/search_results"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/openapi/v2/tracks"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi/openapi/v2/user_collections"
)

type V2 struct {
	Albums          *albums.Albums
	Playlists       *playlists.Playlists
	SearchResults   *search_results.SearchResults
	Tracks          *tracks.Tracks
	UserCollections *user_collections.UserCollections
}

func New(client *internal.Client) *V2 {
	return &V2{
		Albums:          albums.New(client),
		Playlists:       playlists.New(client),
		SearchResults:   search_results.New(client),
		Tracks:          tracks.New(client),
		UserCollections: user_collections.New(client),
	}
}
