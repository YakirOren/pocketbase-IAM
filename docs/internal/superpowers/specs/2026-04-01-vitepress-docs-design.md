# VitePress Documentation Site — Design Spec

## Overview

Add a VitePress documentation site to pocketbase-IAM, serving two audiences: Go developers integrating the library and admins managing permissions via the dashboard. Deployed to GitHub Pages at `/pocketbase-IAM/`.

## Project Setup

**VitePress root:** `docs/`

Existing `docs/plans/` moves to `docs/internal/plans/` to coexist with VitePress content. The `internal/` directory is excluded from the built site.

**package.json:**

```json
{
  "name": "pocketbase-iam-docs",
  "type": "module",
  "scripts": {
    "dev": "vitepress dev",
    "build": "vitepress build",
    "preview": "vitepress preview"
  },
  "devDependencies": {
    "vitepress": "^1.6.4"
  }
}
```

**`.gitignore` (docs-level):**

```
node_modules
.vitepress/cache
.vitepress/dist
```

## VitePress Configuration

File: `docs/.vitepress/config.ts`

- **Title:** "pocketbase-IAM"
- **Description:** "AWS IAM-inspired access control for PocketBase"
- **Base:** `/pocketbase-IAM/`
- **Clean URLs:** enabled
- **Search:** local provider
- **Top nav:** Guide | API | Examples | GitHub (social link to repo)
- **Favicon:** `/pocketbase-IAM/favicon.svg`

## Theme & Branding

Extend the default VitePress theme with custom CSS, matching the Turbine docs style.

File: `docs/.vitepress/theme/index.ts` — imports DefaultTheme + `custom.css`

File: `docs/.vitepress/theme/custom.css`:
- Logo hover effect: subtle scale or glow (no spin)
- Consistent with the security/IAM theme

File: `docs/public/favicon.svg` — shield or lock icon in a neutral color (e.g., zinc-400 `#a1a1aa`)

## Directory Structure

```
docs/
├── .vitepress/
│   ├── config.ts
│   └── theme/
│       ├── index.ts
│       └── custom.css
├── public/
│   └── favicon.svg
├── index.md
├── getting-started.md
├── concepts/
│   ├── policies.md
│   ├── statements.md
│   ├── actions-and-resources.md
│   ├── evaluation-flow.md
│   ├── managed-collections.md
│   ├── roles.md
│   ├── groups.md
│   └── caching.md
├── dashboard/
│   ├── overview.md
│   ├── managing-policies.md
│   ├── managing-roles-and-groups.md
│   └── policy-simulator.md
├── api/
│   ├── setup-options.md
│   ├── check-endpoint.md
│   ├── simulate-endpoint.md
│   └── custom-actions.md
├── examples/
│   ├── basic-crud.md
│   ├── role-based-access.md
│   ├── group-policies.md
│   ├── wildcard-patterns.md
│   └── deny-overrides.md
├── internal/
│   └── plans/          (moved from docs/plans/)
├── .gitignore
└── package.json
```

## Sidebar Configuration

### Get Started (expanded)
- What is pocketbase-IAM? → `/`
- Getting Started → `/getting-started`

### Concepts (expanded)
- Policies → `/concepts/policies`
- Statements → `/concepts/statements`
- Actions & Resources → `/concepts/actions-and-resources`
- Evaluation Flow → `/concepts/evaluation-flow`
- Managed Collections → `/concepts/managed-collections`
- Roles → `/concepts/roles`
- Groups → `/concepts/groups`
- Caching → `/concepts/caching`

### Dashboard (expanded)
- Overview → `/dashboard/overview`
- Managing Policies → `/dashboard/managing-policies`
- Managing Roles & Groups → `/dashboard/managing-roles-and-groups`
- Policy Simulator → `/dashboard/policy-simulator`

### API Reference (collapsed)
- Setup Options → `/api/setup-options`
- Check Endpoint → `/api/check-endpoint`
- Simulate Endpoint → `/api/simulate-endpoint`
- Custom Actions → `/api/custom-actions`

### Examples (collapsed)
- Basic CRUD Policy → `/examples/basic-crud`
- Role-Based Access → `/examples/role-based-access`
- Group Policies → `/examples/group-policies`
- Wildcard Patterns → `/examples/wildcard-patterns`
- Deny Overrides → `/examples/deny-overrides`

## Landing Page (index.md)

No hero section. A simple overview page with:

- H1: "pocketbase-IAM"
- Brief description of what the library is and what it does
- High-level how-it-works summary (evaluation flow in a few sentences)
- Key features as a bullet list (not a features grid)
- Link to Getting Started

## Content Sources

| Section | Source |
|---------|--------|
| index.md | README.md intro + features list |
| getting-started.md | README.md quick start section |
| Concepts (policies, statements, actions) | README.md policy format + CLAUDE.md architecture |
| Concepts (evaluation-flow) | CLAUDE.md evaluation flow section |
| Concepts (managed-collections) | CLAUDE.md key design decisions |
| Concepts (roles, groups) | README.md collections table + new content |
| Concepts (caching) | CLAUDE.md cache description + new content |
| Dashboard pages | New content describing UI workflows |
| API pages | README.md API endpoints + CLAUDE.md setup options |
| Examples | New content with policy JSON snippets and Go code |

## Content Style

- **No frontmatter** — page titles from H1 headings
- **Code-first** — show Go/JSON examples before explanations
- **Callouts** — `:::info`, `:::warning`, `:::tip` for notes and gotchas
- **Tables** for options, fields, and collection schemas
- **Cross-linking** between related pages
- **Concise technical tone** — Go docs style, no fluff
- **Code blocks** — Go for backend, JSON for policies, bash for CLI

## Deployment

GitHub Pages via GitHub Actions. The workflow builds VitePress and deploys to the `gh-pages` branch. Base path `/pocketbase-IAM/` matches the repository name.
