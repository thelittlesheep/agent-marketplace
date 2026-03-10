package main

import (
	"strings"
	"testing"
)

func TestFormatMarkdown_PostOnly(t *testing.T) {
	result := &Result{
		Post: &Post{
			Author:    "testuser",
			Content:   "Hello world!",
			Images:    []string{"https://example.com/img.jpg"},
			SourceURL: "https://x.com/testuser/status/123",
			Platform:  "x",
			Engagement: &Engagement{
				Likes:   intPtr(100),
				Replies: intPtr(10),
			},
		},
	}

	output := FormatMarkdown(result)

	if !strings.Contains(output, "## @testuser") {
		t.Error("missing author header")
	}
	if !strings.Contains(output, "Hello world!") {
		t.Error("missing content")
	}
	if !strings.Contains(output, "![](https://example.com/img.jpg)") {
		t.Error("missing image")
	}
	if !strings.Contains(output, "❤️ 100") {
		t.Error("missing likes")
	}
	if !strings.Contains(output, "**Source**: https://x.com/testuser/status/123") {
		t.Error("missing source URL")
	}
	if strings.Contains(output, "### Comments") {
		t.Error("should not show comments section when no comments")
	}
}

func TestFormatMarkdown_WithComments(t *testing.T) {
	result := &Result{
		Post: &Post{
			Author:    "alice",
			Content:   "Test post",
			SourceURL: "https://x.com/alice/status/123",
			Platform:  "x",
		},
		Comments: []*Comment{
			{
				Author:     "bob",
				Content:    "Nice!",
				Engagement: &Engagement{Likes: intPtr(5)},
			},
		},
	}

	output := FormatMarkdown(result)

	if !strings.Contains(output, "### Comments") {
		t.Error("missing comments section")
	}
	if !strings.Contains(output, "**@bob**: Nice!") {
		t.Error("missing comment content")
	}
}

func TestFormatMarkdown_QuotedPost(t *testing.T) {
	result := &Result{
		Post: &Post{
			Author:  "alice",
			Content: "Interesting take",
			QuotedPost: &Post{
				Author:  "bob",
				Content: "Original thought",
			},
			SourceURL: "https://x.com/alice/status/123",
			Platform:  "x",
		},
	}

	output := FormatMarkdown(result)

	if !strings.Contains(output, "> **@bob**: Original thought") {
		t.Error("missing quoted post")
	}
}
