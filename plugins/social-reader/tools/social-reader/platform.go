package main

import "fmt"

var registry []Platform

// RegisterPlatform adds a platform to the global registry.
// Call this in init() of each platform file.
func RegisterPlatform(p Platform) {
	registry = append(registry, p)
}

// DetectPlatform finds the matching platform for a URL.
func DetectPlatform(rawURL string) (Platform, error) {
	for _, p := range registry {
		if p.MatchURL(rawURL) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("unsupported URL: %s", rawURL)
}
