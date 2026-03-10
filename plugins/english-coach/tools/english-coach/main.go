package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unicode"
)

// State tracks the last processed message hash per session.
type State struct {
	LastHash  string `json:"last_hash"`
	Timestamp int64  `json:"timestamp"`
}

// Available CLI tool name.
type cliTool struct {
	name string
}

var (
	reCodeBlock  = regexp.MustCompile("(?s)```.*?```")
	reInlineCode = regexp.MustCompile("`[^`]+`")
	reURL        = regexp.MustCompile(`https?://\S+`)
	reFilePath   = regexp.MustCompile(`(?:^|\s)(?:[~/.][\w./\\-]+)`)
)

const systemPrompt = `You are an English writing coach. Check the following text for grammar errors, typos, and awkward phrasing. The text may contain Traditional Chinese — only check the English parts.

Rules:
- If the English is all correct, respond with exactly: LGTM
- If there are errors, respond with the single most important correction in this exact format: "wrong" → "correct"
- Only one correction per check
- Focus on: spelling, grammar, word choice, common ESL mistakes
- Ignore: code snippets, file paths, technical terms, URLs, command names
- No explanations, just the correction`

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: english-coach <transcript_path> <session_id>\n")
		os.Exit(1)
	}

	transcriptPath := os.Args[1]
	sessionID := os.Args[2]

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot get home dir: %v\n", err)
		os.Exit(1)
	}

	coachDir := filepath.Join(homeDir, ".claude", "english-coach")
	if err := os.MkdirAll(coachDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create coach dir: %v\n", err)
		os.Exit(1)
	}

	// Acquire per-session file lock.
	lockPath := filepath.Join(coachDir, sessionID+".lock")
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		os.Exit(0) // Silently exit if cannot create lock.
	}
	defer lockFile.Close()

	if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		os.Exit(0) // Another instance is running for this session.
	}
	defer syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)

	// Extract the last user message from transcript.
	rawMessage := extractUserMessage(transcriptPath, sessionID)
	if rawMessage == "" {
		return
	}

	// Filter: skip short messages.
	if len(strings.TrimSpace(rawMessage)) < 5 {
		return
	}

	// Filter: skip system messages.
	if isSystemMessage(rawMessage) {
		return
	}

	// Smart content filtering: strip non-prose, check prose ratio.
	prose := stripNonProse(rawMessage)
	prose = strings.TrimSpace(prose)
	if len(prose) < 5 {
		return
	}
	if proseRatio(prose) < 0.4 {
		return
	}

	// Truncate to avoid overwhelming the LLM with long messages.
	if len(prose) > 500 {
		prose = prose[:500]
	}

	// Hash deduplication.
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(rawMessage)))

	statePath := filepath.Join(coachDir, sessionID+"-state.json")
	tipPath := filepath.Join(coachDir, sessionID+"-tip.txt")

	var state State
	if data, err := os.ReadFile(statePath); err == nil {
		json.Unmarshal(data, &state)
	}
	if state.LastHash == hash {
		return // Already processed this message.
	}

	// Detect available CLI tool.
	cli := detectCLI()
	if cli == nil {
		return // No LLM CLI available.
	}

	// Call LLM.
	tip, err := callLLM(cli, prose)
	if err != nil {
		return // Don't overwrite existing tip on error.
	}

	// Parse and validate output.
	tip = strings.TrimSpace(tip)
	if tip == "" || strings.Contains(tip, "LGTM") {
		tip = ""
	} else if !isValidTip(tip) {
		return // Discard irrelevant LLM output.
	}

	// Write tip file.
	os.WriteFile(tipPath, []byte(tip), 0644)

	// Update state.
	state = State{LastHash: hash, Timestamp: time.Now().Unix()}
	if data, err := json.Marshal(state); err == nil {
		os.WriteFile(statePath, data, 0644)
	}
}

// extractUserMessage finds the last user message in the transcript JSONL.
func extractUserMessage(transcriptPath, sessionID string) string {
	file, err := os.Open(transcriptPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxScanTokenSize = 1024 * 1024
	buf := make([]byte, 0, maxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	var allLines []string
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	// Take last 200 lines.
	start := len(allLines) - 200
	if start < 0 {
		start = 0
	}
	lines := allLines[start:]

	// Search from the end for user message.
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.TrimSpace(line) == "" {
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			continue
		}

		isSidechain, _ := data["isSidechain"].(bool)
		sessionMatch := false
		if sid, ok := data["sessionId"].(string); ok && sid == sessionID {
			sessionMatch = true
		}

		if !isSidechain && sessionMatch {
			if message, ok := data["message"].(map[string]interface{}); ok {
				role, _ := message["role"].(string)
				msgType, _ := data["type"].(string)

				if role == "user" && msgType == "user" {
					if content, ok := message["content"].(string); ok {
						if !isSystemMessage(content) {
							return content
						}
					}
				}
			}
		}
	}

	return ""
}

// isSystemMessage checks if a message is a system/command message.
func isSystemMessage(content string) bool {
	if strings.HasPrefix(content, "[") && strings.HasSuffix(content, "]") {
		return true
	}
	if strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}") {
		return true
	}

	xmlTags := []string{
		"<local-command-stdout>", "<command-name>",
		"<command-message>", "<command-args>",
	}
	for _, tag := range xmlTags {
		if strings.Contains(content, tag) {
			return true
		}
	}

	if strings.HasPrefix(content, "Caveat:") {
		return true
	}

	return false
}

// stripNonProse removes code blocks, inline code, URLs, and file paths.
func stripNonProse(text string) string {
	text = reCodeBlock.ReplaceAllString(text, " ")
	text = reInlineCode.ReplaceAllString(text, " ")
	text = reURL.ReplaceAllString(text, " ")
	text = reFilePath.ReplaceAllString(text, " ")
	return text
}

// proseRatio calculates the ratio of letter characters in the text.
func proseRatio(text string) float64 {
	if len(text) == 0 {
		return 0
	}

	var total, letters int
	for _, r := range text {
		if !unicode.IsSpace(r) {
			total++
			if unicode.IsLetter(r) {
				letters++
			}
		}
	}

	if total == 0 {
		return 0
	}
	return float64(letters) / float64(total)
}

// detectCLI finds an available LLM CLI tool.
func detectCLI() *cliTool {
	if _, err := exec.LookPath("gemini"); err == nil {
		return &cliTool{name: "gemini"}
	}
	if _, err := exec.LookPath("claude"); err == nil {
		return &cliTool{name: "claude"}
	}
	return nil
}

// callLLM invokes the LLM CLI with the combined prompt (system + user message).
func callLLM(cli *cliTool, message string) (string, error) {
	combined := systemPrompt + "\n\nText to check:\n" + message

	var cmd *exec.Cmd
	if cli.name == "gemini" {
		cmd = exec.Command("gemini", "-m", "gemini-2.5-flash", "-p", combined, "--output-format", "text")
	} else {
		cmd = exec.Command("claude", "-p", "--model", "haiku", combined)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return string(output), nil
}

// isValidTip checks if the LLM output is a valid tip (not irrelevant garbage).
func isValidTip(output string) bool {
	// Valid correction must contain → and be a single line.
	if strings.Contains(output, "→") && strings.Count(output, "\n") == 0 {
		return true
	}
	return false
}
