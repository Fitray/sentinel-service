package sentinel_transport_http

import "net/http"

type Route struct {
	Pattern string
	Method  string
	Handler http.HandlerFunc
}

func GetRoutes() []Route {
	return []Route{
		{
			Pattern: "/sentinel",
			Method:  http.MethodGet,
			Handler: GetImagery,
		},
	}
}
