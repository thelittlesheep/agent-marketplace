package main

import (
	"testing"
)

func TestCalculateComplexity(t *testing.T) {
	cfg := DefaultConfig()

	tests := []struct {
		name     string
		stats    DiffStats
		wantMin  int
		wantMax  int
	}{
		{
			name: "simple change",
			stats: DiffStats{
				FileCount:    2,
				LineCount:    50,
				DirCount:     1,
				ChangedFiles: []string{"src/main.go", "src/utils.go"},
			},
			wantMin: 0,
			wantMax: 0,
		},
		{
			name: "many files",
			stats: DiffStats{
				FileCount:    15,
				LineCount:    50,
				DirCount:     1,
				ChangedFiles: []string{"a.go", "b.go", "c.go"},
			},
			wantMin: 2, // file count +2
			wantMax: 2,
		},
		{
			name: "many lines",
			stats: DiffStats{
				FileCount:    2,
				LineCount:    600,
				DirCount:     1,
				ChangedFiles: []string{"src/big.go"},
			},
			wantMin: 2, // line count +2
			wantMax: 2,
		},
		{
			name: "database migration",
			stats: DiffStats{
				FileCount:    3,
				LineCount:    100,
				DirCount:     2,
				ChangedFiles: []string{"db/migration/001_init.sql", "src/model.go"},
			},
			wantMin: 2, // database +2
			wantMax: 3, // database +2, could be +1 for test/config overlap
		},
		{
			name: "api routes",
			stats: DiffStats{
				FileCount:    3,
				LineCount:    100,
				DirCount:     2,
				ChangedFiles: []string{"src/api/users.go", "src/api/auth.go"},
			},
			wantMin: 2, // api +2
			wantMax: 2,
		},
		{
			name: "cross directory",
			stats: DiffStats{
				FileCount:    3,
				LineCount:    50,
				DirCount:     4,
				ChangedFiles: []string{"a/x.go", "b/y.go", "c/z.go", "d/w.go"},
			},
			wantMin: 1, // cross-dir +1
			wantMax: 1,
		},
		{
			name: "complex change",
			stats: DiffStats{
				FileCount:    12,
				LineCount:    600,
				DirCount:     4,
				ChangedFiles: []string{
					"src/api/handler.go",
					"db/migration/002.sql",
					"src/test/api_test.go",
				},
			},
			wantMin: 8, // file(2) + line(2) + test(1) + db(2) + api(2) + dir(1) = 10, but min 8
			wantMax: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := calculateComplexity(tt.stats, cfg)
			score := b.FileCount + b.LineCount + b.TestConfig + b.Database + b.APIRoutes + b.CrossDirectory
			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("score = %d, want [%d, %d], breakdown: %+v", score, tt.wantMin, tt.wantMax, b)
			}
		})
	}
}

func TestDetectSecurity(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  bool
	}{
		{
			name:  "no security files",
			files: []string{"src/main.go", "README.md"},
			want:  false,
		},
		{
			name:  "auth directory",
			files: []string{"src/auth/login.go"},
			want:  true,
		},
		{
			name:  "env file",
			files: []string{".env.production"},
			want:  true,
		},
		{
			name:  "credentials",
			files: []string{"src/credential_store.go"},
			want:  true,
		},
		{
			name:  "crypto",
			files: []string{"pkg/crypto/hash.go"},
			want:  true,
		},
		{
			name:  "token in name",
			files: []string{"src/token_manager.go"},
			want:  true,
		},
		{
			name:  "middleware auth",
			files: []string{"src/middleware/auth_check.go"},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := detectSecurity(tt.files)
			got := len(matched) > 0
			if got != tt.want {
				t.Errorf("got %v, want %v (matched: %v)", got, tt.want, matched)
			}
		})
	}
}

func TestIsTestOrConfig(t *testing.T) {
	tests := []struct {
		file string
		want bool
	}{
		{"src/main_test.go", true},
		{"src/main.go", false},
		{"jest.config.js", true},
		{"vite.config.ts", true},
		{"tsconfig.json", true},
		{"src/utils.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			if got := isTestOrConfig(tt.file); got != tt.want {
				t.Errorf("isTestOrConfig(%q) = %v, want %v", tt.file, got, tt.want)
			}
		})
	}
}
