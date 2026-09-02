// Package main for starting the backend server
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/internal/db"

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
	ctx := context.Background()

	dbConnString := os.Getenv("NEON_DB_STRING")
	if dbConnString == "" {
		log.Fatal("DB env is required")
	}

	pool, err := db.ConnectDB(ctx, dbConnString)
	if err != nil {
		log.Fatalf("Unable to connect with databse: %v\n", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	fmt.Println("DB Connection established")

	if err := db.InitSchema(ctx, pool); err != nil {
		log.Fatalf("DB initialization failed: %v", err)
	}

	// ------ server ------
	server := &Server{db: pool}

	http.HandleFunc("/", server.handler)
	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
