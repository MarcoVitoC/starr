package main

import (
	"log"
	"net/http"
)

var wishlist = []Wish{}

type Config struct {
	port string
}

func NewServer(port string) *Config {
	return &Config{
		port: port,
	}
}

func (s *Config) Run() error {
	mux := http.NewServeMux()

    mux.HandleFunc("GET /wishes", GetWishesHandler)
    mux.HandleFunc("GET /wishes/{id}", GetWishHandler)
    mux.HandleFunc("POST /wishes", SaveWishHandler)
    mux.HandleFunc("PUT /wishes/{id}", UpdateWishHandler)

    server := http.Server{
        Addr: s.port,
        Handler: mux,
    }

    log.Printf("Server is running at %s", s.port)
    return server.ListenAndServe()
}

