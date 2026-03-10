---
name: code-review-router
description: "Use when completing code changes and wanting automated code review - routes to Gemini CLI (fast, simple changes) or OpenCode CLI (deep analysis) based on complexity, security sensitivity, and change scope. Triggers on: code review, review my changes, check code quality, PR review, 程式碼審查, 幫我 review, 檢查程式碼."
---

# Code Review Router

Requires `gemini` or `opencode` CLI on PATH.

## Usage

```bash
review-cli                          # Review uncommitted changes
review-cli --staged                 # Staged changes only
review-cli --commit <hash>          # Single commit changes
review-cli --range <from>..<to>     # Commit range changes
review-cli --dry-run                # Preview routing decision without executing
review-cli --force gemini           # Force Gemini CLI
review-cli --force opencode         # Force OpenCode CLI
review-cli --json                   # Machine-readable output
```

## Config

`~/.config/review-cli/config.toml` — all optional, has built-in defaults.
