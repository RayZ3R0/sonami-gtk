package routes

import (
	"context"
	"errors"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/pkg/gui"
	"codeberg.org/dergs/tidalwave/pkg/tidalapi"
	"github.com/infinytum/injector"
)

func init() {
	router.Register("playlist", Playlist)
}

func Playlist(params router.Params) *router.Response {
	playlistUUID, ok := params["uuid"].(string)
	if !ok {
		return router.FromError("Playlist", errors.New("invalid playlist ID"))
	}

	tidal := injector.MustInject[*tidalapi.TidalAPI]()

	playlist, err := tidal.V1.Playlists.Playlist(context.Background(), playlistUUID)
	if err != nil {
		return router.FromError("Playlist", err)
	}

	return &router.Response{
		PageTitle: playlist.Title,
		View:      gui.Text(playlist.Description),
	}
}
