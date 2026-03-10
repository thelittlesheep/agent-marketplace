package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	RegisterPlatform(&XPlatform{})
}

type XPlatform struct{}

func (x *XPlatform) Name() string { return "x" }

func (x *XPlatform) MatchURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	host := strings.TrimPrefix(strings.ToLower(u.Host), "www.")
	return host == "x.com" || host == "twitter.com"
}

func (x *XPlatform) NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Force https
	u.Scheme = "https"

	// Normalize host
	host := strings.TrimPrefix(strings.ToLower(u.Host), "www.")
	if host == "twitter.com" {
		host = "x.com"
	}
	u.Host = host

	// Strip query params and fragment
	u.RawQuery = ""
	u.Fragment = ""

	// Validate path is a status or article URL
	path := u.Path
	if !xStatusRe.MatchString(path) && !xArticleRe.MatchString(path) {
		return "", fmt.Errorf("unsupported X URL (only tweet and article URLs are supported): %s", rawURL)
	}

	return u.String(), nil
}

var (
	xStatusRe            = regexp.MustCompile(`^/[^/]+/status/\d+`)
	xArticleRe           = regexp.MustCompile(`^/[^/]+/article/`)
	xTimestampRe         = regexp.MustCompile(`^\[\d+:\d+\s+[AP]M\s+·`)
	xViewsRe             = regexp.MustCompile(`(?i)\d[\d.,]*[KkMm]?\s+views`)
	jinaHeaderRe         = regexp.MustCompile(`^(Title|URL Source|Markdown Content|Published Time|Byline):`)
	markdownH1UnderlineRe = regexp.MustCompile(`^=+$`)
	linkOnlyLineRe       = regexp.MustCompile(`^\[.+\]\(.+\)$`)
)

func isXArticleURL(sourceURL string) bool {
	u, err := url.Parse(sourceURL)
	if err != nil {
		return false
	}
	return xArticleRe.MatchString(u.Path)
}

func isLinkOnlyLine(line string) bool {
	return linkOnlyLineRe.MatchString(strings.TrimSpace(line))
}

func (x *XPlatform) ParsePost(markdown string, sourceURL string) (*Post, error) {
	lines := strings.Split(markdown, "\n")
	post := &Post{
		SourceURL: sourceURL,
		Platform:  "x",
	}

	// Find author — look for the first @ pattern in a heading or bold
	for _, line := range lines {
		if m := xAuthorRe.FindStringSubmatch(line); m != nil {
			post.Author = m[1]
			break
		}
	}
	if post.Author == "" {
		// Try to extract from URL
		u, _ := url.Parse(sourceURL)
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) > 0 {
			post.Author = parts[0]
		}
	}

	// Extract content — route by URL type
	if isXArticleURL(sourceURL) {
		post.Content = extractXArticleContent(lines)
	} else {
		post.Content = extractXPostContent(lines)
	}

	// Extract images
	post.Images = extractImages(lines)

	// Extract engagement
	post.Engagement = extractXEngagement(markdown)

	// Extract quoted post
	post.QuotedPost = extractXQuotedPost(lines)

	return post, nil
}

var xAuthorRe = regexp.MustCompile(`@(\w+)`)

func extractXPostContent(lines []string) string {
	var content []string
	inContent := false
	authorFound := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines at start
		if !authorFound {
			if xAuthorRe.MatchString(trimmed) {
				authorFound = true
			}
			continue
		}

		// Skip lines that are just the author handle
		if !inContent {
			if trimmed == "" {
				continue
			}
			// Start collecting content after we pass the author line
			inContent = true
		}

		// Stop at engagement indicators or navigation
		if isXEngagementLine(trimmed) || isNavigationNoise(trimmed) {
			break
		}

		// Stop at "Translate post" line
		if strings.EqualFold(trimmed, "translate post") {
			break
		}

		// Stop at timestamp line: [3:52 AM · Feb 27, 2026](url)
		if xTimestampRe.MatchString(trimmed) {
			break
		}

		// Stop at "Replying to" which marks the comments section
		if strings.HasPrefix(trimmed, "Replying to") {
			break
		}

		// Stop at "Post your reply" or similar
		if strings.Contains(strings.ToLower(trimmed), "post your reply") {
			break
		}

		content = append(content, line)
	}

	return strings.TrimSpace(strings.Join(content, "\n"))
}

func extractXArticleContent(lines []string) string {
	i := 0

	// Phase 1: skip Jina metadata headers
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || jinaHeaderRe.MatchString(trimmed) {
			i++
			continue
		}
		break
	}

	// Phase 2: skip navigation noise (X chrome captured by Jina)
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || isNavigationNoise(trimmed) || isLinkOnlyLine(trimmed) {
			i++
			continue
		}
		break
	}

	// Phase 3: detect article title
	var title string
	if i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		// Setext H1: title line followed by ====
		if i+1 < len(lines) && markdownH1UnderlineRe.MatchString(strings.TrimSpace(lines[i+1])) {
			title = trimmed
			i += 2
		} else if strings.HasPrefix(trimmed, "# ") {
			title = strings.TrimPrefix(trimmed, "# ")
			i++
		}
	}

	// Phase 4: collect content until footer noise
	var content []string
	started := false
	for ; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])

		if !started && trimmed == "" {
			continue
		}
		started = true

		if isNavigationNoise(trimmed) {
			break
		}

		content = append(content, lines[i])
	}

	body := strings.TrimSpace(strings.Join(content, "\n"))
	if title != "" && body != "" {
		return "# " + title + "\n\n" + body
	}
	if title != "" {
		return "# " + title
	}
	return body
}

