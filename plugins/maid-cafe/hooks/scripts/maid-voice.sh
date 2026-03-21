#!/bin/bash
# Play a random WAV file for the given hook event.
# Usage: maid-voice.sh <EventName>
# WAV files expected in: ~/.claude/hooks/maid-voice/<EventName>/*.wav
cat >/dev/null

EVENT="${1:-}"
if [ -z "$EVENT" ]; then
  exit 0
fi

if ! command -v afplay &>/dev/null; then
  exit 0
fi

VOICE_DIR="$HOME/.claude/hooks/maid-voice/$EVENT"
if [ ! -d "$VOICE_DIR" ]; then
  exit 0
fi

shopt -s nullglob
FILES=("$VOICE_DIR"/*.wav)
shopt -u nullglob

if [ ${#FILES[@]} -eq 0 ]; then
  exit 0
fi

RANDOM_INDEX=$((RANDOM % ${#FILES[@]}))
SELECTED="${FILES[$RANDOM_INDEX]}"

nohup afplay "$SELECTED" &>/dev/null &

exit 0
