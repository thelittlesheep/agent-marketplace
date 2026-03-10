#!/bin/bash
# 驗證 code review CLI 工具的可用性
# Exit code 2 = 阻止操作並回報錯誤給 Claude

command -v jq &>/dev/null || exit 0

INPUT=$(cat)
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // empty' 2>/dev/null)

# 檢查是否是 review 相關命令（用 word boundary 避免 false positive）
if [[ "$COMMAND" =~ (^|[[:space:]])gemini([[:space:]]|$) ]]; then
  if ! command -v gemini &>/dev/null; then
    echo "錯誤：gemini CLI 未安裝。請執行 'bun add -g @google/generative-ai-cli'" >&2
    exit 2
  fi
fi

if [[ "$COMMAND" =~ (^|[[:space:]])opencode([[:space:]]|$) ]]; then
  if ! command -v opencode &>/dev/null; then
    echo "錯誤：opencode CLI 未安裝。請參考 https://opencode.ai 安裝" >&2
    exit 2
  fi
fi

exit 0