func extractImages(lines []string) []string {
	var images []string
	imgRe := regexp.MustCompile(`!\[.*?\]\((https?://[^)]+)\)`)
	for _, line := range lines {
		for _, m := range imgRe.FindAllStringSubmatch(line, -1) {
			imgURL := m[1]
			// Skip profile pictures, avatars, emoji SVGs, and hashflags
			if strings.Contains(imgURL, "profile_images") ||
				strings.Contains(imgURL, "s150x150") ||
				strings.Contains(imgURL, "twimg.com/emoji") ||
				strings.Contains(imgURL, "hashflags") {
				continue
			}
			images = append(images, imgURL)
		}
	}
	return images
}

func extractXEngagement(markdown string) *Engagement {
	lines := strings.Split(markdown, "\n")
	e := &Engagement{}
	found := false

	// Strategy 1: Look for emoji-based engagement patterns (e.g. "❤️ 42 · 💬 5 · 🔁 3")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if m := regexp.MustCompile(`❤️\s*(\d[\d,]*)`).FindStringSubmatch(trimmed); m != nil {
			n := parseNumber(m[1])
			e.Likes = &n
			found = true
		}
		if m := regexp.MustCompile(`💬\s*(\d[\d,]*)`).FindStringSubmatch(trimmed); m != nil {
			n := parseNumber(m[1])
			e.Replies = &n
			found = true
		}
		if m := regexp.MustCompile(`🔁\s*(\d[\d,]*)`).FindStringSubmatch(trimmed); m != nil {
			n := parseNumber(m[1])
			e.Reposts = &n
			found = true
		}

		if found {
			return e
		}
	}

	// Strategy 2: Bare numbers after "Views" line
	// Jina renders X engagement as: [12.4K Views](url) then 4 bare-number lines
	// Order: replies, reposts, likes, bookmarks
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if xViewsRe.MatchString(trimmed) {
			nums := collectBareNumbers(lines, i+1, 4)
			if len(nums) >= 3 {
				e.Replies = intPtr(nums[0])
				e.Reposts = intPtr(nums[1])
				e.Likes = intPtr(nums[2])
				return e
			}
			break
		}
	}

	return nil
}

// collectBareNumbers reads up to max consecutive bare integer lines starting from startIdx.
func collectBareNumbers(lines []string, startIdx int, max int) []int {
	var nums []int
	for i := startIdx; i < len(lines) && len(nums) < max; i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			continue
		}
		n, err := strconv.Atoi(strings.ReplaceAll(trimmed, ",", ""))
		if err != nil {
			break
		}
		nums = append(nums, n)
	}
	return nums
}

func isXEngagementLine(line string) bool {
	// Lines that contain engagement metrics
	return strings.Contains(line, "❤️") || strings.Contains(line, "💬") || strings.Contains(line, "🔁")
}

func isNavigationNoise(line string) bool {
	lower := strings.ToLower(line)
	noisePatterns := []string{
		"trending", "who to follow", "what's happening",
		"terms of service", "privacy policy", "cookie",
		"©", "footer", "sidebar",
		"new to x", "sign up", "create account", "log in",
		"don't miss what", "people on x are",
		"personalized timeline",
		"read more replies", "show probable spam",
	}
	for _, p := range noisePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func extractXQuotedPost(lines []string) *Post {
	// Look for blockquote patterns that indicate quoted tweets
	var inQuote bool
	var quoteLines []string
	var quoteAuthor string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "> ") {
			inQuote = true
			quoteLine := strings.TrimPrefix(trimmed, "> ")
			// Try to extract author from first quote line
			if quoteAuthor == "" {
				if m := xAuthorRe.FindStringSubmatch(quoteLine); m != nil {
					quoteAuthor = m[1]
					continue
				}
			}
			quoteLines = append(quoteLines, quoteLine)
		} else if inQuote {
			break
		}
	}

	if len(quoteLines) > 0 {
		return &Post{
			Author:  quoteAuthor,
			Content: strings.TrimSpace(strings.Join(quoteLines, "\n")),
		}
	}
	return nil
}

func (x *XPlatform) ParseComments(markdown string) ([]*Comment, error) {
	return nil, fmt.Errorf("X/Twitter does not support --comments (replies require login and cannot be fetched via Jina Reader)")
}

func parseNumber(s string) int {
	s = strings.ReplaceAll(s, ",", "")
	n, _ := strconv.Atoi(s)
	return n
}
