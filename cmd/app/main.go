package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_middleware "github.com/Fitray/sentinel-service/internal/core/middleware"
	core_server "github.com/Fitray/sentinel-service/internal/core/server"
	sentinel_repository_http "github.com/Fitray/sentinel-service/internal/features/sentinel/repository/http"
	sentinel_service_http "github.com/Fitray/sentinel-service/internal/features/sentinel/service/http"
	sentinel_transport_http "github.com/Fitray/sentinel-service/internal/features/sentinel/transport/http"
)

func main() {
	shutdownCtx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	loggerConfig := core_logger.MustConfig()
	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := logger.Close(); err != nil {
			fmt.Println("failed to close logger file:", err)
		}
	}()

	middlewares := core_middleware.GetMiddlewares(logger)

	httpConfig := core_server.MustConfig()
	httpServer := core_server.NewHTTPServer(httpConfig, logger, middlewares)

	imageryRepository := sentinel_repository_http.NewImageryRepositoryMust()
	imageryService := sentinel_service_http.NewImageryService(&imageryRepository)
	imageryTransport := sentinel_transport_http.NewImageryTransport(&imageryService)
	sentinel_routes := imageryTransport.GetRoutes()

	httpServer.RegisterRoutes(sentinel_routes)

	err = httpServer.Run(shutdownCtx)
	if err != nil {
		logger.Error(err)
	}
}
