package ui

import (
	"fmt"
	"strings"

	"codeberg.org/puregotk/puregotk/v4/adw"
	"codeberg.org/puregotk/puregotk/v4/gtk"
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	"github.com/RayZ3R0/sonami-gtk/internal/localdb"
	"github.com/RayZ3R0/sonami-gtk/internal/notifications"
	"github.com/RayZ3R0/sonami-gtk/internal/secrets"
	"github.com/RayZ3R0/sonami-gtk/internal/settings"
	"github.com/RayZ3R0/sonami-gtk/pkg/schwifty/callback"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
)

func (w *Window) PresentSecretServiceError(err *secrets.ServiceError) {
	if settings.General().ShouldHideSecretServiceWarning() {
		return
	}

	// ConnectResponse is broken with puregotk, so we have to manually hack our way
	AlertDialog(err.Title, err.Body).
		WithCSSClass("no-response").
		ConnectConstruct(func(ad *adw.AlertDialog) {
			checkbox := gtk.NewCheckButtonWithLabel(gettext.Get("Don't show again"))
			checkbox.SetHalign(gtk.AlignBaselineCenterValue)
			checkbox.AddCssClass("space-2")

			ad.SetExtraChild(
				VStack(
					checkbox,
					Button().Label(gettext.Get("Continue")).WithCSSClass("destructive-action").VPadding(10).ConnectClicked(func(b gtk.Button) {
						if checkbox.GetActive() {
							settings.General().SetHideSecretServiceWarning(true)
						}
						ad.Close()
					}),
				).Spacing(12).ToGTK(),
			)
			ad.Present(&w.Widget)
		})()
}

// presentAddToPlaylistDialog shows a dialog for adding a track to an
// existing local playlist or creating a new one on the fly.
func (w *Window) presentAddToPlaylistDialog(trackID, coverURL string) {
	playlists, err := localdb.GetAllPlaylists()
	if err != nil {
		notifications.OnToast.Notify(gettext.Get("Could not load playlists"))
		return
	}

	dialog := adw.NewAlertDialog(gettext.Get("Add to Playlist"), "")
	dialog.AddResponse("cancel", gettext.Get("Cancel"))
	dialog.SetCloseResponse("cancel")

	// closeDialog defers dialog.Close to the next main-loop iteration so that
	// the current signal handler can return before GTK tears down the widget tree.
	closeDialog := func() {
		callback.ScheduleOnMainThreadOncePure(func() { dialog.Close() })
	}

	// --- New playlist entry row ---
	newRow := EntryRow().Title(gettext.Get("Playlist Name…"))()

	createAndAdd := func() {
		name := strings.TrimSpace(newRow.GetText())
		if name == "" {
			return
		}
		pl, err := localdb.CreatePlaylist(name)
		if err != nil {
			notifications.OnToast.Notify(gettext.Get("Failed to create playlist"))
			return
		}
		if err := localdb.AddTrackToPlaylist(pl.ID, trackID, coverURL); err != nil {
			notifications.OnToast.Notify(gettext.Get("Failed to add track to playlist"))
			return
		}
		notifications.OnToast.Notify(fmt.Sprintf(gettext.Get("Created \"%s\" and added track"), pl.Name))
		closeDialog()
	}

	// Connect entry-activated via the schwifty callback infrastructure.
	newRow.ConnectEntryActivated(&entryRowEntryActivated)
	callback.HandleCallback(newRow.Object, "entry-activated", func(_ adw.EntryRow) {
		createAndAdd()
	})

	confirmBtn := Button().IconName("list-add-symbolic").WithCSSClass("flat").ConnectClicked(func(_ gtk.Button) {
		createAndAdd()
	})()
	newRow.AddSuffix(&confirmBtn.Widget)

	newGroupTitle := gettext.Get("New Playlist")
	if len(playlists) == 0 {
		newGroupTitle = gettext.Get("Create Your First Playlist")
	}
	newGroup := adw.NewPreferencesGroup()
	newGroup.SetTitle(newGroupTitle)
	newGroup.Add(&newRow.Widget)

	// --- Outer container ---
	container := gtk.NewBox(gtk.OrientationVerticalValue, 12)
	container.SetMarginTop(4)
	container.SetMarginBottom(4)

	// Existing playlists group (shown only when playlists exist).
	if len(playlists) > 0 {
		existingGroup := adw.NewPreferencesGroup()
		existingGroup.SetTitle(gettext.Get("My Playlists"))

		for _, p := range playlists {
			playlist := p // capture loop variable

			row := ActionRow().
				Title(playlist.Name).
				Subtitle(gettext.GetN("%d Track", "%d Tracks", playlist.TrackCount, playlist.TrackCount)).
				IconName("music-queue-symbolic").
				ConnectActivated(func(_ adw.ActionRow) {
					if err := localdb.AddTrackToPlaylist(playlist.ID, trackID, coverURL); err != nil {
						notifications.OnToast.Notify(gettext.Get("Failed to add to playlist"))
						return
					}
					notifications.OnToast.Notify(fmt.Sprintf(gettext.Get("Added to \"%s\""), playlist.Name))
					closeDialog()
				})()
			row.SetActivatable(true)
			existingGroup.Add(&row.Widget)
		}

		container.Append(&existingGroup.Widget)
	}

	container.Append(&newGroup.Widget)

	clamp := adw.NewClamp()
	clamp.SetMaximumSize(400)
	clamp.SetChild(&container.Widget)

	dialog.SetExtraChild(&clamp.Widget)
	dialog.Present(&w.Widget)
}

// entryRowEntryActivated is a global callback stub for EntryRow's "entry-activated" signal.
// It dispatches to per-widget handlers registered via callback.HandleCallback.
var entryRowEntryActivated = func(widget adw.EntryRow) {
	callback.CallbackHandler[any](widget.Object, "entry-activated", widget)
}
