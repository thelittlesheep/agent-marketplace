package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var (
		domain     string
		stack      string
		maxResults int
		asJSON     bool
	)

	searchCmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search UI/UX design database",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")

			var result SearchResult
			var err error

			if stack != "" {
				result, err = SearchStack(query, stack, maxResults)
			} else {
				result, err = Search(query, domain, maxResults)
			}
			if err != nil {
				return err
			}

			fmt.Print(FormatSearchResult(result, asJSON))
			return nil
		},
	}

	searchCmd.Flags().StringVarP(&domain, "domain", "d", "", "Search domain (style, prompt, color, chart, landing, product, ux, typography)")
	searchCmd.Flags().StringVarP(&stack, "stack", "s", "", "Search stack (html-tailwind, react, nextjs, vue, svelte, swiftui, react-native, flutter)")
	searchCmd.Flags().IntVarP(&maxResults, "max-results", "n", 3, "Maximum results")
	searchCmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")

	domainsCmd := &cobra.Command{
		Use:   "domains",
		Short: "List available search domains",
		Run: func(cmd *cobra.Command, args []string) {
			for _, d := range AllDomains() {
				fmt.Println(d)
			}
		},
	}

	stacksCmd := &cobra.Command{
		Use:   "stacks",
		Short: "List available stacks",
		Run: func(cmd *cobra.Command, args []string) {
			for _, s := range AllStacks() {
				fmt.Println(s)
			}
		},
	}

	rootCmd := &cobra.Command{
		Use:   "uiux-cli",
		Short: "UI/UX design intelligence - search styles, colors, typography, and more",
	}

	rootCmd.AddCommand(searchCmd, domainsCmd, stacksCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
