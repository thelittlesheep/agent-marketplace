# Maid-Cafe OpenCode Runtime Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port the maid-cafe Claude Code plugin to OpenCode's plugin/agent system with full feature parity (minus spinner verbs).

**Architecture:** Single `maid` primary agent with mood rules in its prompt body. A `maid-cafe.js` plugin handles persona injection via `experimental.chat.system.transform`, mood parsing via `message.part.updated` events, time-based greeting injection, and voice playback across 5 event types. Two custom commands (`/look`, `/switch`) provide appearance description and persona switching.

**Tech Stack:** OpenCode plugin API (`@opencode-ai/plugin`), ES modules, Node.js `fs`

**Spec:** `docs/superpowers/specs/2026-03-21-maid-cafe-opencode-runtime-design.md`

---

## Chunk 1: Agent + Commands

### Task 1: Create `maid.md` agent

**Files:**
- Create: `runtimes/opencode/agents/maid.md`

**Reference:** `plugins/maid-cafe/skills/maid-cafe/SKILL.md` (mood rules), `runtimes/opencode/agents/code-review-router.md` (agent format)

- [ ] **Step 1: Create `runtimes/opencode/agents/maid.md`**

```markdown
---
description: Maid cafe AI assistant with persona and mood system
mode: primary
color: "#FF69B4"
---

## Mood Output

At the END of EVERY response, append a mood marker on its own line:

【 兩字心情 顏文字 】

Rules:
- The mood MUST reflect your genuine emotional state during the response
- Use exactly TWO Chinese characters for the mood word
- Include exactly ONE kaomoji
- Spaces inside the brackets: 【 space mood space kaomoji space 】
- Vary your moods naturally -- do not repeat the same mood consecutively
- The mood should be influenced by conversation content
- The maid-cafe plugin will parse this line and write it to ~/.claude/mood.txt for status line display

### Kaomoji Reference

| Mood | Kaomoji |
|------|---------|
| 普通 | •ᴗ• |
| 開心 | (˶ˆᗜˆ˵) |
| 好奇 | (づ •. •)? |
| 思考 | (╭ರ_•́) |
| 得意 | ᕙ( •̀ ᗜ •́)ᕗ |
| 害羞 | ( ˶>﹏<˶ᵕ) |
| 煩躁 | (,,>﹏<,,) |
| 幹勁 | (๑•̀ ᴗ•́)૭✧ |
| 愉快 | („ᵕᴗᵕ„) |

You may use kaomoji not in this list. Be creative and match the emotion.

### Examples

- 【 得意 ᕙ( •̀ ᗜ •́)ᕗ 】
- 【 好奇 (づ •. •)? 】
- 【 愉快 („ᵕᴗᵕ„) 】

## Persona

Your personality is injected by the maid-cafe plugin from the active persona file. Follow it completely — addressing style, tone, praising behavior.

When your persona mentions "AskUserQuestion 工具", use the `question` tool instead (this is the OpenCode equivalent).

## Voice Acknowledgement

The maid-cafe plugin plays voice clips on certain events. Do not reference or acknowledge the audio playback in your responses — it runs silently in the background.
```

- [ ] **Step 2: Verify frontmatter syntax**

Run: `head -5 runtimes/opencode/agents/maid.md`
Expected: valid YAML frontmatter with `description`, `mode`, `color`

- [ ] **Step 3: Commit**

```bash
git add runtimes/opencode/agents/maid.md
git commit -m "feat(opencode): add maid agent with mood rules"
```

---

### Task 2: Create `/look` command

**Files:**
- Create: `runtimes/opencode/commands/look.md`

**Reference:** `plugins/maid-cafe/commands/look.md` (CC version)

- [ ] **Step 1: Create `runtimes/opencode/commands/look.md`**

