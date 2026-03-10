---
name: social-reader
description: "Use when reading, fetching, or extracting content from social media platforms (X/Twitter, Threads). Triggers on: read tweet, read X post, fetch Twitter article, read Threads post, 讀取推文, 讀取 X 文章, 抓取推特內容, 讀取 Threads 貼文, 抓取 Threads 內容, 讀留言, 看留言, read comments, any x.com/status, x.com/article, threads.net/post, or threads.com URL."
---

# Social Reader

Fetch and parse social media posts using the `social-reader` CLI tool. Supports X (Twitter) and Threads.

When a user pastes a social media URL, they usually want to see the original content in full — and decide for themselves whether they need a translation or summary. Summarizing without being asked loses details the user may care about. So the default behavior is to present the fetched content faithfully.

## Output Modes

Determine the mode from the user's message:

| User says | CLI command | Post-processing |
|---|---|---|
| Just a URL (no extra instructions) | `social-reader <url>` | Present output as-is |
| 「讀留言」/「看留言」/ "read comments" (Threads only) | `social-reader <url> --comments` | Present output as-is |
| 「幫我翻譯」/ "translate" | `social-reader <url>` | Translate the CLI output |
| 「幫我總結」/ "summarize" / 「摘要」 | `social-reader <url>` | Summarize the CLI output |
| 「讀留言 + 翻譯」(Threads only) | `social-reader <url> --comments` | Translate the CLI output |
| 「讀留言 + 摘要」(Threads only) | `social-reader <url> --comments` | Summarize the CLI output |

> **Note:** X/Twitter does not support `--comments`. Replies require login and cannot be fetched via Jina Reader. The CLI will show the post content and print a warning to stderr.

When in doubt, default to presenting as-is without translation or summary.

## Workflow

### Step 1: Run the CLI

```bash
social-reader <url>                    # Post only
social-reader <url> --comments         # Post + comment tree
social-reader <url> --json             # JSON output (for programmatic use)
```

The CLI handles:
- URL normalization (twitter.com→x.com, threads.com→threads.net, etc.)
- Fetching via Jina Reader
- Parsing post content, images, engagement metrics
- Comment extraction with nested tree rendering
- Profile picture / noise stripping

### Step 2: Present the Output

Use the CLI's markdown output directly. The output format includes:
- Author and content
- Images (as markdown links)
- Engagement metrics (❤️ likes · 💬 replies · 🔁 reposts)
- Source URL
- Comment tree with box-drawing characters (when --comments)

Then apply post-processing based on the output mode:
- **Raw**: Output as-is
- **Translate**: Translate the full content, preserving structure
- **Summarize**: Replace the full content with translated key points

## Supported Platforms

| Platform | URL patterns |
|---|---|
| X (Twitter) | `x.com/*/status/*`, `twitter.com/*/status/*`, `x.com/*/article/*` |
| Threads | `threads.net/@*/post/*`, `threads.com/@*/post/*`, `threads.net/t/*` |

## Edge Cases

- **Rate limiting**: The CLI retries once on 429. If it still fails, inform the user.
- **404 / deleted content**: Tell the user the post may have been removed.
- **Protected / private accounts**: Jina Reader cannot access these — inform the user.
- **X/Twitter + comments**: X does not support `--comments` — replies require login. The CLI will warn on stderr and still show the post.
- **Articles + comments**: X articles also don't support comment mode.
