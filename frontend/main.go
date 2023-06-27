package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// This should match the homeResponse struct from the backend service.
type backendHomeResponse struct {
	hostname  string
	randomNum int
}

func main() {
	// Register endpoint handlers.
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHomePage)

	// Use PORT environment variable, or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the web server on port and accept requests
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	// Output frontend details.
	host, _ := os.Hostname()
	fmt.Fprintf(w, "This is the frontend!\n")
	fmt.Fprintf(w, "Version: 1.0.0\n")
	fmt.Fprintf(w, "Hostname: %s\n", host)
	// Make GET request to backend.
	backendUrl := os.Getenv("BACKEND_URL")
	if backendUrl == "" {
		backendUrl = "http://my-backend-service:80"
	}
	backendResponse, err := http.Get(backendUrl)
	if err != nil {
		fmt.Fprintf(w, "Failed to reach backend.\n")
		return
	}
	// Read response body from backend.
	backendResponseBodyRaw, err := ioutil.ReadAll(backendResponse.Body)
	if err != nil {
		fmt.Fprintf(w, "Failed to read backend response body.\n")
		return
	}
	// Parse and output response body from backend.
	var backendResponseParsed backendHomeResponse
	json.Unmarshal([]byte(backendResponseBodyRaw), &backendResponse)
	fmt.Fprintf(w, "Response from backend: hostname: %s\n", backendResponseParsed.hostname)
	fmt.Fprintf(w, "Response from backend: randomNum: %d\n", backendResponseParsed.randomNum)
}
