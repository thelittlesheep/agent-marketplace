package main

import (
	"strings"
	"testing"
)

func TestXMatchURL(t *testing.T) {
	x := &XPlatform{}
	tests := []struct {
		url  string
		want bool
	}{
		{"https://x.com/user/status/123", true},
		{"https://twitter.com/user/status/123", true},
		{"https://www.x.com/user/status/123", true},
		{"http://twitter.com/user/status/123", true},
		{"https://threads.net/@user/post/abc", false},
		{"https://example.com/foo", false},
	}

	for _, tt := range tests {
		if got := x.MatchURL(tt.url); got != tt.want {
			t.Errorf("MatchURL(%q) = %v, want %v", tt.url, got, tt.want)
		}
	}
}

func TestXNormalizeURL(t *testing.T) {
	x := &XPlatform{}
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{
			"https://twitter.com/elonmusk/status/123",
			"https://x.com/elonmusk/status/123",
			false,
		},
		{
			"http://www.x.com/user/status/456?s=20",
			"https://x.com/user/status/456",
			false,
		},
		{
			"https://x.com/user/article/789",
			"https://x.com/user/article/789",
			false,
		},
		{
			"https://x.com/user",
			"",
			true, // profile URL not supported
		},
	}

	for _, tt := range tests {
		got, err := x.NormalizeURL(tt.input)
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

func TestXParsePost(t *testing.T) {
	x := &XPlatform{}
	markdown := `# @elonmusk

The future of AI is incredible.

Check out this new development.

![Image](https://pbs.twimg.com/media/test.jpg)

❤️ 42000 · 💬 1500 · 🔁 8000

Post your reply
`

	post, err := x.ParsePost(markdown, "https://x.com/elonmusk/status/123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.Author != "elonmusk" {
		t.Errorf("Author = %q, want %q", post.Author, "elonmusk")
	}
	if post.Platform != "x" {
		t.Errorf("Platform = %q, want %q", post.Platform, "x")
	}
	if post.SourceURL != "https://x.com/elonmusk/status/123" {
		t.Errorf("SourceURL = %q", post.SourceURL)
	}
	if len(post.Images) != 1 {
		t.Errorf("expected 1 image, got %d", len(post.Images))
	}
	if post.Engagement == nil {
		t.Fatal("expected engagement metrics")
	}
	if post.Engagement.Likes == nil || *post.Engagement.Likes != 42000 {
		t.Errorf("expected 42000 likes")
	}
}

func TestXParsePost_RealJinaOutput(t *testing.T) {
	x := &XPlatform{}

	// Realistic Jina Reader output matching actual X page structure
	markdown := `Don't miss what's happening

People on X are the first to know.

[Log in](https://x.com/login)

[Sign up](https://x.com/signup)

Conversation
============

[![Image 1](https://pbs.twimg.com/profile_images/123/avatar_normal.jpg)](https://x.com/pangyusio)

[Pangyu 胖鱼 ![Image 2: 🐠](https://abs-0.twimg.com/emoji/v2/svg/1f420.svg)](https://x.com/pangyusio)

[@pangyusio](https://x.com/pangyusio)

同时用Claude code 和 codex，还真让我有了老板思维了。

看来我老板说的是真的。

Translate post

[3:52 AM · Feb 27, 2026](https://x.com/pangyusio/status/2027230253068980675)

·

[12.4K Views](https://x.com/pangyusio/status/2027230253068980675/analytics)

18
4
65
45

Read 18 replies

New to X?
---------

Sign up now to get your own personalized timeline!

Sign up with Apple

[Create account](https://x.com/i/flow/signup)
`

	post, err := x.ParsePost(markdown, "https://x.com/pangyusio/status/123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Author should be extracted
	if post.Author != "pangyusio" {
		t.Errorf("Author = %q, want %q", post.Author, "pangyusio")
	}

	// Content should NOT include noise
	noiseStrings := []string{
		"Translate post", "3:52 AM", "12.4K Views",
		"Read 18 replies", "New to X?", "Sign up",
		"Create account", "Don't miss",
	}
	for _, noise := range noiseStrings {
		if containsStr(post.Content, noise) {
			t.Errorf("Content should not contain %q, got:\n%s", noise, post.Content)
		}
	}

	// Content SHOULD include the actual tweet text
	if !containsStr(post.Content, "老板思维") {
		t.Errorf("Content should contain tweet text, got:\n%s", post.Content)
	}

	// Images should NOT include emoji SVG
	for _, img := range post.Images {
		if containsStr(img, "twimg.com/emoji") {
			t.Errorf("Images should not contain emoji SVG: %s", img)
		}
	}

	// Engagement should be parsed from bare numbers
	if post.Engagement == nil {
		t.Fatal("expected engagement metrics")
	}
	if post.Engagement.Replies == nil || *post.Engagement.Replies != 18 {
		t.Errorf("expected 18 replies, got %v", post.Engagement.Replies)
	}
	if post.Engagement.Reposts == nil || *post.Engagement.Reposts != 4 {
		t.Errorf("expected 4 reposts, got %v", post.Engagement.Reposts)
	}
	if post.Engagement.Likes == nil || *post.Engagement.Likes != 65 {
		t.Errorf("expected 65 likes, got %v", post.Engagement.Likes)
	}
}

func containsStr(s, substr string) bool {
	return strings.Contains(s, substr)
}


func TestIsXArticleURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://x.com/user/article/123", true},
		{"https://x.com/user/article/abc-slug", true},
		{"https://x.com/user/status/123", false},
		{"https://x.com/user", false},
	}
	for _, tt := range tests {
		if got := isXArticleURL(tt.url); got != tt.want {
			t.Errorf("isXArticleURL(%q) = %v, want %v", tt.url, got, tt.want)
		}
	}
}

