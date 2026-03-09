package shortcut_list

import (
	"time"

	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/router"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty"
	v2 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v2"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/infinytum/injector"
)

func newDeepLink(id string, name string, external bool, url string, coverUrl string) schwifty.Button {
	subtitle := ""
	if external {
		subtitle = gettext.Get("External Link")
	}
	return NewShortcut(
		name,
		subtitle,
		coverUrl,
	).ConnectClicked(func(b gtk.Button) {
		if external {
			gtk.ShowUri(injector.MustInject[*gtk.Window](), url, uint32(time.Now().Unix()))
		} else {
			router.Navigate(url)
		}
	})
}

func NewLegacyDeepLink(deepLink *v2.DeepLinkItemData) schwifty.Button {
	coverUrl := ""
	if deepLink.Id == "tidal://my-collection/tracks" {
		coverUrl = "https://tidal.com/assets/my-tracks-DTG3pLQW.png"
	}

	return newDeepLink(deepLink.Id, deepLink.Title, deepLink.ExternalURL, deepLink.URL, coverUrl)
}
