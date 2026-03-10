#!/bin/bash
# Auto-build Go tools for this plugin if binary is missing
PLUGIN_ROOT="${CLAUDE_PLUGIN_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)}"
for tool_dir in "$PLUGIN_ROOT"/tools/*/; do
  [ -f "$tool_dir/Makefile" ] || continue
  binary=$(basename "$tool_dir")
  if ! command -v "$binary" &>/dev/null; then
    command -v go &>/dev/null && make -C "$tool_dir" install 2>/dev/null || true
  fi
done
exit 0
