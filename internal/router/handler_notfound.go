package router

import (
	. "codeberg.org/dergs/tidalwave/pkg/schwifty/syntax"
)

func notFoundHandler() *Response {
	return &Response{
		PageTitle: "Not Found",
		View: StatusPage().
			Title("Not found").
			Description("The requested deeplink did not have any available handlers.").
			IconName("face-sad-symbolic"),
	}
}
