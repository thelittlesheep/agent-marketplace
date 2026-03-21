# agent-marketplace

[繁體中文](README.zh-TW.md)

A plugin marketplace for AI coding assistants. Each plugin bundles skills, agents, CLI tools, commands, and hooks into a self-contained package that works across multiple AI runtimes.

## Supported Runtimes

| Feature | Claude Code | OpenCode |
|---------|:-----------:|:--------:|
| Skills | `.claude-plugin/` | `~/.config/opencode/skills/` |
| Commands | `.toml` | `.md` |
| Agents | `.md` + hooks | `.md` + permissions |
| Hooks | Plugin manifest | JS plugin |
| CLI Tools | `~/.local/bin/` | `~/.local/bin/` |

## Plugins

| Plugin | What it does |
|--------|-------------|
| **code-quality-suite** | Dead code detection, SOLID violation analysis, code explanation |
| **code-review-router** | Routes code review to Gemini or OpenCode based on diff complexity |
| **social-reader** | Fetches and parses X/Twitter and Threads posts |
| **ui-ux-pro-max** | Local BM25 search over curated UI/UX design knowledge |
| **english-coach** | Background grammar checker that runs as a hook |
| **agent-docs** | `.agent/` documentation methodology for project context |
| **status-line** | Terminal status bar with model info, context usage, rate limits, session time |
| **deep-research** | Source-first research pipeline with inline citations and verification |
| **maid-cafe** | Maid cafe persona - mood display, voice hooks, /look command, time-based greetings |

## Install

### Claude Code

This repo is a Claude Code marketplace. Register and install plugins directly:

```bash
claude plugin marketplace add /path/to/agent-marketplace
claude plugin install code-quality-suite  # install individual plugins
```

### OpenCode

```bash
./install-opencode.sh
```

This will:
- Build all Go CLI tools and install binaries to `~/.local/bin/`
- Install default configs to `~/.config/`
- Symlink skills, commands, agents, and plugin to `~/.config/opencode/`

## Prerequisites

- Go 1.24+
- `gemini` CLI and/or `opencode` CLI (optional, for code-review-router)

## Structure

```
plugins/
├── code-quality-suite/    # skills + agent
├── code-review-router/    # skill + Go CLI + hook
├── social-reader/         # skill + Go CLI
├── ui-ux-pro-max/         # skill + Go CLI + CSV data
├── english-coach/         # Go CLI + hook
├── agent-docs/            # skill + references
├── status-line/           # Go CLI + bash script
├── deep-research/         # skill + references
└── maid-cafe/             # skill + commands + hooks (shell scripts)
```

Each plugin follows the same layout: a `plugin.json` manifest under `.claude-plugin/`, with optional `skills/`, `agents/`, `tools/`, `commands/`, and `hooks/` directories.

## Architecture

```
agent-marketplace/
├── .claude-plugin/          # Claude Code marketplace manifest
├── runtimes/opencode/       # OpenCode-specific commands, agents, plugin
├── plugins/                 # Plugin source (shared across runtimes)
└── install-opencode.sh      # Installs for OpenCode; prints Claude Code instructions
```

### Runtime-Specific Limitations

| Feature | Limitation |
|---------|-----------|
| `english-coach` hook | Claude Code only (depends on transcript lifecycle) |
| `status-line` | Claude Code only (depends on status bar API) |
| Agent model selection | Each runtime manages its own model; agents use suggestions, not enforcement |
| `validate-review-cli` hook | Claude Code only (OpenCode has no hook system) |
| `maid-cafe` voice hooks | macOS only (depends on `afplay`); Claude Code only |
| `maid-cafe` mood display | Requires `status-line` plugin (reads `~/.claude/mood.txt`) |
| `maid-cafe` spinnerVerbs | Auto-injected on first session; requires `jq` |
| `maid-cafe` voice files | Auto-downloaded on first session from `claudecafe.com` |
