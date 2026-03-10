package main

type DomainConfig struct {
	File       string
	SearchCols []string
	OutputCols []string
}

var DomainConfigs = map[string]DomainConfig{
	"style": {
		File:       "data/styles.csv",
		SearchCols: []string{"Style Category", "Keywords", "Best For", "Type"},
		OutputCols: []string{"Style Category", "Type", "Keywords", "Primary Colors", "Effects & Animation", "Best For", "Performance", "Accessibility", "Framework Compatibility", "Complexity"},
	},
	"prompt": {
		File:       "data/prompts.csv",
		SearchCols: []string{"Style Category", "AI Prompt Keywords (Copy-Paste Ready)", "CSS/Technical Keywords"},
		OutputCols: []string{"Style Category", "AI Prompt Keywords (Copy-Paste Ready)", "CSS/Technical Keywords", "Implementation Checklist"},
	},
	"color": {
		File:       "data/colors.csv",
		SearchCols: []string{"Product Type", "Keywords", "Notes"},
		OutputCols: []string{"Product Type", "Keywords", "Primary (Hex)", "Secondary (Hex)", "CTA (Hex)", "Background (Hex)", "Text (Hex)", "Border (Hex)", "Notes"},
	},
	"chart": {
		File:       "data/charts.csv",
		SearchCols: []string{"Data Type", "Keywords", "Best Chart Type", "Accessibility Notes"},
		OutputCols: []string{"Data Type", "Keywords", "Best Chart Type", "Secondary Options", "Color Guidance", "Accessibility Notes", "Library Recommendation", "Interactive Level"},
	},
	"landing": {
		File:       "data/landing.csv",
		SearchCols: []string{"Pattern Name", "Keywords", "Conversion Optimization", "Section Order"},
		OutputCols: []string{"Pattern Name", "Keywords", "Section Order", "Primary CTA Placement", "Color Strategy", "Conversion Optimization"},
	},
	"product": {
		File:       "data/products.csv",
		SearchCols: []string{"Product Type", "Keywords", "Primary Style Recommendation", "Key Considerations"},
		OutputCols: []string{"Product Type", "Keywords", "Primary Style Recommendation", "Secondary Styles", "Landing Page Pattern", "Dashboard Style (if applicable)", "Color Palette Focus"},
	},
	"ux": {
		File:       "data/ux-guidelines.csv",
		SearchCols: []string{"Category", "Issue", "Description", "Platform"},
		OutputCols: []string{"Category", "Issue", "Platform", "Description", "Do", "Don't", "Code Example Good", "Code Example Bad", "Severity"},
	},
	"typography": {
		File:       "data/typography.csv",
		SearchCols: []string{"Font Pairing Name", "Category", "Mood/Style Keywords", "Best For", "Heading Font", "Body Font"},
		OutputCols: []string{"Font Pairing Name", "Category", "Heading Font", "Body Font", "Mood/Style Keywords", "Best For", "Google Fonts URL", "CSS Import", "Tailwind Config", "Notes"},
	},
}

var StackFiles = map[string]string{
	"html-tailwind": "data/stacks/html-tailwind.csv",
	"react":         "data/stacks/react.csv",
	"nextjs":        "data/stacks/nextjs.csv",
	"vue":           "data/stacks/vue.csv",
	"svelte":        "data/stacks/svelte.csv",
	"swiftui":       "data/stacks/swiftui.csv",
	"react-native":  "data/stacks/react-native.csv",
	"flutter":       "data/stacks/flutter.csv",
}

var StackSearchCols = []string{"Category", "Guideline", "Description", "Do", "Don't"}
var StackOutputCols = []string{"Category", "Guideline", "Description", "Do", "Don't", "Code Good", "Code Bad", "Severity", "Docs URL"}

var DomainKeywords = map[string][]string{
	"color":      {"color", "palette", "hex", "#", "rgb"},
	"chart":      {"chart", "graph", "visualization", "trend", "bar", "pie", "scatter", "heatmap", "funnel"},
	"landing":    {"landing", "page", "cta", "conversion", "hero", "testimonial", "pricing", "section"},
	"product":    {"saas", "ecommerce", "e-commerce", "fintech", "healthcare", "gaming", "portfolio", "crypto", "dashboard"},
	"prompt":     {"prompt", "css", "implementation", "variable", "checklist", "tailwind"},
	"style":      {"style", "design", "ui", "minimalism", "glassmorphism", "neumorphism", "brutalism", "dark mode", "flat", "aurora"},
	"ux":         {"ux", "usability", "accessibility", "wcag", "touch", "scroll", "animation", "keyboard", "navigation", "mobile"},
	"typography": {"font", "typography", "heading", "serif", "sans"},
}

func AllDomains() []string {
	return []string{"style", "prompt", "color", "chart", "landing", "product", "ux", "typography"}
}

func AllStacks() []string {
	return []string{"html-tailwind", "react", "nextjs", "vue", "svelte", "swiftui", "react-native", "flutter"}
}
