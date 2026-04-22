package core_server

import "net/http"

type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
	Version string
}

// func NewRoute(
// 	method string,
// 	pattern string,
// 	handler http.HandlerFunc,
// ) Route {
// 	return Route{
// 		Method:  method,
// 		Pattern: pattern,
// 		Handler: handler,
// 	}
// }
