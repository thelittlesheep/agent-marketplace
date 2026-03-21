#!/bin/bash
set -euo pipefail

INPUT=$(cat)
MOOD_FILE="$HOME/.claude/mood.txt"

# Extract transcript path from hook JSON
TRANSCRIPT=$(echo "$INPUT" | jq -r '.transcript_path // ""' 2>/dev/null)

if [ -z "$TRANSCRIPT" ] || [ "$TRANSCRIPT" = "null" ] || [ ! -f "$TRANSCRIPT" ]; then
  exit 0
fi

# Parse mood marker directly from all assistant text in transcript
# Combines extraction + parsing to avoid tail -1 empty-line bug when text ends with \n
MOOD=$(tail -50 "$TRANSCRIPT" | jq -r 'select(.type == "assistant") | .message.content[] | select(.type == "text") | .text' 2>/dev/null | sed -n 's/.*【[[:space:]]*\(.*\)[[:space:]]*】.*/\1/p' | tail -1)

if [ -n "$MOOD" ]; then
  mkdir -p "$(dirname "$MOOD_FILE")"
  echo "$MOOD" >"$MOOD_FILE"
fi

exit 0
