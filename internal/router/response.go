package router

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type Response struct {
	Error     error
	PageTitle string
	Toolbar   gtk.Widgetter
	View      gtk.Widgetter
}

func FromError(pageTitle string, err error) *Response {
	return &Response{
		Error:     err,
		PageTitle: pageTitle,
		Toolbar:   nil,
		View:      nil,
	}
}
