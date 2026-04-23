package sentinel_transport_http

import (
	"context"
	"net/http"

	core_server "github.com/Fitray/sentinel-service/internal/core/server"
)

type ImageryTransport struct {
	imageryService ImageryService
}

type ImageryService interface {
	GetImagery(ctx context.Context, city, from, to string) ([]byte, error)
}

type Route struct {
	Pattern string
	Method  string
	Handler http.HandlerFunc
}

func NewImageryTransport(
	imageryService ImageryService,
) ImageryTransport {
	return ImageryTransport{
		imageryService: imageryService,
	}
}

func (h *ImageryTransport) GetRoutes(
	routes []core_server.Route,
) []core_server.Route {
	for _, v := range []core_server.Route{
		{
			Pattern: "/sentinel/imagery",
			Method:  http.MethodGet,
			Handler: h.GetImagery,
			Version: "v1",
		},
	} {
		routes = append(routes, v)
	}
	return routes
}
