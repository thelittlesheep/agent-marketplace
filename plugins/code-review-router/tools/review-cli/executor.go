package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ExecResult struct {
	Output   string
	UsedTool string // "gemini", "opencode", "opencode (fallback: ...)"
	Err      error
}

func Execute(decision Decision, cfg Config, mode DiffMode) ExecResult {
	diff, err := getDiff(mode)
	if err != nil {
		return ExecResult{Err: fmt.Errorf("failed to get diff: %w", err)}
	}

	switch decision.Route {
	case RouteGemini:
		result := execGemini(diff, cfg)
		if result.Err != nil {
			// Fallback: try opencode chain
			fmt.Fprintln(os.Stderr, "gemini failed, trying opencode...")
			return execOpenCodeChain(diff, cfg)
		}
		return result

	case RouteOpenCode:
		result := execOpenCodeChain(diff, cfg)
		if result.Err != nil {
			// Fallback: try gemini
			fmt.Fprintln(os.Stderr, "opencode chain failed, trying gemini...")
			return execGemini(diff, cfg)
		}
		return result

	default:
		return ExecResult{Err: fmt.Errorf("unknown route: %s", decision.Route)}
	}
}

func getDiff(mode DiffMode) (string, error) {
	out, err := gitCmd(mode.diffArgs()...)
	if err != nil {
		return "", err
	}
	return out, nil
}

func execGemini(diff string, cfg Config) ExecResult {
	if !commandExists("gemini") {
		return ExecResult{Err: fmt.Errorf("gemini CLI not found")}
	}

	cmd := exec.Command("gemini", "-m", cfg.Models.Gemini, "-p", cfg.Prompts.Gemini)
	cmd.Stdin = strings.NewReader(diff)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return ExecResult{Err: fmt.Errorf("gemini failed: %w", err)}
	}

	return ExecResult{
		Output:   string(out),
		UsedTool: "gemini",
	}
}

func execOpenCodeChain(diff string, cfg Config) ExecResult {
	if !commandExists("opencode") {
		return ExecResult{Err: fmt.Errorf("opencode CLI not found")}
	}

	// Write diff to temp file
	tmpFile, err := os.CreateTemp("", "review-cli-*.diff")
	if err != nil {
		return ExecResult{Err: fmt.Errorf("failed to create temp file: %w", err)}
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(diff); err != nil {
		tmpFile.Close()
		return ExecResult{Err: fmt.Errorf("failed to write diff: %w", err)}
	}
	tmpFile.Close()

	// Try primary model
	result := execOpenCode(tmpFile.Name(), cfg.Models.OpenCodePrimary, "high", cfg)
	if result.Err == nil {
		return result
	}

	// Fallback
	fmt.Fprintf(os.Stderr, "opencode primary (%s) failed, trying fallback (%s)...\n",
		cfg.Models.OpenCodePrimary, cfg.Models.OpenCodeFallback)

	result = execOpenCode(tmpFile.Name(), cfg.Models.OpenCodeFallback, "", cfg)
	if result.Err == nil {
		result.UsedTool = fmt.Sprintf("opencode (fallback: %s)", cfg.Models.OpenCodeFallback)
		return result
	}

	return ExecResult{Err: fmt.Errorf("opencode chain failed")}
}

func execOpenCode(diffFile, model, variant string, cfg Config) ExecResult {
	args := []string{"run", "--model", model}
	if variant != "" {
		args = append(args, "--variant", variant)
	}
	args = append(args, cfg.Prompts.OpenCode, "-f", diffFile)

	cmd := exec.Command("opencode", args...)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return ExecResult{Err: err}
	}

	return ExecResult{
		Output:   string(out),
		UsedTool: "opencode",
	}
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
