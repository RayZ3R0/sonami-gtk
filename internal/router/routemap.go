package router

var routeMap = map[string]Handler{}

func Register(path string, handler Handler) {
	routeMap[path] = handler
}
