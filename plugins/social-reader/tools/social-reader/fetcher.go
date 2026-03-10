package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const jinaBaseURL = "https://r.jina.ai/"

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// Fetch retrieves the markdown content of a URL via Jina Reader.
func Fetch(normalizedURL string) (string, error) {
	jinaURL := jinaBaseURL + normalizedURL

	req, err := http.NewRequest("GET", jinaURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "text/markdown")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		// Retry once after a short delay
		time.Sleep(3 * time.Second)
		resp2, err := httpClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("retry failed: %w", err)
		}
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			return "", fmt.Errorf("rate limited (HTTP %d)", resp2.StatusCode)
		}
		resp = resp2
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from Jina Reader", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	return string(body), nil
}
