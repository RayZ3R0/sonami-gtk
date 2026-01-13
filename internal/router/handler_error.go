package router

import (
	. "codeberg.org/dergs/tonearm/pkg/schwifty/syntax"
)

func errorHandler(err error) *Response {
	return &Response{
		PageTitle: "Internal Error",
		View: StatusPage().
			Title("Internal Error").
			Description("Unfortunately an error occurred while loading this view. Please try again later. If the error persists, please open an issue!\n\nError Message: " + err.Error()).
			IconName("sentiment-dissatisfied-symbolic"),
	}
}
