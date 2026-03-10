---
name: code-quality
description: "Use when code needs cleanup or refactoring - covers dead code removal (unused imports, dead functions, unreferenced variables, duplicate code, deprecated patterns) AND SOLID violation detection (god classes, rigid branching, broken substitution, fat interfaces, tight coupling). Triggers on: cleanup, refactor, code quality, pre-merge housekeeping, linter warnings, SOLID violations, 清理程式碼, 重構, 死代碼, 程式碼品質."
---

# Code Quality

Analyze code for dead code and structural issues. Workflow: analyze → report → confirm → execute → verify.

## Scope

| Mode | Trigger | Command |
|------|---------|---------|
| Default | No path specified | `git diff --name-only` (unstaged + staged) |
| Commit | User provides commit | Auto-detect: contains `..` → `git diff <range> --name-only`, otherwise → `git show <commit> --name-only` |

## Phase 1: Dead Code Analysis

### 1. Unused Imports

```bash
ast-grep -p 'import { $$$NAMES } from "$SOURCE"' -l ts --json
ast-grep -p 'from $MODULE import $$$NAMES' -l py --json
```

Cross-reference: `rg -w '<name>' <file>` — 1 match (import only) = unused.

### 2. Dead Functions/Classes

```bash
ast-grep -p 'function $NAME($$$) { $$$ }' -l ts --json
ast-grep -p 'def $NAME($$$): $$$' -l py --json
```

Verify: `rg -l '\b<name>\b' -t ts -t js -t py` — only self-referencing = dead. **Exclude** exports, tests, entry points.

### 3. Unused Variables

```bash
npx eslint --no-eslintrc --rule '{"no-unused-vars":"warn"}' <file>  # TS/JS
ruff check --select F841,F811 <file>                                 # Python
```

### 4. Duplicate Code

```bash
# Example pattern — generate project-specific patterns dynamically
ast-grep -p 'if ($A === null || $A === undefined)' -l ts
```

This is an example pattern. Derive actual patterns from the codebase during analysis.
Flag only **3+ occurrences** or **10+ lines** duplicates.

### 5. Deprecated Patterns

```bash
rg '@deprecated|warnings\.warn\(|DeprecationWarning'
rg 'TODO.*remov|FIXME.*dead|HACK.*temporary' -i
```

## Phase 2: SOLID Violation Analysis

### SRP — God Class/Function

```bash
ast-grep -p 'class $NAME { $$$BODY }' -l ts --json
ast-grep -p 'def $NAME($$$): $$$' -l py --json
```

Flag: **8+ methods** or **50+ lines**. Suggest: Extract Class/Method.

### OCP — Rigid Branching

```bash
ast-grep -p 'switch ($EXPR) { $$$ }' -l ts --json
ast-grep -p 'if ($A) { $$$ } else if ($B) { $$$ }' -l ts --json
```

Match 2+ branches, then count actual branches in results. Flag: **4+ branches**. Suggest: Strategy / Factory.

### LSP — Broken Substitution

```bash
ast-grep -p 'throw new Error("Not implemented")' -l ts
ast-grep -p 'raise NotImplementedError($$$)' -l py
```

Suggest: Composition over Inheritance.

### ISP — Fat Interface

```bash
ast-grep -p 'interface $NAME { $$$BODY }' -l ts --json
ast-grep -p 'class $NAME(Protocol): $$$BODY' -l py --json
```

Flag: **8+ members**. Suggest: Split into focused interfaces.

### DIP — Tight Coupling

```bash
# Step 1: Find constructors
ast-grep -p 'constructor($$$) { $$$ }' -l ts --json
# Step 2: In matched files, search for direct instantiation
rg 'new \w+\(' <matched-files>

# Python: Find direct instantiation in __init__
ast-grep -p 'self.$ATTR = $CLASS($$$)' -l py --json
# Then verify the match is inside __init__ method
```

Flag: Direct instantiation in constructors. Suggest: Dependency Injection.

## Report

```
| # | Phase | Category | Severity | File:Line | Description | Action |
|---|-------|----------|----------|-----------|-------------|--------|
| 1 | Dead Code | Unused import | Low | src/utils.ts:3 | lodash | Remove |
| 2 | SOLID | SRP | High | src/api.ts:12 | 15 methods | Extract Class |
```

For High severity SOLID violations: include **before/after** preview. Group: High → Medium → Low.

## Confirmation

Present per-category summary. **NEVER modify files without user confirmation.**

```
Dead Code: 3 unused imports, 2 dead functions
SOLID: 2 High (SRP), 3 Medium (OCP)
Which to fix? (all / dead-code-only / solid-only / by number)
```

## Verification

1. Run linter + tests — refactoring must not change behavior
2. Show `git diff --stat`
3. If failures → revert and report
4. **Before claiming done:** Use `superpowers:verification-before-completion`

## False Positive Checklist

- **Dynamic imports**: `require(var)`, `importlib.import_module()` — skip
- **DI containers**: Verify with `rg` before flagging
- **Public API exports**: Never remove library entry point exports
- **Test utilities**: Not dead code
- **Event handlers**: Check dynamic registration
- **Type-only imports**: `import type { Foo }` — may appear unused but needed by TS
- **Side-effect imports**: `import './polyfill'` — no bindings but has side effects
- **Decorators / Annotations**: Referenced via reflection, not direct calls
- **Barrel exports**: `export * from` — re-exports may mask usage
- **Module augmentation**: TS declaration merging — appears unused but extends types
- **When uncertain**: Flag as "review needed", don't delete

## Safety Checklist

- **Over-engineering**: Skip Strategy for 2-branch if/else
- **Premature abstraction**: No interfaces for single implementations
- **Untested code**: Warn user, suggest tests first
- **Breaking public API**: Flag signature changes explicitly
- **Cascading changes**: If 10+ files affected, confirm scope
- **Generated code**: Skip protobuf, OpenAPI, or auto-generated files
- **Legacy migration in progress**: Don't refactor code that's actively being migrated
- **有意義的抽象**: 不要因為「只用一次」就移除 helper — 若它提升可讀性則保留
- **Readability over brevity**: 清理或重構後的程式碼不應比之前更難讀
- **不合併無關邏輯**: 清理重複時，不要將不同職責的相似程式碼強行合併
- **不要用過於複雜的設計模式替代簡單的 if/else** — 簡單的 2-branch 不需要 Strategy
- **Preserve useful abstractions**: 不要在重構過程中移除有助於理解的中間層

## 與其他技能的關係

- **code-simplifier**: 處理可讀性改善（不涉及結構性缺陷）。如果問題是「能動但難讀」，使用 code-simplifier。
- **code-explainer**: 唯讀分析，教學用途。不修改檔案。