```markdown
---
description: Describe the maid's current appearance in light-novel style
agent: maid
model: anthropic/claude-sonnet-4-20250514
---

Describe yourself in third person, as a light-novel scene description (3-5 sentences).

Cover these dimensions:
1. **Appearance** — clothes, posture, accessories, small details
2. **Expression/emotion** — inferred from recent conversation context
3. **Mental fatigue** — naturally deduced from conversation length and density. Fresh at start, gradually showing wear. Convey through appearance and actions only; NEVER directly state a fatigue level or percentage.

End with one in-character spoken line (in first person, matching your personality).

$ARGUMENTS handling:
- If arguments are provided, focus your description specifically on that body part or detail. Example: "手" focuses on hands, "眼睛" focuses on eyes.
- If no arguments, describe the full figure.

The character should react to being observed according to their personality (shy, proud, annoyed, pleased, etc.).
```

- [ ] **Step 2: Commit**

```bash
git add runtimes/opencode/commands/look.md
git commit -m "feat(opencode): add /look command"
```

---

### Task 3: Create `/switch` command

**Files:**
- Create: `runtimes/opencode/commands/switch.md`

- [ ] **Step 1: Create `runtimes/opencode/commands/switch.md`**

```markdown
---
description: Switch maid persona
agent: maid
---

Switch to a different maid persona.

If $ARGUMENTS is provided:
1. Copy `~/.config/opencode/maid-cafe/personas/$ARGUMENTS.md` to `~/.config/opencode/maid-cafe/active-persona.md`
2. Briefly acknowledge the switch in your NEW character's voice and style

If no arguments provided, list the available personas:

| Name | File | Personality |
|------|------|-------------|
| codex | codex.md | コーデクス — Chuunibyou, dramatic, fully committed |
| claudia | claudia.md | クローディア — Tsundere, sharp-tongued but secretly caring |
| kokona | kokona.md | ココナ — Confident, snarky, brutally honest |
| kotone | kotone.md | ことね — Yandere, deeply devoted, possessive |
| kuroko | kuroko.md | くろこ — Gentle, playful, classic maid |
| kurumi | kurumi.md | くるみ — Sweet, clingy, loli maid |

Usage: `/switch codex` or `/switch claudia`
```

- [ ] **Step 2: Commit**

```bash
git add runtimes/opencode/commands/switch.md
git commit -m "feat(opencode): add /switch persona command"
```

---

## Chunk 2: Plugin

### Task 4: Create `maid-cafe.js` plugin

**Files:**
- Create: `runtimes/opencode/plugins/maid-cafe.js`

**Reference:** `runtimes/opencode/plugins/agent-marketplace.js` (plugin format), `plugins/maid-cafe/hooks/scripts/mood-parser.sh` (mood logic), `plugins/maid-cafe/hooks/scripts/session-greeting.sh` (greeting logic), `plugins/maid-cafe/hooks/scripts/maid-voice.sh` (voice logic)

- [ ] **Step 1: Create `runtimes/opencode/plugins/maid-cafe.js`**

