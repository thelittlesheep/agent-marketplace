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