func TestXParsePost_ArticleJinaOutput(t *testing.T) {
	x := &XPlatform{}
	markdown := `Title: The Future of Open Source AI
URL Source: https://x.com/techwriter/article/1234567890
Markdown Content:

Don't miss what's happening

People on X are the first to know.

[Log in](https://x.com/login)

[Sign up](https://x.com/signup)

The Future of Open Source AI
============================

Open source AI has transformed the landscape of machine learning research.

In the past year, we have seen unprecedented collaboration between
institutions and individual contributors.

## Key Developments

1. Model weights becoming standard releases
2. Training recipes shared openly
3. Evaluation frameworks standardized

The implications for the industry are profound and far-reaching.

New to X?
---------

Sign up now to get your own personalized timeline!

[Create account](https://x.com/i/flow/signup)`

	post, err := x.ParsePost(markdown, "https://x.com/techwriter/article/1234567890")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Author from URL fallback
	if post.Author != "techwriter" {
		t.Errorf("Author = %q, want %q", post.Author, "techwriter")
	}

	// Content should include title
	if !containsStr(post.Content, "# The Future of Open Source AI") {
		t.Errorf("Content should contain article title, got:\n%s", post.Content)
	}

	// Content should include body
	if !containsStr(post.Content, "Open source AI has transformed") {
		t.Errorf("Content should contain article body, got:\n%s", post.Content)
	}

	// Content should preserve sub-headings
	if !containsStr(post.Content, "## Key Developments") {
		t.Errorf("Content should contain sub-headings, got:\n%s", post.Content)
	}

	// Content should NOT include noise
	noiseStrings := []string{
		"Don't miss what", "Log in", "Sign up",
		"New to X?", "Create account",
		"Title:", "URL Source:", "Markdown Content:",
	}
	for _, noise := range noiseStrings {
		if containsStr(post.Content, noise) {
			t.Errorf("Content should not contain %q, got:\n%s", noise, post.Content)
		}
	}
}

func TestXParsePost_ArticleNoJinaHeaders(t *testing.T) {
	x := &XPlatform{}
	markdown := `Don't miss what's happening

People on X are the first to know.

[Log in](https://x.com/login)

# My Article Title

This is the article body without Jina metadata headers.

It has multiple paragraphs of content.

New to X?
---------

Sign up now`

	post, err := x.ParsePost(markdown, "https://x.com/writer/article/999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !containsStr(post.Content, "# My Article Title") {
		t.Errorf("Content should contain ATX title, got:\n%s", post.Content)
	}

	if !containsStr(post.Content, "article body without Jina") {
		t.Errorf("Content should contain body, got:\n%s", post.Content)
	}

	if containsStr(post.Content, "Don't miss") || containsStr(post.Content, "New to X") {
		t.Errorf("Content should not contain noise, got:\n%s", post.Content)
	}
}

func TestXParseComments_ReturnsError(t *testing.T) {
	x := &XPlatform{}
	comments, err := x.ParseComments("any markdown")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if comments != nil {
		t.Errorf("expected nil comments, got %v", comments)
	}
	if !strings.Contains(err.Error(), "does not support --comments") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"42", 42},
		{"1,234", 1234},
		{"1,234,567", 1234567},
		{"0", 0},
	}
	for _, tt := range tests {
		if got := parseNumber(tt.input); got != tt.want {
			t.Errorf("parseNumber(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
