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

	yaml, err := openYAML(fileVar)
	if err != nil {
		log.Fatalf("Could not open YAML file '%s': %v", fileVar, err)
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mux)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func openYAML(file string) ([]byte, error) {
	data, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func init() {
	flag.StringVar(&fileVar, "file", "input.yaml", "YAML file for path to URL mappings")
	flag.Parse()
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
