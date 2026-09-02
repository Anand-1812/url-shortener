// Package main for starting the backend server
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/internal/db"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/joho/godotenv"
)

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

	// layers
	repository := repository.NewURLRepository(pool)
	service := service.NewURLService(repository)
	handler := handler.NewURLHandler(service)

	// routes
	mux := http.NewServeMux()
	mux := http.HandleFunc("GET /healthcheck", h.HealthCheck)
	mux.HandleFunc("POST /shorten", handler.Shorten)
	mux.HandleFunc("GET /{code}", handler.Redirect)

	// ------ server ------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
