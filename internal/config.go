package internal

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

	WishService := WishService{conn: db}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /wishes", WishService.GetAll)
	mux.HandleFunc("GET /wishes/{id}", WishService.GetById)
	mux.HandleFunc("POST /wishes", WishService.Save)
	mux.HandleFunc("PUT /wishes/{id}", WishService.Update)
	mux.HandleFunc("DELETE /wishes/{id}", WishService.Delete)

	server := http.Server{
		Addr:    s.port,
		Handler: mux,
	}

	log.Printf("INFO: Server is running at %s", s.port)
	return server.ListenAndServe()
}
