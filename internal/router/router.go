package router

import (
	"errors"
	"runtime"

	"github.com/diamondburned/gotk4/pkg/core/glib"
)

type Handler func(params Params) *Response

type Params map[string]any

var routeMap = map[string]Handler{}

func NavigateTo(path string, params Params) {
	OnNavigate.Notify(path)
	history.Push(HistoryEntry{Path: path, Params: params})
	if handler, ok := routeMap[path]; !ok {
		NavigationComplete.Notify(&Response{
			PageTitle: "Not Found",
			View:      getNotFoundView(),
		})
	} else {
		go asyncRouteHandler(handler, params)
	}
}

func asyncRouteHandler(handler Handler, params Params) {
	var response *Response
	var errorResponse *Response
	response = func() *Response {
		defer func() {
			if err := recover(); err != nil {
				errorResponse = &Response{Error: err.(error)}
			}
		}()
		return handler(params)
	}()
	if errorResponse != nil {
		response = errorResponse
	}
	if response == nil {
		response = &Response{Error: errors.New("handler did not generate any response")}
	}
	if response.Error != nil {
		response.View = getErrorView(response.Error)
	}

	glib.IdleAdd(func() bool {
		NavigationComplete.Notify(response)
		runtime.GC()
		return false
	})
}

func Register(path string, handler Handler) {
	routeMap[path] = handler
}
