package state

import (
	"context"
	"sync"
	"time"

	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
	"github.com/infinytum/injector"
)

var (
	favouritesEtag      string
	favourites          *v1.FavouritesIdLists
	favouritesTimeStamp time.Time
	favouritesMutex     sync.RWMutex
)

func Favourites() (*v1.FavouritesIdLists, error) {
	favouritesMutex.RLock()
	if favourites != nil && favouritesTimeStamp.Add(5*time.Minute).After(time.Now()) {
		favouritesMutex.RUnlock()
		return favourites, nil
	}

	favouritesMutex.RUnlock()

	favouritesMutex.Lock()
	defer favouritesMutex.Unlock()

	tidal, _ := injector.Inject[*tidalapi.TidalAPI]()
	userID := secrets.UserID()

	list, etag, noneModified, err := tidal.V1.Favourites.IDsWithCache(context.Background(), userID, favouritesEtag)

	if err != nil {
		return nil, err
	}

	favouritesTimeStamp = time.Now()

	if noneModified {
		return favourites, nil
	}

	favourites = list
	favouritesEtag = etag

	return list, nil
}
