package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	core_auth "github.com/Fitray/sentinel-service/internal/core/auth"
	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	core_postgres "github.com/Fitray/sentinel-service/internal/core/postgres"
	core_postgres_pool "github.com/Fitray/sentinel-service/internal/core/postgres/pool"
	core_server "github.com/Fitray/sentinel-service/internal/core/server"
	sentinel_repository "github.com/Fitray/sentinel-service/internal/features/sentinel/repository"
	sentinel_service "github.com/Fitray/sentinel-service/internal/features/sentinel/service"
	sentinel_transport_http "github.com/Fitray/sentinel-service/internal/features/sentinel/transport/http"
	users_repository_postgres "github.com/Fitray/sentinel-service/internal/features/users/repository/postgres"
	users_service "github.com/Fitray/sentinel-service/internal/features/users/service"
	users_transport_http "github.com/Fitray/sentinel-service/internal/features/users/transport/http"
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

	authService := core_auth.NewAuthServiceMust()

	httpConfig := core_server.MustConfig()
	httpServer := core_server.NewHTTPServer(httpConfig, logger)

	postgreConf := core_postgres.NewConfigMust()
	pool, err := core_postgres_pool.NewConnectionPool(postgreConf)
	if err != nil {
		logger.Error(fmt.Errorf("connection pool: %w", err))
		return
	}

	routes := make([]core_server.Route, 0)

	imageryRepository := sentinel_repository.NewImageryRepositoryMust(pool)
	imageryService := sentinel_service.NewImageryService(&imageryRepository)
	imageryTransport := sentinel_transport_http.NewImageryTransport(&imageryService)
	routes = imageryTransport.GetRoutes(routes, logger, authService)

	usersRepository := users_repository_postgres.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransport := users_transport_http.NewUsersTransport(usersService, &authService)
	routes = usersTransport.GetRoutes(routes, logger, authService)

	httpServer.RegisterRoutes(routes)

	err = httpServer.Run(shutdownCtx)
	if err != nil {
		logger.Error(err)
	}
}
