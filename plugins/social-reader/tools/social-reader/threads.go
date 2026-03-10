package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func init() {
	RegisterPlatform(&ThreadsPlatform{})
}

type ThreadsPlatform struct{}

func (t *ThreadsPlatform) Name() string { return "threads" }

func (t *ThreadsPlatform) MatchURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	host := strings.TrimPrefix(strings.ToLower(u.Host), "www.")
	return host == "threads.net" || host == "threads.com"
}

func (t *ThreadsPlatform) NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Force https
	u.Scheme = "https"

	// Normalize host: threads.com → threads.net, ensure www.
	host := strings.TrimPrefix(strings.ToLower(u.Host), "www.")
	if host == "threads.com" {
		host = "threads.net"
	}
	u.Host = "www." + host

	// Strip query params and fragment
	u.RawQuery = ""
	u.Fragment = ""

	// Validate path
	path := u.Path
	if !threadsPostRe.MatchString(path) && !threadsShortRe.MatchString(path) {
		return "", fmt.Errorf("unsupported Threads URL (only post URLs are supported): %s", rawURL)
	}

	return u.String(), nil
}

var (
	threadsPostRe  = regexp.MustCompile(`^/@[^/]+/post/`)
	threadsShortRe = regexp.MustCompile(`^/t/`)
)

// Profile picture patterns to strip
var (
	profilePicLinkRe = regexp.MustCompile(`\[!\[.*?profile picture.*?\]\(.*?\)\]\(.*?\)`)
	profilePicImgRe  = regexp.MustCompile(`!\[.*?profile picture.*?\]\(.*?\)`)
	translateLineRe  = regexp.MustCompile(`(?i)^Translate\s*$`)
)

func (t *ThreadsPlatform) ParsePost(markdown string, sourceURL string) (*Post, error) {
	cleaned := stripThreadsNoise(markdown)
	lines := strings.Split(cleaned, "\n")

	post := &Post{
		SourceURL: sourceURL,
		Platform:  "threads",
	}

	// Extract author from URL
	u, _ := url.Parse(sourceURL)
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) > 0 {
		post.Author = strings.TrimPrefix(parts[0], "@")
	}

	// Also try to find author in markdown
	if post.Author == "" {
		for _, line := range lines {
			if m := regexp.MustCompile(`@(\w+)`).FindStringSubmatch(line); m != nil {
				post.Author = m[1]
				break
			}
		}
	}

	// Extract post content — first substantial text block after author
	post.Content = extractThreadsPostContent(lines, post.Author)

	// Extract content images (not avatars)
	post.Images = extractThreadsImages(lines)

	// Extract engagement
	post.Engagement = extractThreadsPostEngagement(markdown)

	return post, nil
}

