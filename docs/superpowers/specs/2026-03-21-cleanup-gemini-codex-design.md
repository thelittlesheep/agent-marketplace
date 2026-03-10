# Cleanup: Remove Gemini Extension & Codex Support

## Goal

Narrow project scope to support only **Claude Code** (marketplace/plugin system) and **OpenCode** (install script). Remove Gemini CLI extension infrastructure and Codex CLI support entirely.

## Scope

### Delete

| Target | Reason |
|--------|--------|
| `gemini-extension.json` | Gemini extension manifest |
| `GEMINI.md` | Gemini context file |
| `agents/` (entire directory) | Gemini subagent definitions |
| `commands/` (entire directory) | Gemini CLI .toml symlinks |
| `skills/` (entire directory) | Gemini CLI skill symlinks |
| `hooks/` (entire directory) | Gemini CLI hooks |
| `.codex/` (entire directory) | Codex CLI empty remnant |
| `install.sh` | Replaced by `install-opencode.sh` |

### Create

| Target | Description |
|--------|-------------|
| `install-opencode.sh` | OpenCode-only install script (symlinks to `~/.config/opencode/`) |

### Modify

| Target | Changes |
|--------|---------|
| `.gitignore` | Add `.codex/`, remove `skills/` and Gemini-related entries |
| `README.md` | Update to Claude Code + OpenCode only |
| `README.zh-TW.md` | Same as above |

### Preserve (no changes)

- `plugins/code-review-router/` — Gemini routing (`execGemini()`) stays as internal tool usage
- `plugins/english-coach/` — Gemini CLI detection stays
- `runtimes/opencode/` — OpenCode support layer
- `.claude-plugin/` — Claude Code marketplace manifests
- `shared/` — Shared build scripts

## Approach

Single commit: delete files → create `install-opencode.sh` → update `.gitignore` → update READMEs.
