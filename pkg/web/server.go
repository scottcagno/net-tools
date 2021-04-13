package web

import (
	"net/http"
	"os"
	"time"
)

func init() {
	defaultServer = &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 * MB,
		ErrorLog:       NewLoggerWithPrefix(os.Stderr, "[HTTP Server] "),
	}
}

var defaultServer *http.Server

type Server struct {
	*http.Server
}

func NewServer(s *http.Server) *Server {
	if s == nil {
		return &Server{defaultServer}
	}
	return &Server{s}
}

func (s *Server) WithAddr(addr string) *Server {
	s.Addr = addr
	return s
}

func (s *Server) WithHandler(handler http.Handler) *Server {
	s.Handler = handler
	return s
}

func (s *Server) ListenAndServe() error {
	return s.ListenAndServe()
}

func ListenAndServe(addr string, handler http.Handler) error {
	server := &Server{defaultServer}
	server.Addr = addr
	server.Handler = handler
	return server.ListenAndServe()
}
