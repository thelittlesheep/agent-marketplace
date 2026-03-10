package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JSONOutput struct {
	Route           string   `json:"route"`
	Reasons         []string `json:"reasons"`
	ComplexityScore int      `json:"complexity_score"`
	FilesChanged    int      `json:"files_changed"`
	LinesChanged    int      `json:"lines_changed"`
	DirsChanged     int      `json:"dirs_changed"`
	Security        bool     `json:"security_sensitive"`
	SecurityFiles   []string `json:"security_files,omitempty"`
	ReviewOutput    string   `json:"review_output,omitempty"`
	ReviewedBy      string   `json:"reviewed_by,omitempty"`
}

func FormatDryRun(analysis Analysis, decision Decision, asJSON bool) string {
	if asJSON {
		out := JSONOutput{
			Route:           string(decision.Route),
			Reasons:         decision.Reasons,
			ComplexityScore: analysis.ComplexityScore,
			FilesChanged:    analysis.Stats.FileCount,
			LinesChanged:    analysis.Stats.LineCount,
			DirsChanged:     analysis.Stats.DirCount,
			Security:        analysis.SecuritySensitive,
			SecurityFiles:   analysis.SecurityFiles,
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		return string(b)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Routing Decision\n\n"))
	sb.WriteString(fmt.Sprintf("**Route**: %s\n", decision.Route))
	sb.WriteString(fmt.Sprintf("**Reason**: %s\n", strings.Join(decision.Reasons, "; ")))
	sb.WriteString(fmt.Sprintf("**Complexity Score**: %d/10\n", analysis.ComplexityScore))
	sb.WriteString(fmt.Sprintf("**Files Changed**: %d\n", analysis.Stats.FileCount))
	sb.WriteString(fmt.Sprintf("**Lines Changed**: %d\n", analysis.Stats.LineCount))
	sb.WriteString(fmt.Sprintf("**Dirs Changed**: %d\n", analysis.Stats.DirCount))

	if analysis.SecuritySensitive {
		sb.WriteString(fmt.Sprintf("**Security Files**: %s\n", strings.Join(analysis.SecurityFiles, ", ")))
	}

	// Show breakdown
	b := analysis.Breakdown
	sb.WriteString("\n### Score Breakdown\n")
	sb.WriteString(fmt.Sprintf("  File count (>10):     +%d\n", b.FileCount))
	sb.WriteString(fmt.Sprintf("  Line count (>500):    +%d\n", b.LineCount))
	sb.WriteString(fmt.Sprintf("  Test/config files:    +%d\n", b.TestConfig))
	sb.WriteString(fmt.Sprintf("  Database changes:     +%d\n", b.Database))
	sb.WriteString(fmt.Sprintf("  API route changes:    +%d\n", b.APIRoutes))
	sb.WriteString(fmt.Sprintf("  Cross-directory (3+): +%d\n", b.CrossDirectory))

	return sb.String()
}

func FormatResult(analysis Analysis, decision Decision, result ExecResult, asJSON bool) string {
	if asJSON {
		out := JSONOutput{
			Route:           string(decision.Route),
			Reasons:         decision.Reasons,
			ComplexityScore: analysis.ComplexityScore,
			FilesChanged:    analysis.Stats.FileCount,
			LinesChanged:    analysis.Stats.LineCount,
			DirsChanged:     analysis.Stats.DirCount,
			Security:        analysis.SecuritySensitive,
			SecurityFiles:   analysis.SecurityFiles,
			ReviewOutput:    result.Output,
			ReviewedBy:      result.UsedTool,
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		return string(b)
	}

	var sb strings.Builder
	sb.WriteString("## Code Review Results\n\n")
	sb.WriteString(fmt.Sprintf("**Router Decision**: %s CLI\n", decision.Route))
	sb.WriteString(fmt.Sprintf("**Reason**: %s\n", strings.Join(decision.Reasons, "; ")))
	sb.WriteString(fmt.Sprintf("**Complexity Score**: %d/10\n", analysis.ComplexityScore))
	sb.WriteString(fmt.Sprintf("**Files Changed**: %d\n", analysis.Stats.FileCount))
	sb.WriteString(fmt.Sprintf("**Lines Changed**: %d\n", analysis.Stats.LineCount))
	sb.WriteString("\n---\n\n")
	sb.WriteString(strings.TrimSpace(result.Output))
	sb.WriteString("\n\n---\n\n")
	sb.WriteString(fmt.Sprintf("**Review completed by**: %s\n", result.UsedTool))

	return sb.String()
}
