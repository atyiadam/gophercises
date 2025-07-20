package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"urlshortener"
)

var fileVar string

func main() {
	mux := defaultMux()

	data, err := os.ReadFile(fileVar)

	if err != nil {
		log.Fatalf("Could not open file '%s', '%v'", fileVar, err)
	}

	// yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mux)
	jsonHandler, err := urlshort.JSONHandler(data, mux)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	fmt.Println("Starting the server on :8080")

	http.ListenAndServe(":8080", jsonHandler)
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
