package shortcut_list

import (
	"time"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
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
	return newDeepLink(deepLink.Id, deepLink.Title, deepLink.ExternalURL, deepLink.URL, "")
}
