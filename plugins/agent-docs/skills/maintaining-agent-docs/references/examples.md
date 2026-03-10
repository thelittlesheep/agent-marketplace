# Examples

## Example: E-commerce Site Documentation

```
.agent/
├── README.md
├── System/
│   ├── architecture.md
│   │   - Tech: Next.js frontend, Node.js backend, PostgreSQL, AWS
│   │   - Components: Auth service, Product API, Order service, Payment gateway
│   │   - Integration: Stripe (payments), SendGrid (emails)
│   ├── database-schema.md
│   │   - users, products, orders, payments tables
│   │   - Key relationships and indices
│   └── deployment.md
│       - Production on AWS Lambda
│       - Staging on same Lambda
│       - Rollback procedure
├── Tasks/
│   ├── user-authentication.md
│   │   - PRD: Support email/password login and OAuth
│   │   - What we learned: JWT vs session tokens (chose JWT)
│   │   - Gotchas: Handle token refresh on mobile
│   ├── product-search.md
│   ├── payment-processing.md
│   └── README.md (links to all Tasks)
└── SOP/
    ├── add-new-product-field.md
    │   - Steps: Update schema → Update API → Update UI form
    │   - Checklist: Test in staging, update docs
    ├── deploy-to-production.md
    ├── database-migrations.md
    └── README.md (links to all SOPs)
```

## Example: Update After Implementing "Product Search"

**What changed:**
- Added Elasticsearch to system
- New API endpoint: /api/products/search
- New database index on product fields

**Doc updates:**

1. Tasks/product-search.md → Add final learnings
2. System/architecture.md → Add Elasticsearch to "Tech Stack" and "Integration Points"
3. System/database-schema.md → Add new index
4. SOP/add-new-product-field.md → Add step about maintaining search index
5. .agent/README.md → No new links, but verify all links still work

## Example: CLAUDE.md — Index-Based vs Bloated

**Correct (index-based, ~0.5KB):**

```markdown
# E-commerce App

Next.js e-commerce platform with Stripe payments.

## Commands

- **Dev**: `npm run dev`
- **Test**: `npm test`
- **Build**: `npm run build`

## Rules

- All amounts stored as integers (cents)
- Use kebab-case for filenames

## Docs Index

Prefer retrieval-led reasoning over pre-training-led reasoning.
Complete documentation in `.agent/` directory:

.agent/|root: .agent/README.md
|System:{architecture.md,database-schema.md,deployment.md}
|Tasks:{user-authentication.md,product-search.md,payment-processing.md}
|SOP:{add-new-product-field.md,deploy-to-production.md,database-migrations.md}
```

**Wrong (bloated, 4KB+ and growing):**

```markdown
# E-commerce App

Next.js e-commerce platform with Stripe payments.

## Architecture                    ← belongs in .agent/System/architecture.md
- Frontend: Next.js 14 with App Router, React Server Components...
- Backend: Node.js with Express, REST API endpoints at /api/v2/...
- Database: PostgreSQL 15 with pgvector extension for search...
- Infrastructure: AWS Lambda, API Gateway, CloudFront, S3...
- Auth: JWT with refresh tokens, stored in httpOnly cookies...

## Code Style                      ← belongs in user scope or .agent/
- Use functional components with hooks
- Prefer named exports over default exports
- Always use TypeScript strict mode
- ...30 more rules...

## API Endpoints                   ← belongs in .agent/System/
- GET /api/products - List products with pagination...
- POST /api/orders - Create new order...
- ...20 more endpoints...

## Database Schema                 ← belongs in .agent/System/database-schema.md
- users: id, email, password_hash, created_at...
- products: id, name, price_cents, category_id...
- ...10 more tables with full column definitions...
```

**Why the bloated version is wrong:** Every line above 8KB wastes tokens on every conversation. The agent loads it all but only needs 10% per session. The index version lets agents retrieve exactly what they need.

---

## Real-World Impact

Projects with structured .agent documentation see:
- **Onboarding time:** 50% reduction (days → hours)
- **"Where's the docs?" questions:** 90% reduction
- **Duplicate documentation:** Eliminated (single source of truth)
- **Architecture misalignment:** Caught in doc review (before code review)
- **Knowledge preservation:** New team members have full context without tribal knowledge
