package v2

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/openapi/v2/albums"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/openapi/v2/artists"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/openapi/v2/playlists"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/openapi/v2/search_results"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/openapi/v2/tracks"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/openapi/v2/user_collections"
)

type V2 struct {
	Albums          *albums.Albums
	Artists         *artists.Artists
	Playlists       *playlists.Playlists
	SearchResults   *search_results.SearchResults
	Tracks          *tracks.Tracks
	UserCollections *user_collections.UserCollections
}

func New(client *internal.Client) *V2 {
	return &V2{
		Albums:          albums.New(client),
		Artists:         artists.New(client),
		Playlists:       playlists.New(client),
		SearchResults:   search_results.New(client),
		Tracks:          tracks.New(client),
		UserCollections: user_collections.New(client),
	}
}
