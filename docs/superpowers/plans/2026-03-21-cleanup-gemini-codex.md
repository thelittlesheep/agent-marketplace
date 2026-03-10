# Cleanup: Remove Gemini Extension & Codex Support — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove Gemini CLI extension infrastructure and Codex CLI support, keeping only Claude Code (marketplace) and OpenCode (install script).

**Architecture:** Single-commit cleanup. Delete Gemini/Codex-specific files, create OpenCode-only install script, update .gitignore and READMEs.

**Tech Stack:** Bash, Markdown

---

## Chunk 1: Delete, Create, and Update

### Task 1: Delete Gemini extension infrastructure

**Files:**
- Delete: `gemini-extension.json`
- Delete: `GEMINI.md`
- Delete: `agents/` (entire directory — contains `code-improver.md`, `code-review-router.md`)
- Delete: `commands/` (entire directory — contains .toml symlinks)
- Delete: `skills/` (entire directory — contains skill symlinks)
- Delete: `hooks/` (entire directory — contains `hooks.json`)

- [ ] **Step 1: Delete all Gemini-specific files and directories**

```bash
cd /Users/lsheep/coding/personal/agent-marketplace
rm gemini-extension.json GEMINI.md
rm -rf agents/ commands/ skills/ hooks/
```

- [ ] **Step 2: Verify deletions**

```bash
ls gemini-extension.json GEMINI.md agents/ commands/ skills/ hooks/ 2>&1
```

Expected: all "No such file or directory"

---

### Task 2: Delete Codex CLI remnants

**Files:**
- Delete: `.codex/` (entire directory — empty tmp remnant)

- [ ] **Step 1: Delete .codex directory**

```bash
rm -rf .codex/
```

- [ ] **Step 2: Verify deletion**

```bash
ls -la .codex/ 2>&1
```

Expected: "No such file or directory"

---

### Task 3: Create `install-opencode.sh`

**Files:**
- Delete: `install.sh`
- Create: `install-opencode.sh`

The new script keeps `build_tools()`, `install_config()`, and `install_opencode()` from the original `install.sh`. Removes `install_codex()` and all Gemini/Codex references from output messages.

- [ ] **Step 1: Create `install-opencode.sh`**

```bash
#!/bin/bash
set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Build Go tools -> ~/.local/bin/
build_tools() {
  if ! command -v go &>/dev/null; then
    echo "Go not found. Skipping tool builds."
    return
  fi
  for tool_dir in "$SCRIPT_DIR"/plugins/*/tools/*/; do
    if [ -f "$tool_dir/Makefile" ]; then
      echo "Building $(basename "$tool_dir")..."
      make -C "$tool_dir" install 2>&1 | tail -1
    fi
  done
}

# Install default config (no-clobber)
install_config() {
  local cfg_dir="$HOME/.config/review-cli"
  # Remove dangling symlink if target doesn't exist
  if [ -L "$cfg_dir" ] && [ ! -e "$cfg_dir" ]; then
    rm "$cfg_dir"
  fi
  mkdir -p "$cfg_dir"
  if [ ! -f "$cfg_dir/config.toml" ]; then
    cp "$SCRIPT_DIR/plugins/code-review-router/config/review-cli/config.toml" \
      "$cfg_dir/config.toml"
  fi
}

# OpenCode: symlink skills + commands + agents + plugin
install_opencode() {
  shopt -s nullglob
  # Skills
  mkdir -p "$HOME/.config/opencode/skills"
  for skill_dir in "$SCRIPT_DIR"/plugins/*/skills/*/; do
    ln -sf "$skill_dir" "$HOME/.config/opencode/skills/$(basename "$skill_dir")"
  done

  # Commands (.md format with $ARGUMENTS)
  mkdir -p "$HOME/.config/opencode/commands"
  for cmd in "$SCRIPT_DIR"/runtimes/opencode/commands/*.md; do
    [ -f "$cmd" ] && ln -sf "$cmd" "$HOME/.config/opencode/commands/$(basename "$cmd")"
  done

  # Agents
  mkdir -p "$HOME/.config/opencode/agents"
  for agent in "$SCRIPT_DIR"/runtimes/opencode/agents/*.md; do
    [ -f "$agent" ] && ln -sf "$agent" "$HOME/.config/opencode/agents/$(basename "$agent")"
  done

  # Plugin (system prompt injection)
  mkdir -p "$HOME/.config/opencode/plugins"
  ln -sf "$SCRIPT_DIR/runtimes/opencode/plugins/agent-marketplace.js" \
    "$HOME/.config/opencode/plugins/agent-marketplace.js"

  shopt -u nullglob
  echo "  OpenCode: skills + commands + agents + plugin installed"
}

echo "Installing agent-marketplace for OpenCode..."
echo ""
build_tools || echo "  Warning: some tools failed to build"
install_config || echo "  Warning: config install failed"
echo ""

echo "Linking OpenCode assets..."
install_opencode || echo "  OpenCode: skipped (error during install)"
echo ""

echo "Done!"
echo ""
echo "For Claude Code, use the marketplace system instead:"
echo "  claude plugin marketplace add $SCRIPT_DIR"
echo "  claude plugin install <plugin-name>"
```

- [ ] **Step 2: Make executable and delete old install.sh**

```bash
chmod +x install-opencode.sh
rm install.sh
```

- [ ] **Step 3: Verify**

```bash
head -3 install-opencode.sh && ls -la install-opencode.sh install.sh 2>&1
```

