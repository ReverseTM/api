package server

import (
	"context"
	"errors"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(
	addr string,
	handler http.Handler) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &Server{
		httpServer: srv,
	}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