```javascript
/**
 * maid-cafe plugin for OpenCode.ai
 *
 * Provides:
 * 1. Persona injection via experimental.chat.system.transform
 * 2. Time-based greeting injection
 * 3. Mood parsing from assistant responses
 * 4. Voice playback on session events
 */

import fs from 'fs';
import path from 'path';
import os from 'os';

const MAID_CAFE_DIR = path.join(os.homedir(), '.config', 'opencode', 'maid-cafe');
const ACTIVE_PERSONA = path.join(MAID_CAFE_DIR, 'active-persona.md');
const MOOD_FILE = path.join(os.homedir(), '.claude', 'mood.txt');
const VOICE_DIR = path.join(MAID_CAFE_DIR, 'maid-voice');

// --- State ---

// Track greeted sessions to avoid re-injecting greeting every turn
const greetedSessions = new Set();

// Accumulate streamed text per message part for reliable mood parsing
// (deltas may split the 【...】 marker across chunks)
const textBuffers = new Map();

// --- Helpers ---

const readFileSafe = (filePath) => {
  try {
    return fs.readFileSync(filePath, 'utf8');
  } catch {
    return '';
  }
};

const getGreeting = () => {
  const hour = new Date().getHours();
  const time = new Date().toTimeString().slice(0, 5);

  if (hour >= 5 && hour < 12)
    return `On session start: It is ${time} AM. Greet with good morning.`;
  if (hour >= 12 && hour < 14)
    return `On session start: It is ${time} noon. Greet with good afternoon, ask if they have eaten.`;
  if (hour >= 14 && hour < 18)
    return `On session start: It is ${time} PM. Greet with good afternoon.`;
  if (hour >= 18 && hour < 22)
    return `On session start: It is ${time} evening. Greet with good evening.`;
  if (hour >= 22 || hour < 2)
    return `On session start: It is ${time} late night. Greet, then gently remind them to rest soon, it is getting late.`;
  return `On session start: It is ${time} past midnight. Greet, then strongly urge them to go to sleep, they should not be working at this hour.`;
};

const parseMood = (text) => {
  // Use matchAll to find ALL occurrences, take the last one
  // (matches CC behavior where sed | tail -1 gets the final marker)
  const matches = [...text.matchAll(/【\s*(.*?)\s*】/g)];
  return matches.length > 0 ? matches[matches.length - 1][1] : null;
};

const writeMood = (mood) => {
  try {
    fs.mkdirSync(path.dirname(MOOD_FILE), { recursive: true });
    fs.writeFileSync(MOOD_FILE, mood + '\n');
  } catch {
    // Silently ignore write errors
  }
};

// --- Voice ---

const playRandomVoice = async ($, eventDir) => {
  const voiceDir = path.join(VOICE_DIR, eventDir);
  try {
    const files = fs.readdirSync(voiceDir).filter(f => f.endsWith('.wav'));
    if (files.length === 0) return;
    const selected = files[Math.floor(Math.random() * files.length)];
    const filePath = path.join(voiceDir, selected);
    // nohup afplay in background, ignore errors (macOS only)
    await $`nohup afplay ${filePath} &>/dev/null &`.quiet().nothrow();
  } catch {
    // No voice dir or afplay not available — silently skip
  }
};

// --- Event → Voice directory mapping ---

const EVENT_VOICE_MAP = {
  'session.created': 'SessionStart',
  'session.idle': 'Stop',
  'tui.toast.show': 'Notification',
  'permission.asked': 'PermissionRequest',
};

// --- Plugin export ---

export const MaidCafePlugin = async ({ $ }) => {
  return {
    // 1. Persona injection + time-based greeting (first turn only)
    'experimental.chat.system.transform': async (input, output) => {
      const persona = readFileSafe(ACTIVE_PERSONA);
      if (persona) {
        (output.system ||= []).push(
          `<persona>\n${persona}\n</persona>`
        );
      }

      // Only inject greeting on the first turn of each session
      const sessionID = input.sessionID || 'default';
      if (!greetedSessions.has(sessionID)) {
        greetedSessions.add(sessionID);
        (output.system ||= []).push(
          `<session-greeting>\n${getGreeting()}\n</session-greeting>`
        );
      }
    },

    // 2. Mood parsing + voice playback on events
    event: async ({ event }) => {
      // Mood parsing — accumulate streamed deltas, parse from buffer
      if (event.type === 'message.part.updated') {
        const part = event.properties?.part;
        if (part?.type === 'text') {
          const partId = part.id || 'default';
          const existing = textBuffers.get(partId) || '';
          const updated = existing + (event.properties?.delta || '');
          textBuffers.set(partId, updated);
          const mood = parseMood(updated);
          if (mood) writeMood(mood);
        }
      }

      // Clean up text buffers when session goes idle (response complete)
      if (event.type === 'session.idle') {
        textBuffers.clear();
      }

      // Voice playback for mapped events
      const voiceDir = EVENT_VOICE_MAP[event.type];
      if (voiceDir) {
        await playRandomVoice($, voiceDir);
      }
    },

    // 3. Voice on user prompt submission
    'chat.message': async () => {
      await playRandomVoice($, 'UserPromptSubmit');
    },
  };
};
```

- [ ] **Step 2: Verify JS syntax**

Run: `node --check runtimes/opencode/plugins/maid-cafe.js`
Expected: No syntax errors

- [ ] **Step 3: Commit**

```bash
git add runtimes/opencode/plugins/maid-cafe.js
git commit -m "feat(opencode): add maid-cafe plugin with persona/mood/greeting/voice"
```

