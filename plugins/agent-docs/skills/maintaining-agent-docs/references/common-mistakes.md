# Common Mistakes

### Mistake 1: Creating Docs Without an Index
**Problem:** Engineers find system/architecture.md, then search for database schema for 20 minutes

**Fix:** Always update .agent/README.md when adding documentation. It's your wayfinding system.

### Mistake 2: "Code is the Documentation"
**Problem:** New engineer reads code for 3 hours to understand a 5-minute concept

**Reality Check:**
- Code shows HOW (implementation details)
- Docs show WHY (design decisions, tradeoffs)
- Both are needed

### Mistake 3: Deferring Documentation ("We'll Update It Later")
**Pattern:** "Once we're done coding, we'll write docs"

**Reality:** You won't. Use these rules instead:
- Update docs DURING implementation (as you learn)
- Each PR should include doc updates
- Put "update .agent docs" in your DoD (Definition of Done)

### Mistake 4: Orphaned Documents
**Problem:** Someone creates SOP/add-new-endpoint.md but forgets to add it to .agent/README.md

**Result:**
- Document exists but is undiscoverable
- Duplicate documentation (someone writes it again)
- Knowledge that seems to exist but can't be found

**Fix:** Every new document MUST be added to its category's README (Tasks/README.md, System/README.md, or SOP/README.md) AND updated in .agent/README.md

### Mistake 5: Mixing Narrative with Reference
**Problem:** Architecture doc mixes "here's how it evolved" with "here's what it looks like now"

**Fix:** Keep docs focused:
- **System/** = current state only
- **Tasks/** = decisions and what we learned
- **SOP/** = procedures only

### Mistake 6: Bloated CLAUDE.md
**Problem:** CLAUDE.md grows to 100+ lines with full architecture details, complete code style guides, API documentation, etc.

**Why it matters:** CLAUDE.md is loaded into the context window on every conversation, wasting tokens on content that should be retrieved on demand.

**Research confirms:** An 8KB compressed index performs equivalently to a 40KB full document (Vercel study).

**Fix:** CLAUDE.md contains only the index. Use pipe-delimited format to list `.agent/` file paths. Agents will read full content on demand.

### Mistake 7: Duplicating User Scope Rules in Project Scope
**Problem:** Every project's CLAUDE.md copy-pastes the same rules — "use Traditional Chinese", "use bun not npm", "use ruff for linting" — creating N copies of identical content.

**Why it matters:** Maintaining N copies means changing one rule requires updating N files. Miss one and you have inconsistency.

**Fix:** Before creating project CLAUDE.md, read `~/.claude/CLAUDE.md`. Rules already covered there are NOT repeated. Project CLAUDE.md only contains project-specific content, but its commands and conventions should reflect user scope preferences (e.g., if user scope says use bun, Commands section uses `bun test` not `npm test`).

### Mistake 8: Manually Creating CLAUDE.md Instead of Using `/init`
**Problem:** Handwriting `CLAUDE.md` from scratch misses project-specific conventions that `/init` would auto-detect (build system, test framework, linting tools, project structure).

**Why it matters:** `/init` analyzes the actual project and generates tailored Commands and Rules. Manual creation relies on the agent's assumptions, leading to incomplete or generic content.

**Fix:** Always use `/init` to generate the base `CLAUDE.md`, then supplement with the `## Docs Index` section (which `/init` doesn't produce). Direct Edit is fine for subsequent updates, but initial creation must go through `/init`.
