package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var (
		showComments bool
		asJSON       bool
	)

	rootCmd := &cobra.Command{
		Use:   "social-reader <url>",
		Short: "Fetch and parse social media posts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rawURL := args[0]

			// 1. Detect platform
			platform, err := DetectPlatform(rawURL)
			if err != nil {
				return err
			}

			// 2. Normalize URL
			normalizedURL, err := platform.NormalizeURL(rawURL)
			if err != nil {
				return err
			}

			// 3. Fetch via Jina Reader
			markdown, err := Fetch(normalizedURL)
			if err != nil {
				return fmt.Errorf("fetch failed: %w", err)
			}

			// 4. Parse post
			post, err := platform.ParsePost(markdown, normalizedURL)
			if err != nil {
				return fmt.Errorf("parse failed: %w", err)
			}

			result := &Result{Post: post}

			// 5. Parse comments (if requested)
			if showComments {
				comments, err := platform.ParseComments(markdown)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
				} else {
					// Build tree from flat list
					result.Comments = BuildTree(comments)
				}
			}

			// 6. Output
			if asJSON {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				enc.SetEscapeHTML(false)
				return enc.Encode(result)
			}

			fmt.Print(FormatMarkdown(result))
			return nil
		},
	}

	rootCmd.Flags().BoolVar(&showComments, "comments", false, "Include comments/replies")
	rootCmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
