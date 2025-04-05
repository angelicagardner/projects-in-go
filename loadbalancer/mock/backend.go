package main

import (
	"fmt"
	"log"
	"net/http"
)

func backendHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Hello From Backend Server")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// mock a server for testing
func main() {
	port := "9090"
	http.HandleFunc("/", backendHandler)
	fmt.Printf("Backend server running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
