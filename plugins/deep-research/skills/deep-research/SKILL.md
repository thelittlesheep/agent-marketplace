---
name: deep-research
description: "Use ONLY when explicitly triggered via /deep-research command -- NEVER auto-activate. Performs source-first adaptive research: discovers authoritative sources BEFORE answering, deep-reads them, synthesizes with inline citations, and self-audits for unverified claims. Supports three modes: Research (technical solutions), Verification (fact-checking), and Audit (verify artifacts against trusted sources). This skill operationalizes the principle: prefer retrieval-led reasoning over pre-training-led reasoning."
---

# Deep Research

Source-first adaptive research pipeline. Discover sources BEFORE answering. Never rely on pre-trained knowledge alone.

**Activation**: Manual only via `/deep-research`. Never self-activate.

## Phase 0: Query Analysis

Classify the query before doing anything else.

### Mode Detection

| Signal | Mode | What changes |
|--------|------|-------------|
| `@` file paths provided | **Audit** | Read artifacts, extract claims, verify each one |
| "is it true", "verify", "fact-check", "does X really", "I heard that" | **Verification** | Need 2+ independent sources confirming/denying; actively search for contradicting evidence |
| "how to", "best approach", "compare", "what is", "how does X work" | **Research** | Broad source collection; synthesize multiple perspectives |
| Ambiguous | **Research** (default) | When in doubt, research is safer than verification |

### Scope Tagging

Tag the query with ALL applicable labels -- these determine source priority in Phase 1:

- `[library]` -- involves a specific library/framework (e.g., "React hooks", "SQLAlchemy session")
- `[api]` -- involves a specific API or service
- `[architecture]` -- system design, patterns, tradeoffs
- `[factual]` -- verifiable fact (release date, feature existence, compatibility)
- `[opinion]` -- inherently subjective (best practice debates, tool comparisons)
- `[current]` -- time-sensitive (latest version, recent changes, deprecations)

Output your classification before proceeding:

```
Mode: Research | Verification | Audit
Tags: [library][current]
Query restated: "What is the recommended way to handle authentication in Next.js 15 App Router?"
```

For **Audit mode**, also list the files to be examined.

---

## Research / Verification Mode (no files)

### Phase 1: Source Discovery

#### Source Priority Hierarchy

See `references/source-priority.md` for the full hierarchy with domain lists. Summary:

| Priority | Source type | Tool | When to use |
|----------|-----------|------|-------------|
| P0 | Library docs (context7) | `mcp__context7__resolve-library-id` -> `mcp__context7__query-docs` | Query has `[library]` tag |
| P1 | Official docs (fetched) | `WebFetch` on known official doc URLs | Query has `[api]` or `[library]` tag |
| P2 | Web search (authoritative) | `WebSearch` | Always |
| P3 | Deep-read top results | `WebFetch` on search result URLs | Always (top 3-5 results) |
| P4 | Browser automation | `agent-browser` | Only when P0-P3 all fail or return insufficient results |

#### Step 1: context7 (if `[library]` tagged)

```
1. mcp__context7__resolve-library-id -> get library ID
2. mcp__context7__query-docs -> query with specific question
```

Record what context7 returns. Even if it answers the question fully, proceed to Step 2 for cross-verification -- especially for `[current]` tagged queries where docs may lag behind reality.

#### Step 2: WebSearch (always)

Run 2-3 search queries with different phrasings:

- Query 1: The original question, verbatim
- Query 2: Rephrased with key technical terms (add version numbers, official terminology)
- Query 3 (if `[current]`): Add current year to query

Example for "Next.js 15 authentication":
- Q1: "Next.js 15 App Router authentication best practice"
- Q2: "next.js 15 auth.js middleware session management"
- Q3: "Next.js 15 authentication 2026"

#### Step 3: Deep-read top results

From search results, select 3-5 URLs to deep-read via `WebFetch`.

**Selection criteria** (in order):
1. Official documentation domains (see `references/source-priority.md`)
2. GitHub repos / official blogs
3. Well-known technical publications
4. Recent content (prefer last 12 months for `[current]` queries)
5. Stack Overflow answers with high vote counts

**Skip**: SEO-farm content, AI-generated aggregator sites, outdated (2+ years for `[current]`).