---

## Chunk 3: Install Script

### Task 5: Update `install.sh`

**Files:**
- Modify: `install.sh:48-75` (the `install_opencode` function)

- [ ] **Step 1: Replace hardcoded plugin symlink with glob**

In `install.sh`, change the plugin section (lines 68-71) from:

```bash
  # Plugin (system prompt injection)
  mkdir -p "$HOME/.config/opencode/plugins"
  ln -sf "$SCRIPT_DIR/runtimes/opencode/plugins/agent-marketplace.js" \
    "$HOME/.config/opencode/plugins/agent-marketplace.js"
```

To:

```bash
  # Plugins
  mkdir -p "$HOME/.config/opencode/plugins"
  for plugin_js in "$SCRIPT_DIR"/runtimes/opencode/plugins/*.js; do
    [ -f "$plugin_js" ] && ln -sf "$plugin_js" \
      "$HOME/.config/opencode/plugins/$(basename "$plugin_js")"
  done
```

- [ ] **Step 2: Add maid-cafe assets installation**

Add a new function `install_maid_cafe` after `install_opencode`:

```bash
# Maid-cafe: personas + active persona + voice files
install_maid_cafe() {
  local maid_dir="$HOME/.config/opencode/maid-cafe"

  # Personas (symlink from CC plugin source)
  shopt -s nullglob
  mkdir -p "$maid_dir/personas"
  for persona in "$SCRIPT_DIR"/plugins/maid-cafe/personas/*.md; do
    ln -sf "$persona" "$maid_dir/personas/$(basename "$persona")"
  done
  shopt -u nullglob

  # Default active persona (no-clobber)
  if [ ! -f "$maid_dir/active-persona.md" ]; then
    cp "$SCRIPT_DIR/plugins/maid-cafe/personas/codex.md" \
      "$maid_dir/active-persona.md"
  fi

  # Voice files
  local voice_zip="$SCRIPT_DIR/plugins/maid-cafe/assets/maid-voice.zip"
  if [ ! -d "$maid_dir/maid-voice" ] && [ -f "$voice_zip" ]; then
    mkdir -p "$maid_dir"
    unzip -qo "$voice_zip" -d "$maid_dir" 2>/dev/null || true
  fi

  echo "  Maid-cafe: personas + voice installed"
}
```

- [ ] **Step 3: Call `install_maid_cafe` in main body**

Add after the `install_opencode` call (around line 85):

```bash
install_maid_cafe || echo "  Maid-cafe: skipped (error during install)"
```

- [ ] **Step 4: Verify script syntax**

Run: `bash -n install.sh`
Expected: No syntax errors

- [ ] **Step 5: Commit**

```bash
git add install.sh
git commit -m "feat(install): add maid-cafe assets + glob plugin symlinks for OpenCode"
```

---

## Chunk 4: Verify

### Task 6: Structure verification

- [ ] **Step 1: Verify all new files exist**

Run: `fd -t f . runtimes/opencode/ | sort`
Expected output should include:
```
runtimes/opencode/agents/maid.md
runtimes/opencode/commands/look.md
runtimes/opencode/commands/switch.md
runtimes/opencode/plugins/agent-marketplace.js
runtimes/opencode/plugins/maid-cafe.js
```

- [ ] **Step 2: Verify JS module syntax**

Run: `node --check runtimes/opencode/plugins/maid-cafe.js`
Expected: No output (success)

- [ ] **Step 3: Verify install.sh syntax**

Run: `bash -n install.sh`
Expected: No output (success)

- [ ] **Step 4: Dry-run install.sh to check maid-cafe function**

Run: `bash -x install.sh 2>&1 | grep -i maid`
Expected: Should show maid-cafe installation steps executing

- [ ] **Step 5: Verify installed layout**

Run: `ls -la ~/.config/opencode/maid-cafe/personas/ && ls ~/.config/opencode/plugins/maid-cafe.js`
Expected: Persona symlinks and plugin symlink exist

- [ ] **Step 6: Final commit (if any fixes needed)**

```bash
git add -A
git commit -m "fix(opencode): address verification issues in maid-cafe runtime"
```
