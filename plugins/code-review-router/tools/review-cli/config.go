package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ModelsConfig struct {
	Gemini          string `toml:"gemini"`
	OpenCodePrimary string `toml:"opencode_primary"`
	OpenCodeFallback string `toml:"opencode_fallback"`
}

type ThresholdsConfig struct {
	Complexity int `toml:"complexity"`
	FileCount  int `toml:"file_count"`
	LineCount  int `toml:"line_count"`
	DirCount   int `toml:"dir_count"`
}

type PromptsConfig struct {
	Gemini   string `toml:"gemini"`
	OpenCode string `toml:"opencode"`
}

type Config struct {
	Models     ModelsConfig     `toml:"models"`
	Thresholds ThresholdsConfig `toml:"thresholds"`
	Prompts    PromptsConfig    `toml:"prompts"`
}

func DefaultConfig() Config {
	return Config{
		Models: ModelsConfig{
			Gemini:           "gemini-3-pro-preview",
			OpenCodePrimary:  "openai/gpt-5.3-codex",
			OpenCodeFallback: "kimi/kimi-k2.5-free",
		},
		Thresholds: ThresholdsConfig{
			Complexity: 6,
			FileCount:  20,
			LineCount:  500,
			DirCount:   3,
		},
		Prompts: PromptsConfig{
			Gemini: `Review this code diff. Focus on:
1. Code quality and best practices
2. Potential bugs or issues
3. Suggestions for improvement

Be concise and actionable.`,
			OpenCode: `Perform a thorough code review. Analyze:
1. Security vulnerabilities
2. Architectural concerns
3. Performance implications
4. Error handling
5. Edge cases

Provide detailed feedback with severity levels.`,
		},
	}
}

func LoadConfig() Config {
	cfg := DefaultConfig()

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}

	configPath := filepath.Join(home, ".config", "review-cli", "config.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg
	}

	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return DefaultConfig()
	}

	// Fill in any zero-value fields with defaults
	defaults := DefaultConfig()
	if cfg.Models.Gemini == "" {
		cfg.Models.Gemini = defaults.Models.Gemini
	}
	if cfg.Models.OpenCodePrimary == "" {
		cfg.Models.OpenCodePrimary = defaults.Models.OpenCodePrimary
	}
	if cfg.Models.OpenCodeFallback == "" {
		cfg.Models.OpenCodeFallback = defaults.Models.OpenCodeFallback
	}
	if cfg.Thresholds.Complexity == 0 {
		cfg.Thresholds.Complexity = defaults.Thresholds.Complexity
	}
	if cfg.Thresholds.FileCount == 0 {
		cfg.Thresholds.FileCount = defaults.Thresholds.FileCount
	}
	if cfg.Thresholds.LineCount == 0 {
		cfg.Thresholds.LineCount = defaults.Thresholds.LineCount
	}
	if cfg.Thresholds.DirCount == 0 {
		cfg.Thresholds.DirCount = defaults.Thresholds.DirCount
	}
	if cfg.Prompts.Gemini == "" {
		cfg.Prompts.Gemini = defaults.Prompts.Gemini
	}
	if cfg.Prompts.OpenCode == "" {
		cfg.Prompts.OpenCode = defaults.Prompts.OpenCode
	}

	return cfg
}
