---
name: code-explainer
description: "Use when you want to understand code issues in depth - explains WHY code is problematic, shows bad vs good examples, and teaches best practices as if you're a beginner. Triggers on: review requests, learning requests, code-quality/code-review-router findings, 解釋程式碼, 為什麼這樣寫不好, 程式碼教學."
---

# Code Explainer

分析程式碼問題並教導「為什麼」這樣寫不好。流程：範圍判定 → 偵測 → 用好壞對照教學 → 總結。

**這是一個唯讀技能。絕不修改檔案 — 只分析和教學。**

## 語言規則

**所有輸出必須使用繁體中文（台灣用語）。** 包含：
- 報告標題、說明、教學內容全部用繁體中文
- 程式碼本身、變數名稱、技術識別符保持原文
- 程式碼區塊內的註解使用英文（維持程式碼慣例）
- 技術術語首次出現時以繁體中文解釋，可附上英文原文方便查閱

## Scope

| Mode | Trigger | Command |
|------|---------|---------|
| Default | No path specified | `git diff --name-only` |

## Language Detection

Detect language from file extensions and adapt teaching examples accordingly:

| Extension | Language | Beginner Perspective |
|-----------|----------|---------------------|
| `.py` | Python | Python beginner |
| `.ts`, `.tsx` | TypeScript | TypeScript beginner |
| `.js`, `.jsx` | JavaScript | JavaScript beginner |
| `.sh`, `.bash` | Shell | Shell scripting beginner |
| Other | Infer from extension | General beginner |

## Analysis Categories

Scan using `rg` and `ast-grep` for the following categories. Use language-appropriate patterns.

### 1. Anti-patterns

Language-specific anti-patterns:

```bash
# Python: mutable default arguments
ast-grep -p 'def $NAME($$$, $ARG=$MUT_DEFAULT, $$$): $$$' -l py --json
rg 'def \w+\(.*=\[\]|=\{\}|=set\(\)' -t py

# TypeScript/JavaScript: any abuse, == instead of ===
rg '\bany\b' -t ts
ast-grep -p '$A == $B' -l ts --json
ast-grep -p '$A != $B' -l ts --json

# Shell: unquoted variables
rg '\$\w+[^"]' -t sh
```

### 2. Security

```bash
# SQL injection — string concatenation in queries
rg '(f"|f'\''|\.format\(|% ).*(?i)(SELECT|INSERT|UPDATE|DELETE)' -t py
rg '`.*\$\{.*\}.*(?i)(SELECT|INSERT|UPDATE|DELETE)`' -t ts -t js

# Hardcoded secrets
rg '(?i)(password|secret|api_key|token)\s*=\s*["\x27][^"\x27]+["\x27]' --type-not json --type-not md

# Shell injection
rg 'os\.system\(|subprocess\.call\(.*shell=True' -t py
rg 'child_process\.exec\(' -t ts -t js
```

### 3. Error Handling

```bash
# Python: bare except, pass in except
ast-grep -p 'try: $$$ except: $$$' -l py --json
rg 'except.*:\s*$' -t py
rg 'except.*:\s*pass' -t py

# TypeScript/JavaScript: empty catch
ast-grep -p 'catch ($E) {}' -l ts --json
ast-grep -p 'catch ($E) { }' -l ts --json

# Swallowed errors — catch without logging or re-throwing
rg 'catch\s*\(' -t ts -t js -A 3
```

### 4. Performance

```bash
# Python: list comprehension inside loop (potential N+1)
rg 'for .+ in .+:' -t py -A 5

# N+1 query patterns — query inside loop
rg '(\.query\(|\.find\(|\.get\(|\.fetch\()' -t py -t ts -t js

# Unnecessary re-renders (React)
rg 'useEffect\(\s*\(\)\s*=>' -t tsx -t jsx -A 5
rg 'useState.*new (Object|Array|Map|Set)\(' -t tsx -t jsx
```

### 5. Readability

```bash
# Magic numbers
rg '[^0-9a-zA-Z_](?:[2-9]\d{2,}|[1-9]\d{3,})[^0-9a-zA-Z_.]' -t py -t ts -t js

# Single-letter variables (excluding common loop vars i,j,k,x,y)
ast-grep -p 'const $N = $$$' -l ts --json
ast-grep -p 'let $N = $$$' -l ts --json

# Deeply nested code (4+ levels)
rg '^\s{16,}(if|for|while|switch)' -t py -t ts -t js
```

### 6. Type Safety

```bash
# TypeScript: any abuse
rg ': any\b|as any\b|<any>' -t ts

# Missing null checks
ast-grep -p '$OBJ.$PROP.$NESTED' -l ts --json

# Non-null assertion abuse
rg '!\.' -t ts
rg '!\[' -t ts
```

## 輸出格式

