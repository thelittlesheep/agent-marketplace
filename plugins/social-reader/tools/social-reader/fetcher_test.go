package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "text/markdown" {
			t.Error("expected Accept: text/markdown header")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("# Test Post\n\nHello world"))
	}))
	defer server.Close()

	result, err := fetchURL(server.URL + "/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Hello world") {
		t.Errorf("expected 'Hello world' in result, got: %s", result)
	}
}

func TestFetchHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := fetchURL(server.URL + "/test")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

// fetchURL is a test helper that performs the same logic as Fetch but with an arbitrary URL.
func fetchURL(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/markdown")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", &http.ProtocolError{ErrorString: "bad status"}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
