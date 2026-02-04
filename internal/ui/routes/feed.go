package routes

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"codeberg.org/dergs/tonearm/internal/gettext"
	"codeberg.org/dergs/tonearm/internal/router"
	"codeberg.org/dergs/tonearm/internal/secrets"
	"codeberg.org/dergs/tonearm/pkg/schwifty"
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
	"codeberg.org/dergs/tonearm/pkg/tidalapi"
	"codeberg.org/dergs/tonearm/pkg/tidalapi/v2/feed"
	"codeberg.org/dergs/tonearm/pkg/utils/imgutil"
	"github.com/infinytum/injector"
	"github.com/jwijenbergh/puregotk/v4/gtk"
)

func init() {
	router.Register("feed", Feed)
	router.Register("feed/activities", Feed)
}

func makeEntry(widgets ...any) schwifty.Button {
	return Button().
		Child(
			HStack(widgets...).
				VPadding(6).
				Spacing(16),
		).
		HExpand(true)
}

const (
	todayStage = iota
	lastWeekStage
	lastMonthStage
	olderStage
)

func Feed() *router.Response {
	tidal := injector.MustInject[*tidalapi.TidalAPI]()
	userId := secrets.UserID()

	if userId == "" {
		return &router.Response{
			PageTitle: gettext.Get("Feed"),
			View: StatusPage().
				IconName("avatar-default-symbolic").
				Title(gettext.Get("Connection required")).
				Description(gettext.Get("You need to sign in to your account to access this page")),
		}
	}

	activities, err := tidal.V2.Feed.Activities(context.Background(), userId)

	if err != nil {
		return router.FromError(gettext.Get("Feed"), err)
	}

	body := VStack().Spacing(36).VMargin(20)
	isRead := true

	slices.SortFunc(activities, func(a1, a2 *feed.Activity) int {
		// sub a1 from a2 to sort in descending order
		return int(a2.FollowableActivity.OccuredAt.Unix() - a1.FollowableActivity.OccuredAt.Unix())
	})

	stage := todayStage
	box := VStack(Label(gettext.Get("Today")).WithCSSClass("title-2")).Spacing(12)
	hasElements := false

	for _, activity := range activities {
		if activity.FollowableActivity.OccuredAt.Before(time.Now().Add(-24*time.Hour)) && stage == todayStage {
			stage = lastWeekStage
			if hasElements {
				body = body.Append(box)
			}
			box = VStack(Label(gettext.Get("Last Week")).WithCSSClass("title-2")).Spacing(12)
			hasElements = false
		}

		if activity.FollowableActivity.OccuredAt.Before(time.Now().Add(-7*24*time.Hour)) && stage == lastWeekStage {
			stage = lastMonthStage
			if hasElements {
				body = body.Append(box)
			}
			box = VStack(Label(gettext.Get("Last Month")).WithCSSClass("title-2")).Spacing(12)
			hasElements = false
		}

		if activity.FollowableActivity.OccuredAt.Before(time.Now().Add(-30*24*time.Hour)) && stage == lastMonthStage {
			stage = olderStage
			if hasElements {
				body = body.Append(box)
			}
			box = VStack(Label(gettext.Get("Older")).WithCSSClass("title-2")).Spacing(12)
			hasElements = false
		}

		if isRead && !activity.Seen {
			isRead = false
		}

		if !isRead && activity.Seen {
			sep := gtk.NewSeparator(gtk.OrientationVerticalValue)
			box = box.Append(
				ManagedWidget(&sep.Widget).
					CSS("separator { color: var(--accent-color); height: 5px; }"),
			)
			isRead = true
		}

		switch activity.FollowableActivity.ActivityType {
		case feed.ActivityTypeNewAlbumRelease:
			album := activity.FollowableActivity.Album

			artists := []string{}
			for _, artist := range album.Artists {
				artists = append(artists, artist.Name)
			}

			artistsString := strings.Join(artists, ", ")
			if len(artists) >= 2 {
				artistsString = strings.Join(artists[:len(artists)-1], ", ")

				artistsString = fmt.Sprintf(gettext.Get("%s and %s"), artistsString, artists[len(artists)-1])
			}

			subtitle := ""
			switch album.Type {
			case feed.AlbumTypeSingle:
				subtitle = fmt.Sprintf(gettext.Get("Single by %s"), artistsString)
			case feed.AlbumTypeAlbum:
				subtitle = fmt.Sprintf(gettext.Get("Album by %s"), artistsString)
			case feed.AlbumTypeEP:
				subtitle = fmt.Sprintf(gettext.Get("EP by %s"), artistsString)
			default:
				subtitle = gettext.Get("Unknown album type")
			}

			box = box.Append(
				makeEntry(
					AspectFrame(
						Image().
							ConnectConstruct(func(i *gtk.Image) {
								go func() {
									img, _ := injector.Inject[*imgutil.ImgUtil]()
									img.LoadIntoImage(tidalapi.ImageURL(album.Cover), i)
								}()
							}).
							PixelSize(54),
					).
						Overflow(gtk.OverflowHiddenValue).
						CornerRadius(6),
					VStack(
						Label(album.Title).HAlign(gtk.AlignStartValue).WithCSSClass("heading"),
						Label(subtitle).HAlign(gtk.AlignStartValue).WithCSSClass("body"),
					).
						VAlign(gtk.AlignCenterValue),
				).
					ConnectClicked(func(b gtk.Button) {
						router.Navigate(fmt.Sprintf("album/%s", strconv.Itoa(album.ID)))
					}),
			)
			hasElements = true
		case feed.ActivityTypeNewHistoryMix:
			mix := activity.FollowableActivity.HistoryMix

			img, _ := injector.Inject[*imgutil.ImgUtil]()
			texture, err := img.Load(mix.Images.Medium.Url)
			if err != nil {
				return &router.Response{
					PageTitle: gettext.Get("Feed"),
					Error:     err,
				}
			}

			box = box.Append(
				makeEntry(
					AspectFrame(
						Image().
							FromPaintable(texture).
							PixelSize(54),
					).
						Overflow(gtk.OverflowHiddenValue).
						CornerRadius(6),
					VStack(
						Label(mix.Title).HAlign(gtk.AlignStartValue).WithCSSClass("heading"),
						Label(mix.Subtitle).HAlign(gtk.AlignStartValue).WithCSSClass("body"),
					).
						VAlign(gtk.AlignCenterValue),
				).
					ConnectClicked(func(b gtk.Button) {
						router.Navigate(fmt.Sprintf("playlist/%s", mix.Id))
					}),
			)
			hasElements = true
		default:
			box = box.Append(
				Bin().
					Child(
						Label(gettext.Get("Unsupported activity")),
					),
			)
			hasElements = true
		}
	}

	body = body.Append(box)

	// TODO: Implement read action.
	// Currently, no idea what the API route is. I/We have to wait until
	// an artist releases an album or a mix is released.

	return &router.Response{
		PageTitle: gettext.Get("Feed"),
		View: ScrolledWindow().
			Child(Clamp().Child(body).Orientation(gtk.OrientationHorizontalValue).MaximumSize(800).HPadding(8)).
			Policy(gtk.PolicyNeverValue, gtk.PolicyAutomaticValue),
	}
}
