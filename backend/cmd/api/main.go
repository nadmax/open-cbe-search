package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nadmax/open-cbe-search/core/engine/postgres"
)

func main() {
	conn := os.Getenv("POSTGRES_URL")
	if conn == "" {
		log.Fatalf("No Postgres URL defined")
	}

	db, err := postgres.NewClient(conn)
	if err != nil {
		log.Fatalf("Error connecting to Postgres: %s", err)
	}
	defer db.Close()

	api := NewAPI(db)

	http.HandleFunc("/api/search", api.SearchHandler)
	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}