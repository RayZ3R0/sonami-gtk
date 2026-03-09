package router

import (
	"github.com/RayZ3R0/sonami-gtk/internal/gettext"
	. "github.com/RayZ3R0/sonami-gtk/pkg/schwifty/syntax"
)

func notFoundHandler() *Response {
	return &Response{
		PageTitle: gettext.Get("Not Found"),
		View: StatusPage().
			Title(gettext.Get("Not found")).
			Description(gettext.Get("The requested deeplink did not have any available handlers.")).
			IconName("sentiment-dissatisfied-symbolic"),
	}
}
