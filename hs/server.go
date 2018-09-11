package hs

import (
	"context"
	"net/http"
	"whitetrefoil.com/log-go/logger"
)

type Server struct {
	HttpServer *http.Server
	Handler    *Handler
	Ended      chan int
	logger     *logger.Logger
}

func (s *Server) Start() error {
	go func() {
		s.HttpServer.ListenAndServe()
	}()
	s.logger.Logln("HTTP Server has started...")
	return nil
}

func (s *Server) Stop() {
	(*s.Handler).Shutdown()
	s.HttpServer.Shutdown(context.Background())
	s.logger.Logln("HTTP Server is shutting down...")
	s.Ended <- 1
}

func NewServer(address string, handler Handler, l *logger.Logger) *Server {
	httpServer := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	return &Server{
		HttpServer: httpServer,
		Handler:    &handler,
		Ended:      make(chan int),
		logger:     l,
	}
}
