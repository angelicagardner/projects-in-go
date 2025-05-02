package main

import (
	"net/url"
	"testing"
)

func makeBackend(u string, alive bool) *backend {
	parsed, _ := url.Parse(u)
	b := &backend{URL: parsed, alive: alive}
	return b
}

func TestGetNextAliveBackend(t *testing.T) {
	b1 := makeBackend("http://a", true)
	b2 := makeBackend("http://b", false)
	b3 := makeBackend("http://c", true)
	backends = []*backend{b1, b2, b3}
	current = 0

	// First call → b1
	if got := getNextAliveBackend(); got != b1 {
		t.Errorf("1st: got %v, want %v", got.URL, b1.URL)
	}
	// Second → skips b2, returns b3
	if got := getNextAliveBackend(); got != b3 {
		t.Errorf("2nd: got %v, want %v", got.URL, b3.URL)
	}
	// Next → wraps around to b1 again
	if got := getNextAliveBackend(); got != b1 {
		t.Errorf("3rd: got %v, want %v", got.URL, b1.URL)
	}
}
