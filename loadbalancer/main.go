package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type backend struct {
	URL   *url.URL
	alive bool
	mu    sync.RWMutex
}

var (
	current  int
	mu       sync.Mutex
	backends []*backend
)

var healthCheckInterval = flag.Duration("healthcheck-interval", 10*time.Second, "Health check interval")

func newBackend(rawURL string) *backend {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Invalid backend URL: %s, error %v", rawURL, err)
	}
	return &backend{
		URL:   parsedURL,
		alive: true,
	}
}

func healthCheckLoop(backends []*backend, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		for _, b := range backends {
			resp, err := http.Get(b.URL.String() + "/health")
			b.mu.Lock()
			b.alive = err == nil && resp.StatusCode == http.StatusOK
			b.mu.Unlock()
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}

func getNextAliveBackend() *backend {
	mu.Lock()
	defer mu.Unlock()
	n := len(backends)
	for i := 0; i < n; i++ {
		b := backends[current%n]
		current++
		b.mu.RLock()
		alive := b.alive
		b.mu.RUnlock()
		if alive {
			return b
		}
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nReceived request from %s\n", r.RemoteAddr)
	fmt.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Proto)
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}
	fmt.Println() // Blank line for readability

	backend := getNextAliveBackend()
	if backend == nil {
		http.Error(w, "No healthy backend available", http.StatusServiceUnavailable)
		return
	}

	targetURL := *backend.URL
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

	body, _ := io.ReadAll(resp.Body)
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(body)

	fmt.Printf("Forwarded request to %s with status code %d\n", backend.URL, resp.StatusCode)
}

func main() {
	backendList := flag.String("backends", "", "Comma-separated list of backend URLs")

	flag.Parse()

	if *backendList == "" {
		log.Fatal("You must provide at least one backend using the -backends flag")
	}

	for _, u := range strings.Split(*backendList, ",") {
		backends = append(backends, newBackend(strings.TrimSpace(u)))
	}

	go healthCheckLoop(backends, *healthCheckInterval)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
