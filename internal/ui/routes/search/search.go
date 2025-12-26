package search

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

var promptView *adw.StatusPage

func PromptView() gtk.Widgetter {
	if promptView == nil {
		promptView = adw.NewStatusPage()
		promptView.SetIconName("loupe-symbolic")
		promptView.SetTitle("Search")
		promptView.SetDescription("Start typing in the search bar to search for songs, artists, albums or playlists.")
	}
	return promptView
}

var loadingView *adw.Clamp

func LoadingView() gtk.Widgetter {
	if loadingView == nil {
		spinner := gtk.NewSpinner()
		spinner.SetSpinning(true)
		spinner.Start()

		loadingView = adw.NewClamp()
		loadingView.SetMaximumSize(50)
		loadingView.SetChild(spinner)
	}
	return loadingView
}