使用以下結構產出報告，**所有說明文字必須使用繁體中文**：

```markdown
## 程式碼解說報告

**範圍**: [git diff / 指定路徑 / 全專案]
**語言**: [Python / TypeScript / ...]
**發現問題數**: [N]

---

### 問題 #1: [問題標題]

**檔案**: `path/to/file.py:42`
**分類**: 反模式 / 安全性 / 錯誤處理 / 效能 / 可讀性 / 型別安全
**嚴重度**: 🔴 高 / 🟡 中 / 🟢 低

#### 哪裡有問題？
[1-2 句清楚描述 — 假設讀者是完全沒見過這種寫法的初學者]

#### 為什麼這樣寫不好？
[用真實場景解釋後果，不要說「這不符合最佳實踐」這種空話]
- [具體後果 1 — 例如「攻擊者可以竊取你資料庫中的所有資料」]
- [具體後果 2 — 例如「你的應用程式會無聲地遺失使用者資料」]

#### ❌ 有問題的寫法（你的程式碼）
```<lang>
# the actual problematic code from the file, with context
```

#### ✅ 正確做法
```<lang>
# the root-cause fix — NOT a band-aid / NOT a suppression comment
```

#### 改了什麼？為什麼？
[逐行解釋每一個改動，用初學者能理解的方式寫]
1. [改動 1]: [為什麼這個改動很重要]
2. [改動 2]: [為什麼這個改動很重要]

#### ⚠️ 不要用這些偷懶寫法
[列出常見的治標不治本做法，並解釋為什麼有害]
- `# type: ignore` / `# noqa` / `// @ts-ignore` — [為什麼不該用]
- [其他常見偷懶寫法] — [為什麼不該用]

#### 延伸學習
- **[術語/概念]**（英文原文）: [一句話的初學者友善定義]
```

依嚴重度分組排列：🔴 高 → 🟡 中 → 🟢 低。

## 教學原則

以下規則適用於報告中的每一個問題，**不可違反**：

1. **只修根因** — 永遠修正實際問題。絕不建議使用 `# type: ignore`、`# noqa`、`# eslint-disable`、`// @ts-ignore`、`# nosec` 等壓制警告的做法。這些只是把問題藏起來，不是修好它。

2. **初學者用語** — 不假設讀者知道任何技術術語。第一次使用術語時（例如「可變的 (mutable)」、「注入 (injection)」、「副作用 (side effect)」），必須用白話文解釋，括號內附英文原文。

3. **真實後果** — 解釋在正式環境中實際會發生什麼。不好的寫法：「這不符合最佳實踐。」好的寫法：「攻擊者可以透過注入 SQL 指令來讀取、修改或刪除你資料庫中的所有資料。」

4. **一次只教一個概念** — 每個問題只教一個概念。如果一段程式碼有多個問題，拆成獨立的問題分別說明。

5. **不接受偷懶修法** — 每個問題都必須包含「不要用這些偷懶寫法」區塊，列出常見的治標做法並解釋為什麼有害。

6. **展示實際程式碼** — 「有問題的寫法」必須是使用者的實際程式碼（不是泛用範例）。「正確做法」必須是可以直接替換的根因修正。

7. **解釋每一個改動** — 「改了什麼？為什麼？」區塊必須解釋好壞對照中每一行不同的地方，用從來沒寫過程式的人也能理解的方式說明。

8. **全程繁體中文** — 所有說明、解釋、教學內容必須使用繁體中文（台灣用語）。程式碼和技術識別符保持原文。

## 嚴重度分類

| 嚴重度 | 判定標準 |
|--------|----------|
| 🔴 高 | 安全漏洞、資料遺失風險、會導致程式崩潰的錯誤 |
| 🟡 中 | 效能問題、錯誤處理缺陷、會造成隱性 bug 的反模式 |
| 🟢 低 | 可讀性問題、風格問題、輕微的型別安全改善 |

## 誤報避免

- **魔術數字 (Magic numbers)**: 跳過 0、1、-1、HTTP 狀態碼（200、404、500）、常見 port 號
- **單字母變數**: 跳過迴圈中的 `i`、`j`、`k`；座標中的 `x`、`y`；catch 區塊中的 `e`
- **`any` 型別**: 跳過用在泛用函式庫包裝器或測試工具中的情況
- **深度巢狀**: 跳過模板/JSX 巢狀 — 只標記邏輯巢狀
- **可變預設值 (Mutable defaults)**: 跳過使用 `field(default_factory=...)` 或類似模式包裝的情況
- **字串串接**: 跳過 logging 語句和除錯輸出
- **不確定時**: 仍然列出發現，但註明「這可能是刻意的寫法 — 請與你的團隊確認」

## 下一步

- 要修正問題？使用 `code-quality`（死代碼清理 + SOLID 違規重構）
