package router

import (
	"codeberg.org/dergs/tidalwave/pkg/schwifty"
)

type Response struct {
	Error     error
	PageTitle string
	Toolbar   schwifty.BaseWidgetable
	View      schwifty.BaseWidgetable
}

func FromError(pageTitle string, err error) *Response {
	return &Response{
		Error:     err,
		PageTitle: pageTitle,
		Toolbar:   nil,
		View:      nil,
	}
}
