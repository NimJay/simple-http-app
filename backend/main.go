package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type homeResponse struct {
	// Heads-up: json.Marshal() ignores camelCase fields, so we use PascalCase.
	Hostname  string
	RandomNum int
}

func main() {
	// Register endpoint handlers.
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHome)

	// Use PORT environment variable, or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the web server on port and accept requests
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Create the response
	var response homeResponse
	hostname, _ := os.Hostname()
	response.Hostname = hostname
	response.RandomNum = rand.Intn(100)

	// Convert response to JSON and write response
	json.NewEncoder(w).Encode(response)
}
