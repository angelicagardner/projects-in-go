package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// server mock for testing
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Main handler for application responses
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("Hello From Backend Server on port %s", port)
		if _, err := w.Write([]byte(msg)); err != nil {
			fmt.Printf("Error writing response: %v", err)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			fmt.Printf("Error writing health check response: %v", err)
		}
	})

	fmt.Printf("Backend server running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
