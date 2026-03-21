#!/bin/bash
# One-time setup for maid-cafe plugin.
# Runs on SessionStart; skips if already configured.
set -euo pipefail
cat >/dev/null

PLUGIN_ROOT="${CLAUDE_PLUGIN_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)}"
SETTINGS="$HOME/.claude/settings.json"
VOICE_DIR="$HOME/.claude/hooks/maid-voice"
VOICE_ZIP="$PLUGIN_ROOT/assets/maid-voice.zip"

# --- 1. spinnerVerbs ---
if command -v jq &>/dev/null && [ -f "$SETTINGS" ]; then
  if ! jq -e '.spinnerVerbs' "$SETTINGS" &>/dev/null; then
    VERBS=$(jq -r '.spinnerVerbs' "$PLUGIN_ROOT/references/spinner-verbs.json")
    jq --argjson verbs "$VERBS" '. + {spinnerVerbs: $verbs}' "$SETTINGS" >"$SETTINGS.tmp" \
      && mv "$SETTINGS.tmp" "$SETTINGS"
  fi
elif command -v jq &>/dev/null && [ ! -f "$SETTINGS" ]; then
  mkdir -p "$(dirname "$SETTINGS")"
  jq '{spinnerVerbs: .spinnerVerbs}' "$PLUGIN_ROOT/references/spinner-verbs.json" >"$SETTINGS"
fi

# --- 2. Voice files ---
if [ ! -d "$VOICE_DIR" ] && [ -f "$VOICE_ZIP" ]; then
  if command -v unzip &>/dev/null; then
    mkdir -p "$VOICE_DIR"
    unzip -qo "$VOICE_ZIP" -d "$HOME/.claude/hooks" 2>/dev/null || true
  fi
fi

exit 0
