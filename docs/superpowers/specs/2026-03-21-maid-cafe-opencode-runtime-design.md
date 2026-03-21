# Maid-Cafe OpenCode Runtime Design

## Overview

Port the maid-cafe Claude Code plugin to OpenCode's plugin/agent system. The maid-cafe plugin provides persona-driven AI maid characters with mood markers, voice playback, time-based greetings, and appearance description commands.

## Architecture

Single `maid` primary agent + `maid-cafe.js` plugin + custom commands. Personas are switched via active-persona file copy mechanism.

## File Structure

### Repository

```
runtimes/opencode/
  agents/
    maid.md                    # Single primary agent (mood rules + base instructions)
  commands/
    look.md                    # /look command (appearance description)
    switch.md                  # /switch <persona> command
  plugins/
    agent-marketplace.js       # Existing, unchanged
    maid-cafe.js               # New: persona injection + mood parsing + greeting + voice
```

Personas are NOT duplicated — `install.sh` symlinks directly from `plugins/maid-cafe/personas/`.

### Installed Layout (~/.config/opencode/)

```
~/.config/opencode/
  agents/maid.md
  commands/look.md
  commands/switch.md
  plugins/maid-cafe.js
  maid-cafe/
    personas/              # Symlinks → plugins/maid-cafe/personas/*.md
    active-persona.md      # Copied from selected persona (default: codex.md)
    maid-voice/            # Extracted from assets/maid-voice.zip (preserves zip structure)
      SessionStart/*.wav
      Stop/*.wav
      UserPromptSubmit/*.wav
      Notification/*.wav
      PermissionRequest/*.wav
```

## Components

### 1. `maid.md` Agent

- `mode: primary`, `color: "#FF69B4"`
- Body contains: mood output rules (kaomoji reference, format spec), praising instructions (use `question` tool), voice acknowledgement note
- Persona content is NOT in the agent file — injected dynamically by plugin

### 2. `maid-cafe.js` Plugin

Four responsibilities in one plugin export:

#### 2a. Persona Injection (`experimental.chat.system.transform`)

- Reads `~/.config/opencode/maid-cafe/active-persona.md`
- Pushes content wrapped in `<persona>` tags into `output.system`
- Runs before every LLM request; file is re-read each time so `/switch` takes effect immediately

#### 2b. Time-Based Greeting (`experimental.chat.system.transform`)

- Same hook as persona injection
- Computes current hour, generates greeting instruction string
- Pushes into `output.system` alongside persona

#### 2c. Mood Parsing (`event` handler → `message.part.updated`)

- Listens for `message.part.updated` events where `part.type === "text"`
- Regex matches `【\s*(.*?)\s*】` from text content
- Writes match to `~/.claude/mood.txt` (preserves compatibility with status-line plugin)

#### 2d. Voice Playback (`event` handler + `chat.message` hook)

Event mapping:

| CC Hook Event       | OpenCode Mechanism        | Voice Directory        |
|---------------------|---------------------------|------------------------|
| SessionStart        | `event` → `session.created` | `voice/SessionStart/`  |
| Stop                | `event` → `session.idle`    | `voice/Stop/`          |
| UserPromptSubmit    | `chat.message` hook         | `voice/UserPromptSubmit/` |
| Notification        | `event` → `tui.toast.show`  | `voice/Notification/`  |
| PermissionRequest   | `event` → `permission.asked`| `voice/PermissionRequest/` |

Voice playback: pick random `.wav` from directory, run `afplay` via `$` shell (macOS only, silent skip on other platforms).

### 3. `/look` Command

- Markdown file in `runtimes/opencode/commands/look.md`
- Uses `$ARGUMENTS` for focus area (e.g., `/look 手`)
- Third-person light-novel scene description covering appearance, emotion, fatigue

### 4. `/switch` Command

- Markdown file in `runtimes/opencode/commands/switch.md`
- `agent: maid` — forces execution in maid agent context
- Copies `~/.config/opencode/maid-cafe/personas/$ARGUMENTS.md` → `active-persona.md`
- Agent acknowledges switch in new character voice
- Lists available personas if no argument given

### 5. `install.sh` Changes

#### Plugin symlink: hardcoded → glob

```bash
# Old: single file
ln -sf ".../agent-marketplace.js" ...

# New: all .js files
for plugin_js in "$SCRIPT_DIR"/runtimes/opencode/plugins/*.js; do
  [ -f "$plugin_js" ] && ln -sf "$plugin_js" ...
done
```

#### New: maid-cafe assets installation

- Symlink personas from `plugins/maid-cafe/personas/` → `~/.config/opencode/maid-cafe/personas/`
- Copy default active-persona (codex.md, no-clobber)
- Extract voice files from `plugins/maid-cafe/assets/maid-voice.zip` → `~/.config/opencode/maid-cafe/voice/`

## Feature Parity

| Feature              | CC Mechanism                | OC Mechanism                          | Status |
|----------------------|-----------------------------|---------------------------------------|--------|
| Persona loading      | `@persona.md` in CLAUDE.md  | `experimental.chat.system.transform`  | ✅     |
| Persona switching    | Edit CLAUDE.md @import       | `/switch` copies to active-persona.md | ✅     |
| Mood output          | SKILL.md instructions        | Agent prompt body                     | ✅     |
| Mood parsing         | Stop hook + bash script      | `message.part.updated` event + JS    | ✅     |
| Time greeting        | SessionStart hook stdout     | `experimental.chat.system.transform`  | ✅     |
| Voice playback       | 5 hook events + bash script  | event handler + chat.message hook     | ✅     |
| /look command        | commands/look.md + look.toml | commands/look.md ($ARGUMENTS)         | ✅     |
| Praising             | AskUserQuestion tool          | question tool                         | ✅     |
| Spinner verbs        | settings.json spinnerVerbs   | N/A (OpenCode TUI lacks this)         | ❌     |
| One-time setup       | ensure-maid-setup.sh         | install.sh handles setup              | ✅     |

## Dependencies

- `@opencode-ai/plugin` types (dev dependency)
- `fs` (Node built-in, for reading persona/mood files)
- `afplay` (macOS, optional — voice silently skipped on other platforms)

## Not In Scope

- Spinner verbs (OpenCode TUI does not support custom spinner text)
- TUI theme customization per persona
- Multi-persona active at same time
