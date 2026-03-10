package main

import (
	"strings"
	"testing"
)

func TestDetectDomain(t *testing.T) {
	tests := []struct {
		query string
		want  string
	}{
		{"glassmorphism dark mode", "style"},
		{"color palette hex", "color"},
		{"chart visualization trend", "chart"},
		{"landing page hero cta", "landing"},
		{"saas fintech dashboard", "product"},
		{"font typography serif", "typography"},
		{"accessibility wcag mobile", "ux"},
		{"prompt css tailwind", "prompt"},
		{"random unknown query", "style"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			got := DetectDomain(tt.query)
			if got != tt.want {
				t.Errorf("DetectDomain(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestSearchDomain(t *testing.T) {
	result, err := Search("glassmorphism", "style", 3)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if result.Domain != "style" {
		t.Errorf("expected domain=style, got %s", result.Domain)
	}
	if result.Count == 0 {
		t.Error("expected results for glassmorphism in style domain")
	}
	// First result should contain "glassmorphism" somewhere
	if result.Count > 0 {
		found := false
		for _, v := range result.Results[0] {
			if containsCI(v, "glassmorphism") {
				found = true
				break
			}
		}
		if !found {
			t.Error("first result should contain 'glassmorphism'")
		}
	}
}

func TestSearchStack(t *testing.T) {
	result, err := SearchStack("state useState hook", "react", 3)
	if err != nil {
		t.Fatalf("SearchStack failed: %v", err)
	}
	if result.Stack != "react" {
		t.Errorf("expected stack=react, got %s", result.Stack)
	}
	if result.Count == 0 {
		t.Error("expected results for state useState in react stack")
	}
}

func TestSearchAutoDetect(t *testing.T) {
	result, err := Search("saas dashboard", "", 3)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if result.Domain != "product" {
		t.Errorf("expected auto-detect domain=product, got %s", result.Domain)
	}
}

func TestSearchInvalidStack(t *testing.T) {
	_, err := SearchStack("test", "nonexistent", 3)
	if err == nil {
		t.Error("expected error for invalid stack")
	}
}

func containsCI(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
