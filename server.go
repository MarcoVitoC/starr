package main

import (
	"log"
	"net/http"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Run() error {
	router := http.NewServeMux()

    router.HandleFunc("GET /", WelcomeHandler)

    server := http.Server{
        Addr: s.port,
        Handler: router,
    }

    log.Printf("Server is running at %s", s.port)
    return server.ListenAndServe()
}