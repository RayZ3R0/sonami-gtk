package router

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

var (
	loadingView  *adw.Clamp
	notFoundView *adw.StatusPage
)

func getLoadingView() gtk.Widgetter {
	if loadingView != nil {
		return loadingView
	}

	spinner := gtk.NewSpinner()
	spinner.SetSpinning(true)
	spinner.Start()

	loadingView = adw.NewClamp()
	loadingView.SetMaximumSize(50)
	loadingView.SetChild(spinner)
	return loadingView
}

func getNotFoundView() gtk.Widgetter {
	if notFoundView != nil {
		return notFoundView
	}

	notFoundView = adw.NewStatusPage()
	notFoundView.SetTitle("Not found")
	notFoundView.SetDescription("The requested deeplink did not have any available handlers.")
	notFoundView.SetIconName("face-sad-symbolic")

	return notFoundView
}

func getErrorView(err error) gtk.Widgetter {
	notFoundView = adw.NewStatusPage()
	notFoundView.SetTitle("Internal Error")
	notFoundView.SetDescription("Unfortunately an error occurred while loading this view. Please try again later. If the error persists, please open an issue!\n\nError Message: " + err.Error())
	notFoundView.SetIconName("face-sad-symbolic")

	return notFoundView
}
