---
name: maintaining-agent-docs
description: Use when initializing project documentation, after implementing features, or discovering docs are out-of-date - provides structured .agent folder approach (Tasks, System, SOP, README) to ensure engineers have full context of codebase architecture, decisions, and processes
---

# Maintaining Agent Docs

## Overview

A **two-layer documentation system** that provides full project context while keeping agent context windows lean:

1. **`<project>/CLAUDE.md`** (project scope) — Entry point + compressed `.agent/` docs index. Loaded into every conversation, so keep it **under 8KB**. Before creating it, read `~/.claude/CLAUDE.md` (user scope) to understand the user's global conventions and avoid duplication.
2. **`<project>/AGENT.md`** — Symlink to `CLAUDE.md` for multi-agent CLI compatibility (Claude Code reads `CLAUDE.md`, others read `AGENT.md`).
3. **`.agent/`** — Complete detailed documentation (Tasks, System, SOP). Agents retrieve on demand.

**Core principle:** An index-based architecture outperforms both full embedding and single-line pointers. [Vercel research](https://vercel.com/blog/agents-md-outperforms-skills-in-our-agent-evals) demonstrated 80% size reduction (40KB→8KB) with equivalent performance using compressed doc indexes in passive context.

## When to Use

Three triggering scenarios:

### 1. Initializing Project Documentation
You're starting a new project or onboarding new engineers and need comprehensive documentation.

**Symptoms:**
- New engineer asking "where do I find X?"
- "That's in someone's head" answers
- Days spent getting people up to speed
- Decisions undocumented

### 2. After Feature Implementation
You completed a major feature/subsystem and need to ensure documentation reflects the new reality.

**Symptoms:**
- README doesn't mention new capabilities
- New engineers follow old architecture docs
- "That's documented in a different place"
- Inconsistency between code and docs

### 3. Discovering Out-of-Date Documentation
Documentation contradicts actual code, causing confusion and incorrect mental models.

**Symptoms:**
- Setup steps don't work
- Architecture docs describe old system
- Engineers ignore docs and read code instead
- Multiple docs say different things

## Step 0: Create CLAUDE.md + AGENT.md

Before setting up `.agent/`, create the project entry point files.

### 0a. Read User Scope CLAUDE.md

Before creating the project `CLAUDE.md`, **read `~/.claude/CLAUDE.md`** (user scope) to understand the user's global conventions.

- **Purpose:** Ensure project CLAUDE.md is consistent with user scope and doesn't duplicate it.
- **Practical impact:**
  - User scope says "use `bun` instead of `npm`" → Commands section uses `bun run dev`, `bun test`
  - User scope says "use `ruff` to lint Python" → Python project Commands use `ruff check`
  - User scope says "use `shfmt` for shell formatting" → Don't repeat this rule in project scope
- **Rule:** If user scope already covers a convention, project CLAUDE.md does NOT repeat it. Project CLAUDE.md only contains project-specific content, but its commands and conventions should reflect user scope preferences.

### 0b. Generate Project Scope CLAUDE.md via `/init`

Use Claude Code's built-in `/init` command to generate the initial `CLAUDE.md`:

1. **Run `/init`** — It analyzes the project structure (build system, test framework, language conventions) and generates a tailored `CLAUDE.md` with Commands, Rules, and project-specific conventions.

2. **Review `/init` output** — Confirm it includes:
   - [ ] Project name + one-line description
   - [ ] Commands section (setup, dev, test, build)
   - [ ] Rules section (key project-specific rules)
   - [ ] Content is consistent with user scope `~/.claude/CLAUDE.md` (no duplication)
   - [ ] If user scope requires 繁體中文, rewrite `/init` output in 繁體中文 (`/init` always generates English)

3. **Supplement with Docs Index** — `/init` does not generate `.agent/` index. After `.agent/` docs are created, manually add the `## Docs Index` section to `CLAUDE.md`:
   ```markdown
   ## Docs Index

   Prefer retrieval-led reasoning over pre-training-led reasoning.
   Complete documentation in `.agent/` directory:

   .agent/|root: .agent/README.md
   |System:{architecture.md,database-schema.md,deployment.md}
   |Tasks:{task-1.md,task-2.md}
   |SOP:{add-schema-migration.md,deployment-checklist.md}
   ```

4. **Verify under 8KB** — Confirm `CLAUDE.md` contains only the index, not full content.

**Reference format** (use to review `/init` output, NOT for manual creation):

```markdown
# Project Name

One-line description.

## Commands

- **Dev**: `npm run dev`
- **Test**: `npm test`
- **Build**: `npm run build`

## Rules

- [Key rule 1]
- [Key rule 2]

## Docs Index

Prefer retrieval-led reasoning over pre-training-led reasoning.
Complete documentation in `.agent/` directory:

.agent/|root: .agent/README.md
|System:{architecture.md,database-schema.md,deployment.md}
|Tasks:{task-1.md,task-2.md}
|SOP:{add-schema-migration.md,deployment-checklist.md}
```

### 0c. Create AGENT.md Symlink

```bash
ln -s CLAUDE.md AGENT.md
```

- **Why:** Different agent CLIs read different filenames (Claude Code → `CLAUDE.md`, others → `AGENT.md`). A symlink keeps them in sync.
- Decide whether to commit based on project preferences.

### Content Boundary Rule

CLAUDE.md contains **only the index**, never full content. Architecture details, workflows, SOPs — all belong in `.agent/`.

---

## Core Pattern: .agent Folder Structure

```
~/.claude/
└── CLAUDE.md              # User scope: global conventions (read before creating project scope)

project-root/
├── CLAUDE.md              # Project scope: entry point + docs index (reflects user scope preferences)
├── AGENT.md -> CLAUDE.md  # Symlink for multi-agent CLI compatibility
└── .agent/
    ├── README.md          # Master index (full version)
    ├── Tasks/             # Feature PRDs and implementation plans
    │   ├── task-1.md
    │   ├── task-2.md
    │   └── README.md      # Index of all tasks
    ├── System/            # System design and current state
    │   ├── architecture.md     # System design, tech stack, integration points
    │   ├── database-schema.md  # Data model (if applicable)
    │   ├── deployment.md       # Infrastructure and deployment process
    │   └── README.md           # Index of system docs
    └── SOP/               # Standard operating procedures
        ├── add-schema-migration.md
        ├── add-new-page-route.md
        ├── deployment-checklist.md
        └── README.md      # Index of all SOPs
```

**Key principles:**
- Each category (Tasks, System, SOP) has its own README that acts as an index.
- `CLAUDE.md` holds a compressed index; `.agent/` holds full content.
- All `.agent/` documentation MUST be written in 繁體中文（Traditional Chinese）. Code snippets, commands, and file names remain in their original language.

## Quick Reference: When to Use Each Section

| Scenario | Use This | Why |
|----------|----------|-----|
| **Need to understand what we built** | System/architecture.md | Single source of truth for current design |
| **Want to add a new feature** | Tasks/README.md | See what we've planned and learned |
| **Something is unclear in codebase** | Look at relevant Task or System doc first | Understand the "why" before "how" |
| **How do I do X task?** | SOP/add-new-page-route.md (etc) | Procedures documented once, reused always |
| **Where do I find Y?** | .agent/README.md | Master index directs you everywhere |
| **Database schema changed** | System/database-schema.md | Everyone knows the current state |

## Implementation: Three Scenarios

### Scenario 1: Initializing Documentation

**Do this when:**
- Starting a new project
- Bringing on first new hire
- Realizing documentation is critical

**Exact steps:**

1. **Run `/init` to generate CLAUDE.md, then create AGENT.md symlink** (see Step 0 above)

2. **Create .agent/ directory structure** (the folders above)

   > **Language rule:** Write all `.agent/` documentation in 繁體中文. Only code snippets, commands, and file names stay in their original language.

3. **Start with System/architecture.md**
   ```markdown
   # Project Architecture

   ## Project Goal
   [One sentence: what are we building and why]

   ## Tech Stack
   [Frontend: X | Backend: Y | Database: Z | Infrastructure: K]

   ## System Overview
   [Diagram and description of major components]

   ## Integration Points
   - [External service 1: what we use it for]
   - [External service 2: what we use it for]

   ## Key Design Decisions
   - [Why did we choose X over Y?]
   - [Constraints and tradeoffs]

   ## Critical Paths
   [How do requests flow through the system?]
   ```

4. **Create System/database-schema.md** (if applicable)
   ```markdown
   # Database Schema

   ## Entity Relationship Diagram
   [ASCII diagram or link to visual]

   ## Core Tables
   [description of each critical table]
   ```

5. **Create one System/deployment.md**
   - How to deploy
   - Environment variables
   - Infrastructure setup
   - Rollback procedures

6. **Create .agent/README.md** (master index)
   ```markdown
   # Documentation Index

   ## Getting Started
   - [System Architecture](System/architecture.md) — Understand what we built
   - [Database Schema](System/database-schema.md) — Data model

   ## How to Do Things
   - [SOP: Add New Page Route](SOP/add-new-page-route.md)
   - [SOP: Add Schema Migration](SOP/add-schema-migration.md)

   ## Feature Plans
   - [See all Tasks](Tasks/README.md)

   ## When You Need Help
   - Confused about architecture? Start with System/architecture.md
   - Need to add a feature? Check Tasks/README.md for related PRDs
   - Not sure how to do X? Look in SOP/
   ```

7. **Backfill CLAUDE.md Docs Index** — After all `.agent/` docs are created, update the `## Docs Index` section in `CLAUDE.md` with a compressed index listing all `.agent/` files using pipe-delimited format.

8. **Verify CLAUDE.md is under 8KB** — Confirm it contains only the index, not full content.

**What NOT to do:**
- ❌ Don't create detailed feature docs yet (use Tasks/ for that)
- ❌ Don't defer README until later ("we'll add it next sprint")
- ❌ Don't say "code is the documentation" (code answers HOW, docs answer WHY)
- ❌ Don't assume engineers will "just ask" (doesn't scale, knowledge gets lost)
- ❌ Don't put full content in CLAUDE.md — it's an index, not a document