#### Step 4: agent-browser (fallback only)

Use ONLY when:
- All P0-P3 sources returned nothing useful
- Content requires JavaScript rendering that WebFetch cannot handle
- Content is behind a login that agent-browser can access

#### Source Minimums

| Mode | Minimum sources | Recommended |
|------|----------------|-------------|
| Research | 3 | 5-8 |
| Verification | 2 independent | 3-4 independent |

"Independent" means sources that don't cite each other as their primary reference.

**If you cannot reach the minimum**: State this explicitly in your answer. Do not silently fill gaps with pre-trained knowledge.

### Phase 2: Synthesis

#### Answer Structure

```markdown
## [Direct answer to the query -- 1-3 sentences]

[Detailed explanation with inline citations]

Each factual claim MUST have a citation: "React Server Components
render on the server by default [source](https://react.dev/...)."

### [Subtopic 1]
[Content with citations...]

### [Subtopic 2]
[Content with citations...]

### Caveats
- [Important limitations, version-specific notes, common pitfalls]

### Sources Consulted
| # | Source | Type | URL |
|---|--------|------|-----|
| 1 | React Documentation | Official docs | [link](url) |
| 2 | Vercel Blog | Official blog | [link](url) |
```

#### Citation Rules

1. **Every factual claim needs a citation** -- no exceptions for "common knowledge"
2. **Inline format**: `text [source](url)` -- immediately after the claim
3. **Same source, multiple claims**: repeat the full `[source](url)` each time
4. **No footnote-style citations** -- the reader should never need to scroll
5. **Code examples**: cite the source above the code block, not inside it

#### Confidence Marking

| Marker | Meaning | When to use |
|--------|---------|-------------|
| *(no mark)* | Confirmed by 2+ sources | Default for well-sourced claims |
| `(single source)` | Only 1 source found | Append after citation |
| `⚠️ unverified` | No source found; from pre-trained knowledge | Append after the claim |
| `⚠️ conflicting` | Sources disagree | Present both sides with citations |

#### Handling Conflicts

When sources disagree, always surface both sides:

```markdown
Source A states X [source](url-a), while Source B states Y [source](url-b).
The discrepancy likely stems from [version difference / context difference /
one being outdated]. Based on [reasoning], Y appears more current.
```

Never silently pick one side.

### Phase 3: Self-Audit

Run this checklist before delivering:

- [ ] Every factual claim has a citation -- scan paragraph by paragraph
- [ ] No orphan claims -- claims without `[source](url)` or `⚠️ unverified`
- [ ] Source minimum met -- check against Phase 1 table
- [ ] No stale sources -- for `[current]` queries, primary sources are within 12 months
- [ ] Pre-trained knowledge declared -- any claim from memory is marked `⚠️ unverified`
- [ ] Code examples verified -- code came from a source, not generated from memory
- [ ] Version numbers checked -- no version numbers cited from memory alone
- [ ] Conflicts surfaced -- if sources disagreed, both sides are shown

**Verification mode additional checks:**
- [ ] At least 2 independent sources for the core claim
- [ ] Actively searched for counter-evidence (not just confirming sources)
- [ ] Verdict is explicit: Confirmed / Likely true / Disputed / Likely false / Cannot verify

**Fix before delivering:**
If audit finds issues -- search for a source or mark `⚠️ unverified`. Do not deliver with unchecked claims.

---

## Audit Mode (with `@` file paths)

For verifying technical correctness of existing artifacts (spec changes, design docs, code, etc.).

### Phase 1: Read & Extract Claims

1. Read all provided files using the Read tool
2. Extract every **verifiable technical claim**, including:
   - Pattern/approach choices: "Use X pattern to handle Y"
   - API/feature claims: "This API supports Z"
   - Recommendations: "Prefer library A over B because..."
   - Code correctness: API call signatures, parameter usage, deprecated methods
   - Version-specific claims: "Available since v3.0"
   - Configuration values: default settings, environment variables
3. Categorize each claim with scope tags (`[library]`, `[factual]`, `[api]`, etc.)

Output the extracted claims as a numbered list before proceeding:

