package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	port string
}

func NewServer(port string) *Config {
	return &Config{
		port: port,
	}
}

func (s *Config) Run() error {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("ERROR: Failed to load .env file")
    }

    db, err := InitDB(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal("ERROR: Unable to connect to database: ", err)
    }
    defer db.Close()

    wishHandler := WishHandler{conn: db}

	mux := http.NewServeMux()
    mux.HandleFunc("GET /wishes", wishHandler.GetAll)
    mux.HandleFunc("GET /wishes/{id}", wishHandler.GetById)
    mux.HandleFunc("POST /wishes", wishHandler.Save)
    mux.HandleFunc("PUT /wishes/{id}", wishHandler.Update)
    mux.HandleFunc("DELETE /wishes/{id}", wishHandler.Delete)

    server := http.Server{
        Addr: s.port,
        Handler: mux,
    }

    log.Printf("INFO: Server is running at %s", s.port)
    return server.ListenAndServe()
}

