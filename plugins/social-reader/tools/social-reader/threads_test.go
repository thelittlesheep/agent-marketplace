package main

import (
	"testing"
)

func TestThreadsMatchURL(t *testing.T) {
	th := &ThreadsPlatform{}
	tests := []struct {
		url  string
		want bool
	}{
		{"https://www.threads.net/@zuck/post/abc123", true},
		{"https://threads.net/@user/post/xyz", true},
		{"https://threads.com/@user/post/xyz", true},
		{"https://www.threads.net/t/abc123", true},
		{"https://x.com/user/status/123", false},
		{"https://example.com/foo", false},
	}

	for _, tt := range tests {
		if got := th.MatchURL(tt.url); got != tt.want {
			t.Errorf("MatchURL(%q) = %v, want %v", tt.url, got, tt.want)
		}
	}
}

func TestThreadsNormalizeURL(t *testing.T) {
	th := &ThreadsPlatform{}
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{
			"https://threads.com/@zuck/post/abc123",
			"https://www.threads.net/@zuck/post/abc123",
			false,
		},
		{
			"http://threads.net/@user/post/xyz?hl=en",
			"https://www.threads.net/@user/post/xyz",
			false,
		},
		{
			"https://www.threads.net/t/abc123",
			"https://www.threads.net/t/abc123",
			false,
		},
		{
			"https://www.threads.net/@user",
			"",
			true, // profile URL not supported
		},
	}

	for _, tt := range tests {
		got, err := th.NormalizeURL(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("NormalizeURL(%q) expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("NormalizeURL(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("NormalizeURL(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestThreadsParsePost(t *testing.T) {
	th := &ThreadsPlatform{}
	markdown := `[![zuck's profile picture](https://scontent-xxx.cdninstagram.com/v/s150x150/avatar.jpg)](/@zuck)

@zuck

This is an amazing update to Threads!

![Image 1: No photo description available](https://scontent-xxx.cdninstagram.com/v/media.jpg)

❤️ 15000 · 💬 2000 · 🔁 500 · 🔄 100

Translate
`

	post, err := th.ParsePost(markdown, "https://www.threads.net/@zuck/post/abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.Author != "zuck" {
		t.Errorf("Author = %q, want %q", post.Author, "zuck")
	}
	if post.Platform != "threads" {
		t.Errorf("Platform = %q, want %q", post.Platform, "threads")
	}
	if post.Engagement == nil {
		t.Fatal("expected engagement metrics")
	}
	if post.Engagement.Likes == nil || *post.Engagement.Likes != 15000 {
		t.Errorf("expected 15000 likes, got %v", post.Engagement.Likes)
	}
}

func TestThreadsParseComments(t *testing.T) {
	th := &ThreadsPlatform{}
	markdown := `@zuck

Main post content here

❤️ 1000 · 💬 50

@alice
Great post!

❤️ 10 · 💬 1

@bob ·Author
Thanks everyone!

❤️ 5 · 💬 0
`

	comments, err := th.ParseComments(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(comments) < 2 {
		t.Fatalf("expected at least 2 comments, got %d", len(comments))
	}

	// All comments should be root-level (conservative nesting)
	for i, c := range comments {
		if c.parentIndex != -1 {
			t.Errorf("comment %d (%s) has parentIndex %d, want -1 (flat)", i, c.Author, c.parentIndex)
		}
	}

	// Find the ·Author comment
	foundAuthor := false
	for _, c := range comments {
		if c.IsAuthor {
			foundAuthor = true
			if c.Author != "bob" {
				t.Errorf("expected author comment from bob, got %s", c.Author)
			}
		}
	}
	if !foundAuthor {
		t.Error("expected to find a comment with IsAuthor=true")
	}
}

func TestStripThreadsNoise(t *testing.T) {
	input := `[![user's profile picture](https://scontent-xxx.cdninstagram.com/s150x150/pic.jpg)](/@user)

@user

Hello world

Translate

Some content`

	result := stripThreadsNoise(input)

	if containsString(result, "profile picture") {
		t.Error("should have stripped profile picture link")
	}
	if containsString(result, "Translate") {
		t.Error("should have stripped Translate line")
	}
	if !containsString(result, "Hello world") {
		t.Error("should have kept content")
	}
}

func containsString(haystack, needle string) bool {
	return len(haystack) > 0 && len(needle) > 0 && indexOf(haystack, needle) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
