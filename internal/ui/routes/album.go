package routes

import (
	"strings"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/player"
	"codeberg.org/dergs/tonearm/internal/resources"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/signals"
	"codeberg.org/dergs/tonearm/internal/ui/components/tracklist"
	"codeberg.org/dergs/tonearm/internal/ui/pages"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	"codeberg.org/dergs/tonearm/pkg/schwifty/state"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tonearm"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

var canPlayAlbumState = state.NewStateful(false)

func init() {
	router.Register("album/:id", Album)
	player.PlaybackStateChanged.On(func(ps *player.PlaybackState) bool {
		canPlayAlbumState.SetValue(!ps.Loading)
		return signals.Continue
	})
}

func Album(albumId string) *router.Response {
	service, err := injector.Inject[tonearm.Service]()
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	album, err := service.GetAlbum(albumId)
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	coverUrl := album.Cover(160)

	trackPaginator, err := service.GetAlbumTracks(albumId)
	if err != nil {
		return router.FromError(gettext.Get("Album"), err)
	}

	page, err := pages.NewPaginatedTracklistPage(
		trackPaginator,
		func() *tracklist.TrackList {
			return tracklist.NewTrackList(
				tracklist.GroupedColumn(2, gtk.AlignStartValue, tracklist.PositionColumn, tracklist.TitleColumn),
				tracklist.ArtistsColumn,
				tracklist.ExpandCustomButtonColumn(1, func(trackId string, position, _ int) {
					go player.PlayAlbum(albumId, false, position)
				}),
				tracklist.GroupedColumn(1, gtk.AlignEndValue, tracklist.DurationColumn, tracklist.ControlsColumn),
			)
		}, func(tl *tracklist.TrackList) schwifty.BaseWidgetable {
			return tl.HMargin(30).VAlign(gtk.AlignStartValue)
		},
	)

	return &router.Response{
		PageTitle: album.Title(),
		Error:     err,
		View: VStack(
			HStack(
				AspectFrame(
					Image().
						PixelSize(146).
						FromPaintable(resources.MissingAlbum()).
						ConnectConstruct(func(i *gtk.Image) {
							if coverUrl != "" {
								injector.MustInject[*imgutil.ImgUtil]().LoadIntoImage(coverUrl, i)
							}
						}),
				).CornerRadius(10).Overflow(gtk.OverflowHiddenValue),
				VStack(
					Label(album.Title()).
						WithCSSClass("title-2").
						HAlign(gtk.AlignStartValue),
					Label(strings.Join(album.Artists().Names(), ", ")).
						WithCSSClass("heading").WithCSSClass("dimmed").
						PaddingTop(10).
						HAlign(gtk.AlignStartValue),
					Label(album.ReleasedAt().Format("2006")).
						WithCSSClass("heading").WithCSSClass("dimmed").
						HAlign(gtk.AlignStartValue),
					Label(gettext.GetN("%d Track (%s)", "%d Tracks (%s)", album.Count(), album.Count(), tidalapi.FormatDuration(album.Duration()))).
						WithCSSClass("heading").WithCSSClass("dimmed").
						HAlign(gtk.AlignStartValue).
						MarginTop(10),
				).MarginStart(20).VAlign(gtk.AlignCenterValue),
				Spacer().VExpand(false),
				HStack(
					Button().
						TooltipText(gettext.Get("Shuffle Album")).
						IconName("playlist-shuffle-symbolic").
						MinWidth(81).
						CornerRadius(21).
						Padding(9).
						ConnectClicked(func(b gtk.Button) {
							go player.PlayAlbum(albumId, true, 0)
						}).
						BindSensitive(canPlayAlbumState),
					Button().
						TooltipText(gettext.Get("Play Album")).
						IconName("play-symbolic").
						MinWidth(81).
						CornerRadius(21).
						Padding(9).
						CSS(`
							button {
								background-color: var(--accent-bg-color);
							}

							button:hover {
								background-color: var(--accent-color);
							}
						`).
						ConnectClicked(func(b gtk.Button) {
							go player.PlayAlbum(albumId, false, 0)
						}).
						BindSensitive(canPlayAlbumState),
				).
					VAlign(gtk.AlignCenterValue).
					Spacing(5),
			).HMargin(40),
			page.VExpand(true).MarginTop(20),
		).VMargin(20),
	}
}
