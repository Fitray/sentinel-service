package core_server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	core_logger "github.com/Fitray/sentinel-service/internal/core/logger"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5500",
			"http://127.0.0.1:5500",
		},

		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},

		ExposedHeaders: []string{
			"Link",
		},

		AllowCredentials: false,

		MaxAge: 300,
	}))
	return &HTTPServer{
		Config: config,
		Mux:    r,
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

func (s *HTTPServer) RegisterRoutes(routes []Route) {
	for _, route := range routes {
		route := route

		s.Mux.Group(func(r chi.Router) {
			r.Use(route.Middlewares...)

			pattern := "/api/" + route.Version + route.Pattern

			s.Logger.Debug(
				"route registered",
				slog.String("pattern", pattern),
			)

			r.MethodFunc(route.Method, pattern, route.Handler)
		})
	}
}

func (s *HTTPServer) RegisterWeb() {
	s.Mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/pages/index.html")
	})

	s.Mux.Handle("/*", http.FileServer(http.Dir("./web")))
}
