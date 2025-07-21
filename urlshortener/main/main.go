package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"urlshortener"
	"urlshortener/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var fileVar string

const (
	dbUrl = "postgresql://user:password@localhost:5432/urlshortener_db"
)

func main() {
	pool, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		log.Fatalf("Could not connect to database. %v", err)
	}

	defer pool.Close()

	urlRepo := repository.NewPostgresRepository(pool)

	mux := defaultMux()

	dbHandler := urlshort.DBHandler(urlRepo, mux)

	// data, err := os.ReadFile(fileVar)

	if err != nil {
		log.Fatalf("Could not open file '%s', '%v'", fileVar, err)
	}

	// yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mux)
	// jsonHandler, err := urlshort.JSONHandler(data, mux)
	// if err != nil {
	// 	log.Fatalf("Error parsing YAML: %v", err)
	// }

	fmt.Println("Starting the server on :8080")

	http.ListenAndServe(":8080", dbHandler)
	// http.ListenAndServe(":8080", jsonHandler)
	// http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func init() {
	flag.StringVar(&fileVar, "file", "input.json", "JSON file for path to URL mappings")
	flag.Parse()
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
