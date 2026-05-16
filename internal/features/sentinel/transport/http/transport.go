package sentinel_transport_http

import (
	"context"
	"net/http"

	core_auth "github.com/Fitray/sentinel-service/internal/core/auth"
	core_domain "github.com/Fitray/sentinel-service/internal/core/domain"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_middleware "github.com/Fitray/sentinel-service/internal/core/middleware"
	core_server "github.com/Fitray/sentinel-service/internal/core/server"
)

type ImageryTransport struct {
	imageryService ImageryService
}

type ImageryService interface {
	GetImagery(
		ctx context.Context, imageryRequest core_domain.ImageryRequest,
	) (core_domain.ImageryResponce, error)
	GetHistory(
		requestFilter core_domain.FilterRequest,
	) ([]core_domain.NewImagery, error)
	GetRequestFromID(
		user_id string, idStr string,
	) (core_domain.NewImagery, error)
}

// type Route struct {
// 	Pattern string
// 	Method  string
// 	Handler http.HandlerFunc
// }

func NewImageryTransport(
	imageryService ImageryService,
) ImageryTransport {
	return ImageryTransport{
		imageryService: imageryService,
	}
}

func (h *ImageryTransport) GetRoutes(
	routes []core_server.Route,
	logger core_logger.Logger,
	auth core_auth.AuthService,
) []core_server.Route {
	for _, v := range []core_server.Route{
		{
			Pattern:     "/sentinel/imagery",
			Method:      http.MethodPost,
			Handler:     h.GetImagery,
			Version:     "v1",
			Middlewares: core_middleware.AuthGroupMiddlewares(logger, auth),
		},
		{
			Pattern:     "/sentinel/requests",
			Method:      http.MethodGet,
			Handler:     h.GetHistory,
			Version:     "v1",
			Middlewares: core_middleware.AuthGroupMiddlewares(logger, auth),
		},
		{
			Pattern:     "/sentinel/requests/{id}",
			Method:      http.MethodGet,
			Handler:     h.GetRequestFromID,
			Version:     "v1",
			Middlewares: core_middleware.AuthGroupMiddlewares(logger, auth),
		},
	} {
		routes = append(routes, v)
	}
	return routes
}
