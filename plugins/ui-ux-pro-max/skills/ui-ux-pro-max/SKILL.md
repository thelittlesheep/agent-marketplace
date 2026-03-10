---
name: ui-ux-pro-max
description: "Use when building, designing, or reviewing any UI/UX — websites, dashboards, landing pages, mobile apps, SaaS, e-commerce, or any .html/.tsx/.vue/.svelte component. Covers style selection, color palettes, typography, chart types, layout patterns, and accessibility. Supports React, Next.js, Vue, Svelte, SwiftUI, React Native, Flutter, Tailwind. Always use this skill when the user mentions UI design, frontend styling, component aesthetics, dark mode, responsive layout, or visual polish."
---

# UI/UX Pro Max - Design Intelligence

Searchable database of UI styles, color palettes, font pairings, chart types, product recommendations, UX guidelines, and stack-specific best practices.

## Search

```bash
uiux-cli search "<keyword>" --domain <domain> [-n <max_results>] [--json]
uiux-cli search "<keyword>" --stack <stack> [-n <max_results>] [--json]
uiux-cli domains   # list available domains
uiux-cli stacks    # list available stacks
```

**Domains:** style, prompt, color, chart, landing, product, ux, typography
**Stacks:** html-tailwind (default), react, nextjs, vue, svelte, swiftui, react-native, flutter

Omit `--domain` and `--stack` to auto-detect domain from query keywords.

## Workflow

When user requests UI/UX work, search multiple domains to build a complete design:

1. **Product** — style recommendations for product type
2. **Style** — detailed style guide (colors, effects, frameworks)
3. **Typography** — font pairings with Google Fonts imports
4. **Color** — color palette (Primary, Secondary, CTA, Background, Text, Border)
5. **Landing** — page structure (if landing page)
6. **Chart** — chart recommendations (if dashboard/analytics)
7. **UX** — best practices and anti-patterns
8. **Stack** — stack-specific guidelines (default: html-tailwind)

---

## Common Rules for Professional UI

### Icons & Visual Elements

| Rule | Do | Don't |
|------|----|----- |
| **No emoji icons** | Use SVG icons (Heroicons, Lucide, Simple Icons) | Use emojis like 🎨 🚀 ⚙️ as UI icons |
| **Stable hover states** | Use color/opacity transitions on hover | Use scale transforms that shift layout |
| **Correct brand logos** | Research official SVG from Simple Icons | Guess or use incorrect logo paths |
| **Consistent icon sizing** | Use fixed viewBox (24x24) with w-6 h-6 | Mix different icon sizes randomly |

### Interaction & Cursor

| Rule | Do | Don't |
|------|----|----- |
| **Cursor pointer** | Add `cursor-pointer` to all clickable/hoverable cards | Leave default cursor on interactive elements |
| **Hover feedback** | Provide visual feedback (color, shadow, border) | No indication element is interactive |
| **Smooth transitions** | Use `transition-colors duration-200` | Instant state changes or too slow (>500ms) |

### Light/Dark Mode Contrast

| Rule | Do | Don't |
|------|----|----- |
| **Glass card light mode** | Use `bg-white/80` or higher opacity | Use `bg-white/10` (too transparent) |
| **Text contrast light** | Use `#0F172A` (slate-900) for text | Use `#94A3B8` (slate-400) for body text |
| **Muted text light** | Use `#475569` (slate-600) minimum | Use gray-400 or lighter |
| **Border visibility** | Use `border-gray-200` in light mode | Use `border-white/10` (invisible) |

### Layout & Spacing

| Rule | Do | Don't |
|------|----|----- |
| **Floating navbar** | Add `top-4 left-4 right-4` spacing | Stick navbar to `top-0 left-0 right-0` |
| **Content padding** | Account for fixed navbar height | Let content hide behind fixed elements |
| **Consistent max-width** | Use same `max-w-6xl` or `max-w-7xl` | Mix different container widths |

---

## Pre-Delivery Checklist

### Visual Quality
- [ ] No emojis used as icons (use SVG instead)
- [ ] All icons from consistent icon set (Heroicons/Lucide)
- [ ] Brand logos are correct (verified from Simple Icons)
- [ ] Hover states don't cause layout shift
- [ ] Use theme colors directly (bg-primary) not var() wrapper

### Interaction
- [ ] All clickable elements have `cursor-pointer`
- [ ] Hover states provide clear visual feedback
- [ ] Transitions are smooth (150-300ms)
- [ ] Focus states visible for keyboard navigation

### Light/Dark Mode
- [ ] Light mode text has sufficient contrast (4.5:1 minimum)
- [ ] Glass/transparent elements visible in light mode
- [ ] Borders visible in both modes
- [ ] Test both modes before delivery

### Layout
- [ ] Floating elements have proper spacing from edges
- [ ] No content hidden behind fixed navbars
- [ ] Responsive at 320px, 768px, 1024px, 1440px
- [ ] No horizontal scroll on mobile

### Accessibility
- [ ] All images have alt text
- [ ] Form inputs have labels
- [ ] Color is not the only indicator
- [ ] `prefers-reduced-motion` respected
