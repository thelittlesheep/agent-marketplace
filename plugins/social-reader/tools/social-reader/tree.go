package main

import (
	"fmt"
	"strings"
)

// BuildTree converts a flat list of comments (with parentIndex) into a nested tree.
// Comments with parentIndex == -1 are root-level.
// Comments with parentIndex >= 0 become children of the comment at that index.
func BuildTree(flat []*Comment) []*Comment {
	if len(flat) == 0 {
		return nil
	}

	var roots []*Comment
	for i, c := range flat {
		if c.parentIndex < 0 || c.parentIndex >= len(flat) || c.parentIndex == i {
			roots = append(roots, c)
		} else {
			parent := flat[c.parentIndex]
			parent.Children = append(parent.Children, c)
		}
	}
	return roots
}

// RenderTree renders a comment tree with box-drawing characters.
func RenderTree(comments []*Comment) string {
	var sb strings.Builder
	for i, c := range comments {
		isLast := i == len(comments)-1
		renderNode(&sb, c, "", isLast, true)
	}
	return sb.String()
}

func renderNode(sb *strings.Builder, c *Comment, prefix string, isLast bool, isRoot bool) {
	// Write the connector
	if !isRoot {
		if isLast {
			sb.WriteString(prefix + "└── ")
		} else {
			sb.WriteString(prefix + "├── ")
		}
	}

	// Author line
	author := c.Author
	if c.IsAuthor {
		author += " 🧵"
	}
	sb.WriteString(fmt.Sprintf("**@%s**: %s\n", author, c.Content))

	// Engagement line
	engStr := formatEngagement(c.Engagement)
	if engStr != "" {
		if isRoot {
			sb.WriteString(engStr + "\n")
		} else {
			childPrefix := prefix
			if isLast {
				childPrefix += "    "
			} else {
				childPrefix += "│   "
			}
			sb.WriteString(childPrefix + engStr + "\n")
		}
	}

	// Children
	newPrefix := prefix
	if !isRoot {
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
	}

	for i, child := range c.Children {
		childIsLast := i == len(c.Children)-1
		renderNode(sb, child, newPrefix, childIsLast, false)
	}
}

func formatEngagement(e *Engagement) string {
	if e == nil {
		return ""
	}
	var parts []string
	if e.Likes != nil {
		parts = append(parts, fmt.Sprintf("❤️ %d", *e.Likes))
	}
	if e.Replies != nil {
		parts = append(parts, fmt.Sprintf("💬 %d", *e.Replies))
	}
	if e.Reposts != nil {
		parts = append(parts, fmt.Sprintf("🔁 %d", *e.Reposts))
	}
	if e.Quotes != nil {
		parts = append(parts, fmt.Sprintf("🔄 %d", *e.Quotes))
	}
	return strings.Join(parts, " · ")
}
