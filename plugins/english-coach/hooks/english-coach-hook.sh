#!/bin/bash
INPUT=$(cat)
TRANSCRIPT=$(echo "$INPUT" | jq -r '.transcript_path // ""')
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // ""')

if [ -n "$TRANSCRIPT" ] && [ "$TRANSCRIPT" != "null" ]; then
  nohup english-coach \
    "$TRANSCRIPT" "$SESSION_ID" &>/dev/null &
fi
# Exit 0 with no stdout = allow tool execution
