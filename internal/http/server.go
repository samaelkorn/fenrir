package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	Address string
	server  *http.Server
}

func NewServer(host, port string) *Server {
	return &Server{
		Address: net.JoinHostPort(host, port),
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	spa := spaHandler{staticPath: "client/build", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	s.server = &http.Server{
		Addr:         s.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := s.server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "start server error")
	}

	select {

	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return errors.New("server is nil")
	}
	if err := s.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "stop server error")
	}
	return nil
}
