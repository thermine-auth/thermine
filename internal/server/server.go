// Package server wires up the HTTP server for the auth service.
package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"xermess/internal/config"
)

// Server serves the auth API.
type Server struct {
	cfg    config.Config
	logger *slog.Logger
}

// New returns a Server ready to Run.
func New(cfg config.Config, logger *slog.Logger) *Server {
	return &Server{cfg: cfg, logger: logger}
}

// Run starts the HTTP server and blocks until ctx is cancelled, then shuts
// down gracefully.
func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: s.routes(),
	}

	errs := make(chan error, 1)
	go func() {
		s.logger.Info("listening", "addr", s.cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- err
			return
		}
		errs <- nil
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		s.logger.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	return mux
}