```
Extracted claims from [filename]:
1. [library][factual] "useOptimistic hook accepts an update function as second parameter"
2. [api][current] "Next.js 15 Server Actions support FormData natively"
3. [architecture] "Repository pattern is recommended for data access layer"
...
```

### Phase 2: Verify Each Claim

For each extracted claim, run the Source Discovery workflow (same as Research mode Phase 1):
- context7 for `[library]` claims
- WebSearch + WebFetch for others
- Record the verification result for each claim

When artifacts contain many claims, prioritize by risk:
1. **High risk** (verify first): API calls, version numbers, security-related, deprecation claims
2. **Medium risk**: architectural patterns, best practice recommendations
3. **Lower risk**: general descriptions, well-established concepts

### Phase 3: Audit Report

Produce a structured report:

```markdown
## Audit Report: [filename(s)]

### Summary
- Total claims extracted: N
- ✅ Verified: X
- ⚠️ Unverified: Y
- ❌ Incorrect: Z
- 📝 Outdated: W

### Detailed Findings

#### ✅ Verified Claims
1. "useOptimistic hook accepts an update function" [source](url)
2. ...

#### ❌ Incorrect Claims
1. "Server Actions require 'use server' at the top of the file"
   - **Actual**: 'use server' can be placed at the top of the file OR before individual functions [source](url)
   - **Suggested fix**: Update the description to mention both placement options

#### ⚠️ Unverified Claims
1. "This approach reduces bundle size by 30%"
   - No authoritative source found to confirm this specific number

#### 📝 Outdated Claims
1. "Use getServerSideProps for server-side data fetching"
   - **Status**: Deprecated in Next.js 13+ App Router [source](url)
   - **Current approach**: Use Server Components or fetch in page/layout [source](url)

### Sources Consulted
| # | Source | Type | URL |
|---|--------|------|-----|
| 1 | ... | ... | ... |
```

### Phase 4: Self-Audit

Same checklist as Research mode Phase 3, plus:
- [ ] All extracted claims have a verification status (none left unchecked)
- [ ] Incorrect claims include suggested fixes with sources
- [ ] Outdated claims include the current recommended approach

---

## Anti-Patterns

These undermine the entire purpose of this skill. The reasoning explains why each matters -- understanding this helps avoid subtle violations too.

| Anti-pattern | Why it's harmful | Do this instead |
|-------------|-----------------|-----------------|
| Answer first, search later | Confirmation bias -- you'll only search for supporting evidence | Search first, answer from findings |
| Single search query | Missing perspectives, echo chamber | 2-3 queries with different phrasings |
| Cite without reading | URL might not say what you think it does | Always `WebFetch` before citing |
| Mix memory with sources | Reader can't distinguish verified from unverified | Mark ALL memory-based claims `⚠️ unverified` |
| Skip context7 for libraries | Missing the most authoritative, up-to-date source | Always check context7 first for `[library]` queries |
| Use agent-browser first | Slow, heavy, wastes tokens | Only after P0-P3 fail |
| Omit Sources Consulted table | Reader can't assess source quality | Always include the summary table |
| Silently drop conflicting info | Gives false certainty | Surface ALL conflicts |
| Generate code from memory | May contain outdated APIs or hallucinated patterns | Code must come from a cited source |

## Edge Cases

### Zero useful results
```
I searched for [query] using [N] search queries and consulted context7,
but found no authoritative sources addressing this specific question.

Based on pre-trained knowledge ⚠️ unverified: [best-effort answer]

Suggested next steps:
- Check [specific official channel]
- Try [alternative search terms]
- This may require direct experimentation
```

### Very new feature (< 1 month old)
Prioritize: official announcement blog posts, changelog entries, GitHub release notes.
Tag as `[current]` and search with exact version numbers.

### Mixed library-specific and general architecture
Split into sub-queries:
1. Library-specific part -> context7 + official docs
2. Architecture part -> WebSearch for design discussions + case studies

### Query is in Chinese
Translate to English for searching (most technical sources are English).
Deliver the answer in the user's language with English-language citations.

### Too many claims in Audit mode
Cap at the most critical 20 claims per run. Prioritize: API calls > version numbers > security > architecture > general descriptions. Mention: "Additional claims exist but were not verified in this pass."
