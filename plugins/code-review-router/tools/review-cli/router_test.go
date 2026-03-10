package main

import (
	"testing"
)

func TestDecide(t *testing.T) {
	cfg := DefaultConfig()

	tests := []struct {
		name      string
		analysis  Analysis
		wantRoute Route
	}{
		{
			name: "simple change routes to gemini",
			analysis: Analysis{
				Stats:             DiffStats{FileCount: 2, LineCount: 50, DirCount: 1},
				ComplexityScore:   0,
				SecuritySensitive: false,
			},
			wantRoute: RouteGemini,
		},
		{
			name: "high complexity routes to opencode",
			analysis: Analysis{
				Stats:           DiffStats{FileCount: 5, LineCount: 200, DirCount: 2},
				ComplexityScore: 7,
			},
			wantRoute: RouteOpenCode,
		},
		{
			name: "security sensitive routes to opencode",
			analysis: Analysis{
				Stats:             DiffStats{FileCount: 2, LineCount: 30, DirCount: 1},
				ComplexityScore:   1,
				SecuritySensitive: true,
				SecurityFiles:     []string{"src/auth/login.go"},
			},
			wantRoute: RouteOpenCode,
		},
		{
			name: "many files routes to opencode",
			analysis: Analysis{
				Stats:           DiffStats{FileCount: 25, LineCount: 100, DirCount: 3},
				ComplexityScore: 3,
			},
			wantRoute: RouteOpenCode,
		},
		{
			name: "db migration routes to opencode",
			analysis: Analysis{
				Stats:          DiffStats{FileCount: 3, LineCount: 50, DirCount: 2},
				ComplexityScore: 2,
				HasDBMigration: true,
			},
			wantRoute: RouteOpenCode,
		},
		{
			name: "boundary: complexity exactly at threshold routes to opencode",
			analysis: Analysis{
				Stats:           DiffStats{FileCount: 5, LineCount: 200, DirCount: 2},
				ComplexityScore: 6,
			},
			wantRoute: RouteOpenCode,
		},
		{
			name: "boundary: complexity just below threshold routes to gemini",
			analysis: Analysis{
				Stats:           DiffStats{FileCount: 5, LineCount: 200, DirCount: 2},
				ComplexityScore: 5,
			},
			wantRoute: RouteGemini,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := Decide(tt.analysis, cfg)
			if decision.Route != tt.wantRoute {
				t.Errorf("got route %s, want %s (reasons: %v)", decision.Route, tt.wantRoute, decision.Reasons)
			}
		})
	}
}
