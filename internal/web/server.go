package web

import (
	"fmt"
	"net/http"
)

type Server struct {
	port int
	mux  *http.ServeMux
}

func NewServer(port int) *Server {
	s := &Server{
		port: port,
		mux:  http.NewServeMux(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Server listening on %s\n", addr)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", HandleHome)
	s.mux.HandleFunc("/api/timer/start", HandleStartTimer)
	s.mux.HandleFunc("/api/timer/stop", HandleStopTimer)
	s.mux.HandleFunc("/ws", HandleWebSocket)
}
