package routes

import (
	"errors"
	"fmt"

	"codeberg.org/dergs/tidalwave/internal/router"
	"codeberg.org/dergs/tidalwave/pkg/gui"
)

func init() {
	router.Register("album", Album)
}

func Album(params router.Params) *router.Response {
	albumId, ok := params["id"].(int)
	if !ok {
		return router.FromError("Album", errors.New("invalid album ID"))
	}

	//tidal := injector.MustInject[*tidalapi.TidalAPI]()
	return &router.Response{
		PageTitle: "Album",
		View:      gui.Text(fmt.Sprint(albumId)),
	}
}
