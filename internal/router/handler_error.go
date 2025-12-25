package router

import "github.com/diamondburned/gotk4-adwaita/pkg/adw"

func errorHandler(err error) *Response {
	errorView := adw.NewStatusPage()
	errorView.SetTitle("Internal Error")
	errorView.SetDescription("Unfortunately an error occurred while loading this view. Please try again later. If the error persists, please open an issue!\n\nError Message: " + err.Error())
	errorView.SetIconName("face-sad-symbolic")

	return &Response{
		PageTitle: "Internal Error",
		View:      errorView,
	}
}
