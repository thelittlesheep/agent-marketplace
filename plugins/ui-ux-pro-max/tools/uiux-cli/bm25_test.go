package main

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  int // expected token count
	}{
		{"hello world", 0},           // both <= 2 chars after filtering? no, "hello"=5, "world"=5
		{"a b c", 0},                 // all single chars
		{"glassmorphism dark mode", 2}, // "glassmorphism", "dark"=4, "mode"=4
		{"UI/UX design", 3},          // "ui"=2, "ux"=2, "design"=6
		{"", 0},
	}

	// Fix: "hello world" has 2 tokens (hello=5, world=5)
	tests[0].want = 2
	// "glassmorphism dark mode" = glassmorphism(13) + dark(4) + mode(4) = 3 tokens
	tests[2].want = 3

	for _, tt := range tests {
		tokens := Tokenize(tt.input)
		if len(tokens) != tt.want {
			t.Errorf("Tokenize(%q) = %v (len=%d), want len=%d", tt.input, tokens, len(tokens), tt.want)
		}
	}
}

func TestTokenizeFiltersShortWords(t *testing.T) {
	tokens := Tokenize("a UI is ok for UX")
	for _, tok := range tokens {
		if len([]rune(tok)) <= 1 {
			t.Errorf("token %q should have been filtered (len <= 1)", tok)
		}
	}
}

func TestBM25Score(t *testing.T) {
	docs := []string{
		"glassmorphism transparent blur effect modern design",
		"minimalism clean simple whitespace swiss style",
		"brutalism raw bold industrial design heavy",
		"glassmorphism frosted glass blur transparency",
	}

	bm25 := NewBM25(docs)
	results := bm25.Score("glassmorphism blur")

	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}

	// Doc 0 and 3 should be top results (both contain glassmorphism + blur)
	topIdx := results[0].Index
	if topIdx != 0 && topIdx != 3 {
		t.Errorf("expected top result to be doc 0 or 3, got %d", topIdx)
	}

	// Doc 1 (minimalism) should not appear (score = 0)
	for _, r := range results {
		if r.Index == 1 {
			t.Errorf("doc 1 (minimalism) should not match glassmorphism query, score=%f", r.Score)
		}
	}
}

func TestBM25EmptyCorpus(t *testing.T) {
	bm25 := NewBM25([]string{})
	results := bm25.Score("test")
	if len(results) != 0 {
		t.Errorf("expected no results for empty corpus, got %d", len(results))
	}
}
