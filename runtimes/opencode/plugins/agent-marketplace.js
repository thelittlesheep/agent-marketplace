/**
 * agent-marketplace plugin for OpenCode.ai
 *
 * Injects available skills listing into system prompt so the agent
 * knows what agent-marketplace skills are available.
 * Skills are discovered via OpenCode's native skill tool from symlinked directories.
 */

import path from 'path';
import fs from 'fs';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

// Walk up from plugins/ -> opencode/ -> runtimes/ -> agent-marketplace/
const marketplaceRoot = path.resolve(__dirname, '../../..');

// Simple frontmatter extraction — only handles single-line key: value pairs.
// Multi-line YAML values (| or >) are not supported. This matches superpowers' approach.
const extractFrontmatter = (content) => {
  const match = content.match(/^---\n([\s\S]*?)\n---\n([\s\S]*)$/);
  if (!match) return { frontmatter: {}, content };

  const frontmatterStr = match[1];
  const body = match[2];
  const frontmatter = {};

  for (const line of frontmatterStr.split('\n')) {
    const colonIdx = line.indexOf(':');
    if (colonIdx > 0) {
      const key = line.slice(0, colonIdx).trim();
      const value = line.slice(colonIdx + 1).trim().replace(/^["']|["']$/g, '');
      frontmatter[key] = value;
    }
  }

  return { frontmatter, content: body };
};

const discoverSkills = () => {
  const pluginsDir = path.join(marketplaceRoot, 'plugins');
  const skills = [];

  try {
    const plugins = fs.readdirSync(pluginsDir, { withFileTypes: true });
    for (const plugin of plugins) {
      if (!plugin.isDirectory()) continue;
      const skillsDir = path.join(pluginsDir, plugin.name, 'skills');
      if (!fs.existsSync(skillsDir)) continue;

      const skillDirs = fs.readdirSync(skillsDir, { withFileTypes: true });
      for (const skillDir of skillDirs) {
        if (!skillDir.isDirectory()) continue;
        const skillMd = path.join(skillsDir, skillDir.name, 'SKILL.md');
        if (!fs.existsSync(skillMd)) continue;

        const content = fs.readFileSync(skillMd, 'utf8').replace(/\r\n/g, '\n');
        const { frontmatter } = extractFrontmatter(content);
        if (frontmatter.name && frontmatter.description) {
          skills.push({ name: frontmatter.name, description: frontmatter.description.trim() });
        }
      }
    }
  } catch {
    // Silently ignore if plugins directory doesn't exist
  }

  return skills;
};

export const AgentMarketplacePlugin = async () => {
  const skills = discoverSkills();
  if (skills.length === 0) return {};

  const skillList = skills.map(s => `- **${s.name}**: ${s.description}`).join('\n');

  return {
    'experimental.chat.system.transform': async (_input, output) => {
      (output.system ||= []).push(
        `<agent-marketplace>\nThe following agent-marketplace skills are available. Use the skill tool to load any of them when relevant:\n${skillList}\n</agent-marketplace>`
      );
    }
  };
};
