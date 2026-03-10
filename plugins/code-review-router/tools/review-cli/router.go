package main

import "strconv"

type Route string

const (
	RouteGemini   Route = "gemini"
	RouteOpenCode Route = "opencode"
)

type Decision struct {
	Route   Route
	Reasons []string
}

func Decide(analysis Analysis, cfg Config) Decision {
	var reasons []string

	if analysis.SecuritySensitive {
		reasons = append(reasons, "security-sensitive files detected")
	}
	if analysis.Stats.FileCount >= cfg.Thresholds.FileCount {
		reasons = append(reasons, "high file count ("+itoa(analysis.Stats.FileCount)+")")
	}
	if analysis.HasDBMigration {
		reasons = append(reasons, "database migration detected")
	}
	if analysis.ComplexityScore >= cfg.Thresholds.Complexity {
		reasons = append(reasons, "complexity score "+itoa(analysis.ComplexityScore)+"/10 >= threshold "+itoa(cfg.Thresholds.Complexity))
	}

	if len(reasons) > 0 {
		return Decision{Route: RouteOpenCode, Reasons: reasons}
	}

	return Decision{
		Route:   RouteGemini,
		Reasons: []string{"low complexity (" + itoa(analysis.ComplexityScore) + "/10), fast review"},
	}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
