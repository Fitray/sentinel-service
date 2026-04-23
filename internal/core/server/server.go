package core_server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	"github.com/go-chi/chi"
)

type HTTPServer struct {
	Config Config
	Mux    *chi.Mux
	Logger core_logger.Logger
}

func NewHTTPServer(
	config Config,
	log core_logger.Logger,
) *HTTPServer {
	return &HTTPServer{
		Config: config,
		Mux:    chi.NewRouter(),
		Logger: log,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := http.Server{
		Addr:    s.Config.Addr,
		Handler: s.Mux,
	}

	ch := make(chan error, 1)

	go func() {
		s.Logger.Debug("Starting HTTP server", slog.String("Addr", s.Config.Addr))
		defer close(ch)
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				ch <- err
			}
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("an error occured while HTTP server was running: %w", err)
		}
	case <-ctx.Done():
		shutdownTimeout, cancel := context.WithTimeout(
			context.Background(),
			s.Config.Timeout,
		)
		defer cancel()

		s.Logger.Debug("HTTP server shutdown")

		if err := server.Shutdown(shutdownTimeout); err != nil {
			server.Close()
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
	}
	return nil
}

func (s *HTTPServer) RegisterRoutes(routes []Route, middlewares chi.Middlewares) {
	s.Mux.Group(func(r chi.Router) {
		r.Use(middlewares...)
		for _, route := range routes {
			pattern, err := url.JoinPath("/", "api", route.Version, route.Pattern)
			if err != nil {
				s.Logger.Warn(
					fmt.Sprintf("failed to register route: %s", err.Error()),
					slog.String("pattern", route.Pattern),
				)
			}
			r.MethodFunc(route.Method, pattern, route.Handler)
		}
	})
}