Expected: shebang shown, `install-opencode.sh` exists with +x, `install.sh` not found

---

### Task 4: Update `.gitignore`

**Files:**
- Modify: `.gitignore`

Remove `skills/` entry (no longer needed — directory deleted). Remove Codex/OpenCode comment. Add `.codex/`.

- [ ] **Step 1: Replace `.gitignore` contents**

New content:

```
.DS_Store

# Compiled Go binaries
plugins/*/tools/*/review-cli
plugins/*/tools/*/social-reader
plugins/*/tools/*/uiux-cli
plugins/*/tools/*/english-coach
plugins/*/tools/*/status-line

# Codex CLI temp files
.codex/
```

---

### Task 5: Update `README.md`

**Files:**
- Modify: `README.md`

Remove all Codex/Gemini columns from tables, Gemini install instructions, Gemini architecture entries. Update install section to reference `install-opencode.sh`. Update prerequisites.

- [ ] **Step 1: Replace `README.md` contents**

```markdown
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
└── status-line/           # Go CLI + bash script
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
```

---

### Task 6: Update `README.zh-TW.md`

**Files:**
- Modify: `README.zh-TW.md`

Same changes as README.md but in Traditional Chinese.

- [ ] **Step 1: Replace `README.zh-TW.md` contents**

```markdown
# agent-marketplace

[English](README.md)

給 AI coding assistant 用的 plugin marketplace。每個 plugin 把 skills、agents、CLI 工具、commands 和 hooks 打包成獨立套件，可以跨多個 AI runtime 使用。

## 支援的 Runtime

| 功能 | Claude Code | OpenCode |
|------|:-----------:|:--------:|
| Skills | `.claude-plugin/` | `~/.config/opencode/skills/` |
| Commands | `.toml` | `.md` |
| Agents | `.md` + hooks | `.md` + permissions |
| Hooks | Plugin manifest | JS plugin |
| CLI Tools | `~/.local/bin/` | `~/.local/bin/` |

## Plugins

| Plugin | 說明 |
|--------|------|
| **code-quality-suite** | 死碼偵測、SOLID 原則分析、程式碼解說 |
| **code-review-router** | 根據 diff 複雜度自動分派 code review 給 Gemini 或 OpenCode |
| **social-reader** | 抓取並解析 X/Twitter 和 Threads 貼文 |
| **ui-ux-pro-max** | 本地 BM25 搜尋，查詢 UI/UX 設計知識庫 |
| **english-coach** | 以 hook 形式在背景執行的英文文法檢查 |
| **agent-docs** | `.agent/` 文件方法論，維護專案上下文 |
| **status-line** | 終端機狀態列，顯示模型資訊、context 用量、rate limit、session 時間 |

## 安裝

### Claude Code

本 repo 是 Claude Code marketplace，直接註冊並安裝 plugin：

```bash
claude plugin marketplace add /path/to/agent-marketplace
claude plugin install code-quality-suite  # 安裝個別 plugin
```

### OpenCode

```bash
./install-opencode.sh
```

這會：
- 編譯所有 Go CLI 工具，安裝到 `~/.local/bin/`
- 安裝預設設定檔到 `~/.config/`
- Symlink skills、commands、agents 和 plugin 到 `~/.config/opencode/`

## 前置需求

- Go 1.24+
- `gemini` CLI 和/或 `opencode` CLI（選用，給 code-review-router 用）

## 結構

```
plugins/
├── code-quality-suite/    # skills + agent
├── code-review-router/    # skill + Go CLI + hook
├── social-reader/         # skill + Go CLI
├── ui-ux-pro-max/         # skill + Go CLI + CSV data
├── english-coach/         # Go CLI + hook
├── agent-docs/            # skill + references
└── status-line/           # Go CLI + bash script
```

每個 plugin 遵循相同結構：`.claude-plugin/` 下放 `plugin.json` manifest，搭配選用的 `skills/`、`agents/`、`tools/`、`commands/`、`hooks/` 目錄。

## 架構

```
agent-marketplace/
├── .claude-plugin/          # Claude Code marketplace manifest
├── runtimes/opencode/       # OpenCode 專屬 commands、agents、plugin
├── plugins/                 # Plugin 原始碼（所有 runtime 共用）
└── install-opencode.sh      # 安裝 OpenCode；印出 Claude Code 指示
```

### Runtime 專屬限制

| 功能 | 限制 |
|------|------|
| `english-coach` hook | 僅 Claude Code（依賴 transcript 生命週期） |
| `status-line` | 僅 Claude Code（依賴 status bar API） |
| Agent model 選擇 | 各 runtime 自行管理；agent 用建議而非強制 |
| `validate-review-cli` hook | 僅 Claude Code（OpenCode 無 hook 系統） |
```

---

### Task 7: Commit

- [ ] **Step 1: Stage all changes**

```bash
git add -A
```

- [ ] **Step 2: Commit**

```bash
git commit -m "refactor: remove Gemini extension and Codex CLI support

Drop Gemini CLI extension infrastructure (gemini-extension.json, GEMINI.md,
agents/, commands/, skills/, hooks/) and Codex CLI support (.codex/,
install_codex). Project now supports Claude Code (marketplace) and
OpenCode (install script) only.

- Replace install.sh with install-opencode.sh (OpenCode-only)
- Update .gitignore to remove skills/ and add .codex/
- Update both READMEs to reflect two-runtime architecture
- Preserve Gemini CLI usage within code-review-router and english-coach"
```
