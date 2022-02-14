package internalhttp

import (
	"context"
	"log"
	"net"
	"net/http"
	"path"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	Address string
	server  *http.Server
}

func (s *Server) NewServer(host, port string) *Server {
	return &Server{
		Address: net.JoinHostPort(host, port),
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	router.HandleFunc("/", s.ServeHTTP)

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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(path.Join("client", "build", "index.html"))

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusMethodNotAllowed)
		log.Print("Internal Server Error")
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusMethodNotAllowed)
		log.Print("Internal Server Error")
		return
	}
}
