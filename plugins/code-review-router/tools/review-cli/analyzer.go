package main

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type DiffMode struct {
	Staged bool
	Commit string // single commit hash
	Range  string // commit range "from..to"
}

func (m DiffMode) diffArgs() []string {
	switch {
	case m.Commit != "":
		return []string{"show", "--format=", m.Commit}
	case m.Range != "":
		return []string{"diff", m.Range}
	case m.Staged:
		return []string{"diff", "HEAD"}
	default:
		return []string{"diff"}
	}
}

func (m DiffMode) statsArgs(flag string) []string {
	base := m.diffArgs()
	return append(base, flag)
}

type DiffStats struct {
	FileCount    int
	LineCount    int
	DirCount     int
	ChangedFiles []string
}

type ComplexityBreakdown struct {
	FileCount      int
	LineCount      int
	TestConfig     int
	Database       int
	APIRoutes      int
	CrossDirectory int
}

type Analysis struct {
	Stats             DiffStats
	ComplexityScore   int
	Breakdown         ComplexityBreakdown
	SecuritySensitive bool
	SecurityFiles     []string
	HasDBMigration    bool
}

func Analyze(cfg Config, mode DiffMode) (Analysis, error) {
	stats, err := collectDiffStats(mode)
	if err != nil {
		return Analysis{}, err
	}

	breakdown := calculateComplexity(stats, cfg)
	score := breakdown.FileCount + breakdown.LineCount + breakdown.TestConfig +
		breakdown.Database + breakdown.APIRoutes + breakdown.CrossDirectory

	securityFiles := detectSecurity(stats.ChangedFiles)
	hasDB := breakdown.Database > 0

	return Analysis{
		Stats:             stats,
		ComplexityScore:   score,
		Breakdown:         breakdown,
		SecuritySensitive: len(securityFiles) > 0,
		SecurityFiles:     securityFiles,
		HasDBMigration:    hasDB,
	}, nil
}

func collectDiffStats(mode DiffMode) (DiffStats, error) {
	// Get changed files
	out, err := gitCmd(mode.statsArgs("--name-only")...)
	if err != nil {
		return DiffStats{}, err
	}

	files := splitLines(out)

	// Get line counts
	out, err = gitCmd(mode.statsArgs("--numstat")...)
	if err != nil {
		return DiffStats{}, err
	}

	lineCount := 0
	for _, line := range splitLines(out) {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			add, _ := strconv.Atoi(parts[0])
			del, _ := strconv.Atoi(parts[1])
			lineCount += add + del
		}
	}

	// Get unique directories
	dirs := make(map[string]bool)
	for _, f := range files {
		dirs[filepath.Dir(f)] = true
	}

	return DiffStats{
		FileCount:    len(files),
		LineCount:    lineCount,
		DirCount:     len(dirs),
		ChangedFiles: files,
	}, nil
}

func calculateComplexity(stats DiffStats, cfg Config) ComplexityBreakdown {
	var b ComplexityBreakdown

	if stats.FileCount > 10 {
		b.FileCount = 2
	}
	if stats.LineCount > cfg.Thresholds.LineCount {
		b.LineCount = 2
	}
	if stats.DirCount >= cfg.Thresholds.DirCount {
		b.CrossDirectory = 1
	}

	for _, f := range stats.ChangedFiles {
		lower := strings.ToLower(f)

		// Test/config
		if isTestOrConfig(lower) {
			b.TestConfig = 1
		}
		// Database
		if isDatabase(lower) {
			b.Database = 2
		}
		// API routes
		if isAPIRoute(lower) {
			b.APIRoutes = 2
		}
	}

	return b
}

func detectSecurity(files []string) []string {
	var matched []string
	for _, f := range files {
		if isSecuritySensitive(f) {
			matched = append(matched, f)
		}
	}
	return matched
}

func isTestOrConfig(f string) bool {
	patterns := []string{
		"test", "_test.", ".test.", "spec.",
		"config.", "configuration.", ".config",
		"jest.config", "vitest.config", "tsconfig",
		"webpack.config", "vite.config",
	}
	for _, p := range patterns {
		if strings.Contains(f, p) {
			return true
		}
	}
	return false
}

func isDatabase(f string) bool {
	patterns := []string{
		"migration", "migrate",
		"schema", "seed",
		"prisma/", "drizzle/",
		".sql",
	}
	for _, p := range patterns {
		if strings.Contains(f, p) {
			return true
		}
	}
	return false
}

func isAPIRoute(f string) bool {
	patterns := []string{
		"/api/", "/routes/", "/route/",
		"/endpoints/", "/endpoint/",
		"/handlers/", "/handler/",
		"/controllers/", "/controller/",
		"router.", "routes.",
	}
	for _, p := range patterns {
		if strings.Contains(f, p) {
			return true
		}
	}
	return false
}

func isSecuritySensitive(f string) bool {
	lower := strings.ToLower(f)
	patterns := []string{
		"/auth/", "/authentication/", "/authorization/",
		"credential", "secret", "token",
		".env", "config/secrets",
		"/security/", "/crypto/",
		"/middleware/auth",
	}
	for _, p := range patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func gitCmd(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func splitLines(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
