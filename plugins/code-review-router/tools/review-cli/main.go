package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var (
		dryRun    bool
		force     string
		staged    bool
		commit    string
		rangeFlag string
		asJSON    bool
	)

	rootCmd := &cobra.Command{
		Use:   "review-cli",
		Short: "Intelligently route code reviews to Gemini or OpenCode CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Mutual exclusion: --staged, --commit, --range
			modeCount := 0
			if staged {
				modeCount++
			}
			if commit != "" {
				modeCount++
			}
			if rangeFlag != "" {
				modeCount++
			}
			if modeCount > 1 {
				return fmt.Errorf("--staged, --commit, and --range are mutually exclusive")
			}

			mode := DiffMode{
				Staged: staged,
				Commit: commit,
				Range:  rangeFlag,
			}

			cfg := LoadConfig()

			// Verify git repo
			if _, err := gitCmd("rev-parse", "--is-inside-work-tree"); err != nil {
				return fmt.Errorf("not inside a git repository")
			}

			// Check at least one CLI is available
			hasGemini := commandExists("gemini")
			hasOpenCode := commandExists("opencode")
			if !hasGemini && !hasOpenCode {
				return fmt.Errorf("neither gemini nor opencode CLI found in PATH")
			}

			// Analyze
			analysis, err := Analyze(cfg, mode)
			if err != nil {
				return fmt.Errorf("analysis failed: %w", err)
			}

			if analysis.Stats.FileCount == 0 {
				fmt.Println("No changes to review.")
				return nil
			}

			// Route
			var decision Decision
			if force != "" {
				decision = Decision{
					Route:   Route(force),
					Reasons: []string{"forced by --force flag"},
				}
			} else {
				decision = Decide(analysis, cfg)
			}

			// Adjust if forced tool not available
			if decision.Route == RouteGemini && !hasGemini {
				decision = Decision{Route: RouteOpenCode, Reasons: []string{"gemini not available, using opencode"}}
			} else if decision.Route == RouteOpenCode && !hasOpenCode {
				decision = Decision{Route: RouteGemini, Reasons: []string{"opencode not available, using gemini"}}
			}

			// Dry run
			if dryRun {
				fmt.Println(FormatDryRun(analysis, decision, asJSON))
				return nil
			}

			// Execute
			result := Execute(decision, cfg, mode)
			if result.Err != nil {
				return fmt.Errorf("review failed: %w", result.Err)
			}

			fmt.Println(FormatResult(analysis, decision, result, asJSON))
			return nil
		},
	}

	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show routing decision without executing")
	rootCmd.Flags().StringVar(&force, "force", "", "Force a specific CLI (gemini or opencode)")
	rootCmd.Flags().BoolVar(&staged, "staged", false, "Review staged changes only")
	rootCmd.Flags().StringVar(&commit, "commit", "", "Review a single commit's changes")
	rootCmd.Flags().StringVar(&rangeFlag, "range", "", "Review a commit range (e.g., abc123..def456)")
	rootCmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
