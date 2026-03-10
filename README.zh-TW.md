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
