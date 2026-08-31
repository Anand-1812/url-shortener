package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Server struct {
	db *pgxpool.Pool
}

// method for the handler
func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is from %s", r.URL.Path[1:])
}

func main() {
	_ = godotenv.Load()

	dbConnString := os.Getenv("NEON_DB_STRING")
	if dbConnString == "" {
		log.Fatal("DB env is required")
	}

	pool, err := pgxpool.New(context.Background(), dbConnString)
	if err != nil {
		log.Fatal("Error connecting to db")
	}
	defer pool.Close()

	fmt.Println("Connection established")

	server := &Server{db: pool}

	http.HandleFunc("/", server.handler)
	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
