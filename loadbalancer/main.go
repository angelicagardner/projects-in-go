package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request from %s\n", r.RemoteAddr)
	fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}
	fmt.Println() // Blank line for readability

	backendURL := "http://localhost:9090"

	targetURL, err := url.Parse(backendURL)
	if err != nil {
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}
	targetURL.Path = r.URL.Path
	targetURL.RawQuery = r.URL.RawQuery

	req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response from server: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading backend response", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(body))

	w.WriteHeader(resp.StatusCode)
}

func main() {
	port := "8080"

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
