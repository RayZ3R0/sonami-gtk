package router

import "github.com/diamondburned/gotk4-adwaita/pkg/adw"

func notFoundHandler(params Params) *Response {
	notFoundView := adw.NewStatusPage()
	notFoundView.SetTitle("Not found")
	notFoundView.SetDescription("The requested deeplink did not have any available handlers.")
	notFoundView.SetIconName("face-sad-symbolic")

	return &Response{
		PageTitle: "Not Found",
		View:      notFoundView,
	}
}
