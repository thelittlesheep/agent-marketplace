package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

const maxValueLen = 300

func FormatSearchResult(result SearchResult, asJSON bool) string {
	if asJSON {
		b, _ := json.MarshalIndent(result, "", "  ")
		return string(b)
	}

	var sb strings.Builder

	if result.Stack != "" {
		sb.WriteString(fmt.Sprintf("## Stack: %s | Query: \"%s\"\n\n", result.Stack, result.Query))
	} else {
		sb.WriteString(fmt.Sprintf("## Domain: %s | Query: \"%s\"\n\n", result.Domain, result.Query))
	}

	if result.Count == 0 {
		sb.WriteString("No results found.\n")
		return sb.String()
	}

	for i, entry := range result.Results {
		sb.WriteString(fmt.Sprintf("### Result %d\n", i+1))

		// Use output cols order from config to maintain consistent ordering
		var orderedCols []string
		if result.Stack != "" {
			orderedCols = StackOutputCols
		} else if cfg, ok := DomainConfigs[result.Domain]; ok {
			orderedCols = cfg.OutputCols
		}

		if orderedCols != nil {
			for _, col := range orderedCols {
				val, ok := entry[col]
				if ok && val != "" {
					sb.WriteString(fmt.Sprintf("- **%s:** %s\n", col, truncate(val)))
				}
			}
		} else {
			for k, v := range entry {
				if v != "" {
					sb.WriteString(fmt.Sprintf("- **%s:** %s\n", k, truncate(v)))
				}
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func truncate(s string) string {
	if utf8.RuneCountInString(s) <= maxValueLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxValueLen]) + "..."
}
