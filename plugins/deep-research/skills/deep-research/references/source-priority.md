# Source Priority Hierarchy

## Tier 1: Highest Authority (P0-P1)

### context7 Library Documentation
- **When**: Query involves a specific library or framework
- **Tool**: `mcp__context7__resolve-library-id` -> `mcp__context7__query-docs`
- **Why**: Curated, version-aware documentation with code examples
- **Caveat**: May not have bleeding-edge features; cross-verify for `[current]` queries

### Official Documentation Sites
- **Pattern**: `*.dev`, `docs.*`, `developer.*`
- **Examples**:
  - react.dev, vuejs.org, svelte.dev, angular.dev
  - docs.python.org, doc.rust-lang.org, go.dev/doc
  - developer.mozilla.org (MDN)
  - nodejs.org/docs, deno.land/manual, bun.sh/docs
  - nextjs.org/docs, nuxt.com/docs, astro.build/docs
  - tailwindcss.com/docs
  - docs.github.com, docs.aws.amazon.com, cloud.google.com/docs
  - supabase.com/docs, firebase.google.com/docs
  - typescriptlang.org/docs

### Official Blogs & Changelogs
- **Examples**: react.dev/blog, nextjs.org/blog, blog.rust-lang.org
- **When**: `[current]` tagged queries, release announcements
- **Caveat**: Blog posts can be aspirational; verify features actually shipped

## Tier 2: High Authority (P2)

### GitHub Repositories
- **Source code**: definitive for "does X support Y?" questions
- **Issues/Discussions**: real-world problem reports and solutions
- **Release notes**: version-specific changes
- **README.md**: project-maintained documentation

### Recognized Technical Publications
- **Examples**: web.dev (Google), engineering blogs from major companies
- **RFCs and specs**: tc39.es, w3c.org, whatwg.org, ietf.org

### Curated Community Knowledge
- **Stack Overflow**: answers with 10+ upvotes, accepted answers
- **GitHub Discussions**: maintainer-endorsed answers

## Tier 3: Medium Authority (P3)

### Technical Blogs & Tutorials
- **Good signals**: author has GitHub presence, content is recent, includes working code
- **Examples**: kentcdodds.com, tkdodo.eu, joshwcomeau.com
- **Caveat**: individual opinions, may not represent consensus

### Conference Talks & Podcasts (transcripts)
- **When**: architecture decisions, design philosophy
- **Caveat**: opinions can evolve; check date

## Tier 4: Low Authority -- Use with Caution

### AI-generated content aggregators
- **Red flags**: no author attribution, suspiciously comprehensive, generic examples
- **Rule**: never use as sole source; cross-verify everything

### Forum posts without verification
- **Examples**: Reddit without upvotes, Discord messages, random blog comments
- **Rule**: treat as leads for further research, not as sources

### Outdated documentation (2+ years for fast-moving tech)
- **Rule**: if query is `[current]`, deprioritize; if about fundamentals, still usable

## Source Evaluation Quick-Check

Before citing any source, verify:

1. **Authorship**: Is the author identifiable? Are they associated with the project?
2. **Date**: When was this published/updated? Is it relevant to the query's timeframe?
3. **Specificity**: Does it directly address the question, or is it tangential?
4. **Corroboration**: Does at least one other source agree?
5. **Code examples**: Does the code actually work, or is it illustrative pseudo-code?

## Domain Allowlist (high-trust)

```
react.dev, vuejs.org, svelte.dev, angular.dev, nextjs.org,
docs.python.org, doc.rust-lang.org, go.dev, nodejs.org, bun.sh,
developer.mozilla.org, web.dev, typescriptlang.org,
github.com, stackoverflow.com,
docs.aws.amazon.com, cloud.google.com, learn.microsoft.com,
supabase.com, firebase.google.com, vercel.com/docs
```

## Domain Blocklist (skip these)

```
w3schools.com (frequent errors),
geeksforgeeks.org (often outdated/inaccurate),
tutorialspoint.com (outdated),
medium.com articles without verified author (quality varies wildly)
```
