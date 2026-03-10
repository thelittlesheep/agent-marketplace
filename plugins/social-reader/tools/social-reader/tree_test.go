package main

import (
	"strings"
	"testing"
)

func TestBuildTree_RootOnly(t *testing.T) {
	flat := []*Comment{
		{Author: "alice", Content: "Hello", parentIndex: -1},
		{Author: "bob", Content: "World", parentIndex: -1},
	}
	roots := BuildTree(flat)
	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(roots))
	}
	if roots[0].Author != "alice" || roots[1].Author != "bob" {
		t.Error("unexpected root order")
	}
}

func TestBuildTree_Nested(t *testing.T) {
	flat := []*Comment{
		{Author: "alice", Content: "Parent", parentIndex: -1},       // 0
		{Author: "bob", Content: "Reply to alice", parentIndex: 0},  // 1
		{Author: "carol", Content: "Reply to bob", parentIndex: 1},  // 2
		{Author: "dave", Content: "Another root", parentIndex: -1},  // 3
	}
	roots := BuildTree(flat)
	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(roots))
	}

	// alice should have 1 child (bob)
	if len(roots[0].Children) != 1 {
		t.Fatalf("expected alice to have 1 child, got %d", len(roots[0].Children))
	}
	if roots[0].Children[0].Author != "bob" {
		t.Error("expected bob as alice's child")
	}

	// bob should have 1 child (carol)
	if len(roots[0].Children[0].Children) != 1 {
		t.Fatalf("expected bob to have 1 child, got %d", len(roots[0].Children[0].Children))
	}
	if roots[0].Children[0].Children[0].Author != "carol" {
		t.Error("expected carol as bob's child")
	}
}

func TestBuildTree_Empty(t *testing.T) {
	roots := BuildTree(nil)
	if roots != nil {
		t.Error("expected nil for empty input")
	}
}

func TestBuildTree_InvalidParentIndex(t *testing.T) {
	flat := []*Comment{
		{Author: "alice", Content: "Hello", parentIndex: 5}, // invalid index → becomes root
	}
	roots := BuildTree(flat)
	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}
}

func TestRenderTree_Simple(t *testing.T) {
	comments := []*Comment{
		{
			Author:     "user_a",
			Content:    "This is great!",
			Engagement: &Engagement{Likes: intPtr(42), Replies: intPtr(2)},
			Children: []*Comment{
				{
					Author:     "user_b",
					Content:    "Totally agree!",
					Engagement: &Engagement{Likes: intPtr(10), Replies: intPtr(1)},
					Children: []*Comment{
						{
							Author:     "user_a",
							Content:    "Thanks!",
							Engagement: &Engagement{Likes: intPtr(3)},
						},
					},
				},
				{
					Author:     "user_c",
					Content:    "Nice work",
					Engagement: &Engagement{Likes: intPtr(5)},
				},
			},
		},
	}

	output := RenderTree(comments)

	// Check structure
	if !strings.Contains(output, "**@user_a**: This is great!") {
		t.Error("missing root comment")
	}
	if !strings.Contains(output, "├── **@user_b**: Totally agree!") {
		t.Error("missing first child with ├──")
	}
	if !strings.Contains(output, "└── **@user_a**: Thanks!") {
		t.Error("missing nested child with └──")
	}
	if !strings.Contains(output, "└── **@user_c**: Nice work") {
		t.Error("missing last child with └──")
	}
	if !strings.Contains(output, "❤️ 42 · 💬 2") {
		t.Error("missing engagement metrics")
	}
}

func TestRenderTree_IsAuthor(t *testing.T) {
	comments := []*Comment{
		{
			Author:   "zuck",
			Content:  "Original post reply",
			IsAuthor: true,
		},
	}

	output := RenderTree(comments)
	if !strings.Contains(output, "🧵") {
		t.Error("expected author indicator 🧵")
	}
}

func TestFormatEngagement(t *testing.T) {
	tests := []struct {
		name string
		e    *Engagement
		want string
	}{
		{"nil", nil, ""},
		{"likes only", &Engagement{Likes: intPtr(10)}, "❤️ 10"},
		{"all fields", &Engagement{Likes: intPtr(10), Replies: intPtr(5), Reposts: intPtr(3), Quotes: intPtr(1)}, "❤️ 10 · 💬 5 · 🔁 3 · 🔄 1"},
		{"zero values", &Engagement{Likes: intPtr(0)}, "❤️ 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatEngagement(tt.e)
			if got != tt.want {
				t.Errorf("formatEngagement() = %q, want %q", got, tt.want)
			}
		})
	}
}
