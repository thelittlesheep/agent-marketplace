#!/bin/bash
set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Build Go tools -> ~/.local/bin/
build_tools() {
  if ! command -v go &>/dev/null; then
    echo "Go not found. Skipping tool builds."
    return
  fi
  for tool_dir in "$SCRIPT_DIR"/plugins/*/tools/*/; do
    if [ -f "$tool_dir/Makefile" ]; then
      echo "Building $(basename "$tool_dir")..."
      make -C "$tool_dir" install 2>&1 | tail -1
    fi
  done
}

# Install default config (no-clobber)
install_config() {
  local cfg_dir="$HOME/.config/review-cli"
  # Remove dangling symlink if target doesn't exist
  if [ -L "$cfg_dir" ] && [ ! -e "$cfg_dir" ]; then
    rm "$cfg_dir"
  fi
  mkdir -p "$cfg_dir"
  if [ ! -f "$cfg_dir/config.toml" ]; then
    cp "$SCRIPT_DIR/plugins/code-review-router/config/review-cli/config.toml" \
      "$cfg_dir/config.toml"
  fi
}

# OpenCode: symlink skills + commands + agents + plugin
install_opencode() {
  shopt -s nullglob
  # Skills
  mkdir -p "$HOME/.config/opencode/skills"
  for skill_dir in "$SCRIPT_DIR"/plugins/*/skills/*/; do
    ln -sf "$skill_dir" "$HOME/.config/opencode/skills/$(basename "$skill_dir")"
  done

  # Commands (.md format with $ARGUMENTS)
  mkdir -p "$HOME/.config/opencode/commands"
  for cmd in "$SCRIPT_DIR"/runtimes/opencode/commands/*.md; do
    [ -f "$cmd" ] && ln -sf "$cmd" "$HOME/.config/opencode/commands/$(basename "$cmd")"
  done

  # Agents
  mkdir -p "$HOME/.config/opencode/agents"
  for agent in "$SCRIPT_DIR"/runtimes/opencode/agents/*.md; do
    [ -f "$agent" ] && ln -sf "$agent" "$HOME/.config/opencode/agents/$(basename "$agent")"
  done

  # Plugins
  mkdir -p "$HOME/.config/opencode/plugins"
  for plugin_js in "$SCRIPT_DIR"/runtimes/opencode/plugins/*.js; do
    [ -f "$plugin_js" ] && ln -sf "$plugin_js" \
      "$HOME/.config/opencode/plugins/$(basename "$plugin_js")"
  done

  shopt -u nullglob
  echo "  OpenCode: skills + commands + agents + plugin installed"
}

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

echo "Installing agent-marketplace for OpenCode..."
echo ""
build_tools || echo "  Warning: some tools failed to build"
install_config || echo "  Warning: config install failed"
echo ""

echo "Linking OpenCode assets..."
install_opencode || echo "  OpenCode: skipped (error during install)"
install_maid_cafe || echo "  Maid-cafe: skipped (error during install)"
echo ""

echo "Done!"
echo ""
echo "For Claude Code, use the marketplace system instead:"
echo "  claude plugin marketplace add $SCRIPT_DIR"
echo "  claude plugin install <plugin-name>"