func stripThreadsNoise(markdown string) string {
	// Remove profile picture links
	result := profilePicLinkRe.ReplaceAllString(markdown, "")
	result = profilePicImgRe.ReplaceAllString(result, "")

	// Remove standalone "Translate" lines
	var lines []string
	for _, line := range strings.Split(result, "\n") {
		if translateLineRe.MatchString(strings.TrimSpace(line)) {
			continue
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func extractThreadsPostContent(lines []string, author string) string {
	var content []string
	foundAuthor := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip until we find the author
		if !foundAuthor {
			if strings.Contains(trimmed, "@"+author) || strings.Contains(trimmed, author) {
				foundAuthor = true
			}
			continue
		}

		// Skip empty lines immediately after author
		if len(content) == 0 && trimmed == "" {
			continue
		}

		// Stop at engagement metrics
		if isThreadsEngagementLine(trimmed) {
			break
		}

		// Stop at reply section markers
		if strings.Contains(trimmed, "·Author") && len(content) > 0 {
			break
		}

		// Stop at navigation noise
		if isNavigationNoise(trimmed) {
			break
		}

		if trimmed != "" {
			content = append(content, trimmed)
		}
	}

	return strings.TrimSpace(strings.Join(content, "\n"))
}

func extractThreadsImages(lines []string) []string {
	var images []string
	imgRe := regexp.MustCompile(`!\[([^\]]*)\]\((https?://[^)]+)\)`)
	for _, line := range lines {
		for _, m := range imgRe.FindAllStringSubmatch(line, -1) {
			alt := m[1]
			url := m[2]
			// Skip profile pictures
			if strings.Contains(strings.ToLower(alt), "profile picture") {
				continue
			}
			if strings.Contains(url, "s150x150") {
				continue
			}
			images = append(images, url)
		}
	}
	return images
}

func isThreadsEngagementLine(line string) bool {
	return strings.Contains(line, "❤️") || strings.Contains(line, "💬") ||
		strings.Contains(line, "🔁") || strings.Contains(line, "🔄")
}

func extractThreadsPostEngagement(markdown string) *Engagement {
	e := &Engagement{}
	found := false

	for _, line := range strings.Split(markdown, "\n") {
		trimmed := strings.TrimSpace(line)
		if !isThreadsEngagementLine(trimmed) {
			continue
		}

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
		if m := regexp.MustCompile(`🔄\s*(\d[\d,]*)`).FindStringSubmatch(trimmed); m != nil {
			n := parseNumber(m[1])
			e.Quotes = &n
			found = true
		}

		if found {
			return e
		}
	}

	return nil
}

func (t *ThreadsPlatform) ParseComments(markdown string) ([]*Comment, error) {
	cleaned := stripThreadsNoise(markdown)
	lines := strings.Split(cleaned, "\n")

	var comments []*Comment
	var currentAuthor string
	var currentContent []string
	var currentIsAuthor bool
	var currentEngagement *Engagement

	// Skip past the main post engagement line to find comments
	inComments := false
	seenMainEngagement := false

	flushComment := func() {
		if currentAuthor == "" {
			return
		}
		content := strings.TrimSpace(strings.Join(currentContent, "\n"))
		if content == "" {
			currentAuthor = ""
			currentContent = nil
			currentIsAuthor = false
			currentEngagement = nil
			return
		}

		// Conservative nesting: all comments are root-level
		comments = append(comments, &Comment{
			Author:      currentAuthor,
			Content:     content,
			Engagement:  currentEngagement,
			IsAuthor:    currentIsAuthor,
			parentIndex: -1,
		})

		currentAuthor = ""
		currentContent = nil
		currentIsAuthor = false
		currentEngagement = nil
	}

	authorLineRe := regexp.MustCompile(`@(\w+)`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect end of main post
		if !seenMainEngagement && isThreadsEngagementLine(trimmed) {
			seenMainEngagement = true
			continue
		}

		if !seenMainEngagement {
			continue
		}

		if !inComments {
			inComments = true
		}

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Skip navigation noise
		if isNavigationNoise(trimmed) {
			continue
		}

		// Engagement line → attach to current comment and flush
		if isThreadsEngagementLine(trimmed) {
			eng := &Engagement{}
			if m := regexp.MustCompile(`❤️\s*(\d[\d,]*)`).FindStringSubmatch(trimmed); m != nil {
				n := parseNumber(m[1])
				eng.Likes = &n
			}
			if m := regexp.MustCompile(`💬\s*(\d[\d,]*)`).FindStringSubmatch(trimmed); m != nil {
				n := parseNumber(m[1])
				eng.Replies = &n
			}
			currentEngagement = eng
			flushComment()
			continue
		}

		// Check for ·Author marker
		isAuthor := strings.Contains(trimmed, "·Author")
		if isAuthor {
			trimmed = strings.ReplaceAll(trimmed, "·Author", "")
			trimmed = strings.TrimSpace(trimmed)
		}

		// New comment author line
		if authorM := authorLineRe.FindStringSubmatch(trimmed); authorM != nil && currentAuthor == "" {
			flushComment()
			currentAuthor = authorM[1]
			currentIsAuthor = isAuthor
			// If there's content after the @author on the same line, grab it
			afterAuthor := strings.TrimSpace(trimmed[strings.Index(trimmed, "@"+currentAuthor)+len("@"+currentAuthor):])
			if afterAuthor != "" {
				currentContent = append(currentContent, afterAuthor)
			}
			continue
		}

		// Content line
		if currentAuthor != "" {
			currentContent = append(currentContent, trimmed)
		}
	}

	flushComment()
	return comments, nil
}