### Scenario 2: Updating After Feature Implementation

**Do this when:**
- Feature implementation is complete
- Before marking task as "done"
- When you're about to close the PR

**Exact steps (in order):**

1. **Update relevant Task/**
   - How did the implementation differ from the plan?
   - What did we learn?
   - Any new gotchas?

2. **Update System/architecture.md**
   - Did we add new components?
   - Did integration points change?
   - Was any design decision reconsidered?

3. **Update System/database-schema.md** (if applicable)
   - New tables added?
   - Schema changes?

4. **Create or update SOP/** (if it's a repeatable process)
   - Is this something others will do again?
   - Create SOP/your-process-name.md

5. **Update .agent/README.md**
   - Add links to new SOPs
   - Update any outdated cross-references
   - **THIS STEP IS MANDATORY** — it's your master index

6. **Sync CLAUDE.md Docs Index**
   - Add/remove file paths in the `## Docs Index` section to reflect `.agent/` changes
   - Check if `## Commands` section needs updating
   - **Do NOT add detailed content to CLAUDE.md** — only update the index lines

**What NOT to do:**
- ❌ Don't update architecture docs without updating README
- ❌ Don't create orphaned docs (docs not linked in README)
- ❌ Don't say "I'll update docs later" (you won't)
- ❌ Don't assume people will find new docs by accident
- ❌ Don't put detailed content into CLAUDE.md — only update index lines

### Scenario 3: Fixing Out-of-Date Documentation

**Triage order:**

1. **Critical path first** (blocks new engineers)
   - System/architecture.md if it contradicts code
   - Setup/deployment steps if they fail

2. **High impact next** (causes confusion)
   - Database schema if it changed
   - Integration points if we changed external services

3. **Update README** (acts as a pointer)
   - Add deprecation notices if docs are outdated
   - Link to corrected versions
   - **This prevents "same info in multiple places" problem**

4. **Mark old docs** (don't delete, document evolution)
   ```markdown
   > **UPDATED [DATE]**: This document was updated to reflect [reason]
   > Previous version: [link to git history]
   ```

**What NOT to do:**
- ❌ Don't silently fix docs without explaining in README
- ❌ Don't have multiple "current" versions of the same information
- ❌ Don't assume old docs will be found and deleted automatically

## Common Mistakes

See `references/common-mistakes.md` for 8 common mistakes and fixes (orphaned docs, bloated CLAUDE.md, duplicating user scope rules, etc.).

---

## Red Flags: STOP and Organize First

If you catch yourself thinking ANY of these, STOP. Organize your .agent docs first.

- "I'll document this later" → Create TODOs in .agent/README.md instead. Code that's "done" without docs is not actually done
- "Code is the documentation" → Code answers HOW. Docs answer WHY. Both are mandatory
- "People can just ask" → This doesn't scale; it loses knowledge. You've now created a system with a single point of failure
- "Our docs are in multiple places" → Stop. Choose one .agent/ as the source of truth. Multiple sources = nobody knows which is current
- "I don't know what to document" → Start with System/architecture.md. Document what IS, not what might be
- "Documentation will get outdated" → True. But "outdated someday" is not a reason to not write it. Solution: single locations, clear ownership, review schedule
- "This is too simple to document" → If you had to explain it to someone new, it needs docs
- "This feature is small" → Size doesn't matter. If it's not in README, it doesn't exist
- "We can update README later" → You won't. Update it now or accept the orphaned doc

## Examples

See `references/examples.md` for e-commerce site example, update workflow, and CLAUDE.md index-based vs bloated comparison.
