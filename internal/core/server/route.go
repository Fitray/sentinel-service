package core_server

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

type Route struct {
	Method      string
	Pattern     string
	Handler     http.HandlerFunc
	Version     string
	Middlewares chi.Middlewares
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
