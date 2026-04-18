package sentinel_transport_http

import (
	"context"
	"net/http"
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

func (h *ImageryTransport) GetRoutes() []Route {
	return []Route{
		{
			Pattern: "/sentinel/imagery",
			Method:  http.MethodGet,
			Handler: h.GetImagery,
		},
	}
}
