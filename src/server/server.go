package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/senfix/logger"
)

type Server struct {
	http.Server
	logger      logger.Log
	shutdownReq chan bool
	running     bool
}

func NewServer(listen string) Server {
	return Server{
		Server: http.Server{
			Addr:         listen,
			ReadTimeout:  10 * time.Minute,
			WriteTimeout: 10 * time.Minute,
		},
		shutdownReq: make(chan bool),
	}
}

func (s *Server) Stop() {
	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Shutdown request error: %v", err)
	}
}

func (s *Server) Start(controller ...Controller) (err error) {

	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	handler := mux.NewRouter()
	for _, c := range controller {
		c.Register(handler)
	}
	s.Handler = handlers.CORS(originsOk, headersOk, methodsOk)(handler)
	s.running = true
	return s.ListenAndServe()
}
