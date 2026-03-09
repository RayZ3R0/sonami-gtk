package v2

import (
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/internal"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/v2/artist"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/v2/favourites"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/v2/feed"
	"github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/v2/home"
)

type V2 struct {
	client     *internal.Client
	Artist     *artist.Artist
	Home       *home.Home
	Favourites *favourites.Favourites
	Feed       *feed.Feed
}

func New(client *internal.Client) *V2 {
	return &V2{
		client:     client,
		Artist:     artist.New(client),
		Home:       home.New(client),
		Favourites: favourites.New(client),
		Feed:       feed.New(client),
	}
}
