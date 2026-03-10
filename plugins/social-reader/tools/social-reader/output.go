package main

import (
	"fmt"
	"strings"
)

// FormatMarkdown renders a Result as human-readable markdown.
func FormatMarkdown(r *Result) string {
	var sb strings.Builder

	// Post
	p := r.Post
	sb.WriteString(fmt.Sprintf("## @%s\n\n", p.Author))
	sb.WriteString(p.Content + "\n")

	if p.QuotedPost != nil {
		sb.WriteString(fmt.Sprintf("\n> **@%s**: %s\n", p.QuotedPost.Author, p.QuotedPost.Content))
	}

	if len(p.Images) > 0 {
		sb.WriteString("\n**Images**:\n")
		for _, img := range p.Images {
			sb.WriteString(fmt.Sprintf("![](%s)\n", img))
		}
	}

	if p.Engagement != nil {
		engStr := formatEngagement(p.Engagement)
		if engStr != "" {
			sb.WriteString("\n" + engStr + "\n")
		}
	}

	sb.WriteString(fmt.Sprintf("**Source**: %s\n", p.SourceURL))

	// Comments
	if len(r.Comments) > 0 {
		sb.WriteString("\n---\n\n### Comments\n\n")
		sb.WriteString(RenderTree(r.Comments))
	}

	return sb.String()
}
