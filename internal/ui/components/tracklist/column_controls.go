package tracklist

import (
	"fmt"
	"log/slog"
	"strconv"

	"codeberg.org/dergs/tonearm/internal/g"
	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/state"
	favouritebutton "codeberg.org/dergs/tonearm/internal/ui/components/favourite_button"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/models/openapi"
	v2 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v2"
	"github.com/jwijenbergh/puregotk/v4/gio"
	"github.com/jwijenbergh/puregotk/v4/glib"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var logger = slog.With("module", "components/tracklist")

type lightArtist struct {
	Name string
	ID   string
}

func controlsColumn(trackId, albumId string, artistId []lightArtist, grid *gtk.Grid, position int, column int) int {
	model := gio.NewMenu()

	item := gio.NewMenuItem(gettext.Get("Navigate to Album"), "win.route.album")
	item.SetActionAndTargetValue("win.route.album", glib.NewVariantString(albumId))
	model.AppendItem(item)

	if len(artistId) > 1 {
		submenu := gio.NewMenu()
		for _, artist := range artistId {
			item := gio.NewMenuItem(fmt.Sprintf(gettext.Get("Navigate to %s"), artist.Name), "win.route.artist")
			item.SetActionAndTargetValue("win.route.artist", glib.NewVariantString(artist.ID))
			submenu.AppendItem(item)
		}
		model.AppendSubmenu(gettext.Get("Navigate to Artist"), &submenu.MenuModel)
	} else if len(artistId) == 1 {
		item := gio.NewMenuItem(gettext.Get("Navigate to Artist"), "win.route.artist")
		item.SetActionAndTargetValue("win.route.artist", glib.NewVariantString(artistId[0].ID))
		model.AppendItem(item)
	}

	popover := gtk.NewPopoverMenuFromModel(&model.MenuModel)

	grid.Attach(
		HStack(
			favouritebutton.FavouriteButton(state.TracksCache, trackId).
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue),
			Button().
				TooltipText(gettext.Get("Add to Queue")).
				IconName("queue-symbolic").
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue).
				ActionName("win.player.queue-track").
				ActionTargetValue(glib.NewVariantString(trackId)).
				WithCSSClass("flat"),
			MenuButton().
				TooltipText(gettext.Get("More…")).
				IconName("view-more-symbolic").
				HAlign(gtk.AlignCenterValue).
				VAlign(gtk.AlignCenterValue).
				Popover(popover).
				WithCSSClass("flat"),
		).
			Margin(10).
			HAlign(gtk.AlignEndValue).
			ToGTK(),
		column,
		0,
		1,
		1,
	)
	return 1
}

func ControlsColumn(track *openapi.Track, grid *gtk.Grid, position int, column int) int {
	if track == nil {
		grid.Attach(
			Box(gtk.OrientationHorizontalValue).ToGTK(),
			column,
			0,
			1,
			1,
		)
		return 1
	}
	return controlsColumn(
		track.Data.ID,
		track.Included.Albums(track.Data.Relationships.Albums.Data...)[0].Data.ID,
		g.Map(
			track.Included.Artists(track.Data.Relationships.Artists.Data...),
			func(artist openapi.Artist) lightArtist {
				return lightArtist{
					Name: artist.Data.Attributes.Name,
					ID:   artist.Data.ID,
				}
			},
		),
		grid,
		position,
		column,
	)
}

func LegacyControlsColumn(track *v2.TrackItemData, grid *gtk.Grid, position int, column int) int {
	if track == nil {
		grid.Attach(
			Box(gtk.OrientationHorizontalValue).ToGTK(),
			column,
			0,
			1,
			1,
		)
		return 1
	}
	return controlsColumn(
		strconv.Itoa(track.ID),
		strconv.Itoa(track.Album.ID),
		g.Map(track.Artists, func(artist v2.TrackItemDataArtist) lightArtist {
			return lightArtist{
				Name: artist.Name,
				ID:   strconv.Itoa(artist.ID),
			}
		}),
		grid,
		position,
		column,
	)
}
