package main

import "testing"

func TestDetectPlatform(t *testing.T) {
	tests := []struct {
		url      string
		wantName string
		wantErr  bool
	}{
		{"https://x.com/user/status/123", "x", false},
		{"https://twitter.com/user/status/123", "x", false},
		{"https://www.threads.net/@user/post/abc", "threads", false},
		{"https://threads.com/@user/post/abc", "threads", false},
		{"https://example.com/foo", "", true},
	}

	for _, tt := range tests {
		p, err := DetectPlatform(tt.url)
		if tt.wantErr {
			if err == nil {
				t.Errorf("DetectPlatform(%q) expected error, got %v", tt.url, p)
			}
			continue
		}
		if err != nil {
			t.Errorf("DetectPlatform(%q) unexpected error: %v", tt.url, err)
			continue
		}
		if p.Name() != tt.wantName {
			t.Errorf("DetectPlatform(%q) = %q, want %q", tt.url, p.Name(), tt.wantName)
		}
	}
}
