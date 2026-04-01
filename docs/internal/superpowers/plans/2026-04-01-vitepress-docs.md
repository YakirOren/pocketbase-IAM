# VitePress Documentation Site — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a VitePress documentation site to pocketbase-IAM, deployed to GitHub Pages, covering both developer integration and admin dashboard usage.

**Architecture:** VitePress docs live in `docs/` at the repo root. Existing `docs/plans/` moves to `docs/internal/plans/`. The site uses VitePress's default theme with custom CSS for branding (shield favicon, subtle hover effects). Content is organized into 5 sidebar sections: Get Started, Concepts, Dashboard, API Reference, and Examples.

**Tech Stack:** VitePress 1.6.4, Node.js, GitHub Pages

---

### Task 1: Scaffold VitePress Project

**Files:**
- Create: `docs/package.json`
- Create: `docs/.gitignore`
- Move: `docs/plans/` → `docs/internal/plans/`

- [ ] **Step 1: Move existing plans directory**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM
mkdir -p docs/internal
git mv docs/plans docs/internal/plans
```

- [ ] **Step 2: Create package.json**

Create `docs/package.json`:

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

- [ ] **Step 3: Create .gitignore**

Create `docs/.gitignore`:

```
node_modules
.vitepress/cache
.vitepress/dist
```

- [ ] **Step 4: Install dependencies**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM/docs
npm install
```

Expected: `node_modules` created, `package-lock.json` generated.

- [ ] **Step 5: Commit**

```bash
git add docs/package.json docs/package-lock.json docs/.gitignore docs/internal/
git commit -m "scaffold VitePress docs project"
```

---

### Task 2: VitePress Configuration and Theme

**Files:**
- Create: `docs/.vitepress/config.ts`
- Create: `docs/.vitepress/theme/index.ts`
- Create: `docs/.vitepress/theme/custom.css`
- Create: `docs/public/favicon.svg`

- [ ] **Step 1: Create VitePress config**

Create `docs/.vitepress/config.ts`:

```ts
import { defineConfig } from "vitepress";

export default defineConfig({
  title: "pocketbase-IAM",
  description: "AWS IAM-inspired access control for PocketBase",
  base: "/pocketbase-IAM/",
  cleanUrls: true,
  head: [
    [
      "link",
      {
        rel: "icon",
        type: "image/svg+xml",
        href: "/pocketbase-IAM/favicon.svg",
      },
    ],
  ],
  themeConfig: {
    logo: "/favicon.svg",
    nav: [
      { text: "Guide", link: "/getting-started" },
      { text: "API", link: "/api/setup-options" },
      { text: "Examples", link: "/examples/basic-crud" },
    ],
    sidebar: [
      {
        text: "Get Started",
        items: [
          { text: "What is pocketbase-IAM?", link: "/" },
          { text: "Getting Started", link: "/getting-started" },
        ],
      },
      {
        text: "Concepts",
        items: [
          { text: "Policies", link: "/concepts/policies" },
          { text: "Statements", link: "/concepts/statements" },
          {
            text: "Actions & Resources",
            link: "/concepts/actions-and-resources",
          },
          { text: "Evaluation Flow", link: "/concepts/evaluation-flow" },
          {
            text: "Managed Collections",
            link: "/concepts/managed-collections",
          },
          { text: "Roles", link: "/concepts/roles" },
          { text: "Groups", link: "/concepts/groups" },
          { text: "Caching", link: "/concepts/caching" },
        ],
      },
      {
        text: "Dashboard",
        items: [
          { text: "Overview", link: "/dashboard/overview" },
          { text: "Managing Policies", link: "/dashboard/managing-policies" },
          {
            text: "Managing Roles & Groups",
            link: "/dashboard/managing-roles-and-groups",
          },
          { text: "Policy Simulator", link: "/dashboard/policy-simulator" },
        ],
      },
      {
        text: "API Reference",
        collapsed: true,
        items: [
          { text: "Setup Options", link: "/api/setup-options" },
          { text: "Check Endpoint", link: "/api/check-endpoint" },
          { text: "Simulate Endpoint", link: "/api/simulate-endpoint" },
          { text: "Custom Actions", link: "/api/custom-actions" },
        ],
      },
      {
        text: "Examples",
        collapsed: true,
        items: [
          { text: "Basic CRUD Policy", link: "/examples/basic-crud" },
          { text: "Role-Based Access", link: "/examples/role-based-access" },
          { text: "Group Policies", link: "/examples/group-policies" },
          { text: "Wildcard Patterns", link: "/examples/wildcard-patterns" },
          { text: "Deny Overrides", link: "/examples/deny-overrides" },
        ],
      },
    ],
    socialLinks: [
      {
        icon: "github",
        link: "https://github.com/YakirOren/pocketbase-IAM",
      },
    ],
    search: {
      provider: "local",
    },
  },
});
```

- [ ] **Step 2: Create theme entry**

Create `docs/.vitepress/theme/index.ts`:

```ts
import DefaultTheme from "vitepress/theme";
import "./custom.css";

export default DefaultTheme;
```

- [ ] **Step 3: Create custom CSS**

Create `docs/.vitepress/theme/custom.css`:

```css
/* Nav logo: subtle scale on hover */
.VPNavBarTitle .logo {
  transition: transform 300ms ease-in-out;
}

.VPNavBarTitle:hover .logo {
  transform: scale(1.1);
}
```

- [ ] **Step 4: Create favicon SVG**

Create `docs/public/favicon.svg` — a shield icon in zinc-400:

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 100 100" fill="#a1a1aa">
  <path d="M50 8 L82 24 V52 C82 72 68 88 50 96 C32 88 18 72 18 52 V24 Z" />
  <path d="M50 18 L72 30 V52 C72 66 62 78 50 84 C38 78 28 66 28 52 V30 Z" fill="white" />
  <path d="M50 28 L62 34 V52 C62 60 58 68 50 72 C42 68 38 60 38 52 V34 Z" />
</svg>
```

- [ ] **Step 5: Verify dev server starts**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM/docs
npx vitepress dev --port 5174
```

Expected: VitePress dev server starts at `http://localhost:5174/pocketbase-IAM/`. You'll see a 404 until we add `index.md` in the next task. Kill the server after verifying it starts.

- [ ] **Step 6: Commit**

```bash
git add docs/.vitepress/ docs/public/
git commit -m "add VitePress config, theme, and favicon"
```

---

### Task 3: Get Started Pages

**Files:**
- Create: `docs/index.md`
- Create: `docs/getting-started.md`

- [ ] **Step 1: Create index.md (What is pocketbase-IAM?)**

Create `docs/index.md`:

```md
# What is pocketbase-IAM?

pocketbase-IAM is an AWS IAM-inspired access control library for [PocketBase](https://pocketbase.io). It adds policy-based RBAC to any PocketBase application with a single function call.

## How It Works

1. Register collections for IAM enforcement (opt-in)
2. Define JSON policy documents with Allow/Deny statements
3. Attach policies to users directly, via roles, or via groups
4. IAM evaluates every request — Deny always overrides Allow

Superusers and unauthenticated requests bypass IAM entirely. Non-managed collections use PocketBase's native rules as usual.

## Key Features

- **Policy-based RBAC** — JSON policy documents with Allow/Deny statements
- **Deny overrides Allow** — explicit Deny always wins, matching AWS IAM evaluation
- **Multiple attachment paths** — policies attach to users directly, via roles, or via groups
- **Wildcard matching** — `*` patterns in actions and resources
- **Opt-in enforcement** — only registered "managed collections" are gated
- **Action registry** — register and discover custom actions
- **Policy simulator** — test permissions before deploying (superuser-only)
- **Admin dashboard** — built-in React UI for managing everything

## Install

```bash
go get github.com/yakiroren/pocketbase-IAM/iam
```

## What's Next?

- [Getting Started](/getting-started) — set up IAM in your PocketBase app
- [Policies](/concepts/policies) — how policy documents work
- [Evaluation Flow](/concepts/evaluation-flow) — how requests are evaluated
- [Dashboard](/dashboard/overview) — manage permissions through the UI
```

- [ ] **Step 2: Create getting-started.md**

Create `docs/getting-started.md`:

```md
# Getting Started

## Prerequisites

- Go 1.24+
- Node 18+ (for the admin dashboard, optional)

## Install

```bash
go get github.com/yakiroren/pocketbase-IAM/iam
```

## Quick Start

```go
package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/yakiroren/pocketbase-IAM/iam"
)

func main() {
	app := pocketbase.New()

	if err := iam.Setup(app, iam.DefaultOptions()); err != nil {
		log.Fatalf("Failed to setup IAM: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
```

::: info What happens on first launch
When PocketBase starts, IAM automatically:

1. Creates 11 collections it needs (`iam_policies`, `iam_roles`, `iam_groups`, etc.)
2. Syncs rules on any already-managed collections
3. Registers enforcement hooks for all CRUD operations
:::

## Your First Policy

1. Open `http://localhost:8090/_/` and create a superuser account
2. Navigate to `/_/iam/` to open the IAM dashboard
3. Register a collection as managed (e.g., `posts`)
4. Create a policy:

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowReadPosts",
      "effect": "Allow",
      "action": ["collections:read", "collections:list"],
      "resource": ["posts"]
    }
  ]
}
```

5. Attach the policy to a user, role, or group

Now authenticated users with this policy can read and list posts, but all other operations are implicitly denied.

## What's Next?

- [Policies](/concepts/policies) — policy document format and validation
- [Statements](/concepts/statements) — Allow and Deny statement semantics
- [Actions & Resources](/concepts/actions-and-resources) — what actions and resources mean
- [Evaluation Flow](/concepts/evaluation-flow) — how the engine decides Allow or Deny
- [Managed Collections](/concepts/managed-collections) — how opt-in enforcement works
```

- [ ] **Step 3: Verify pages render**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM/docs
npx vitepress dev --port 5174
```

Expected: `http://localhost:5174/pocketbase-IAM/` shows "What is pocketbase-IAM?" page. Navigation and sidebar work. Kill server after verifying.

- [ ] **Step 4: Commit**

```bash
git add docs/index.md docs/getting-started.md
git commit -m "add Get Started pages"
```

---

### Task 4: Concepts Pages (Part 1 — Policies, Statements, Actions & Resources)

**Files:**
- Create: `docs/concepts/policies.md`
- Create: `docs/concepts/statements.md`
- Create: `docs/concepts/actions-and-resources.md`

- [ ] **Step 1: Create policies.md**

Create `docs/concepts/policies.md`:

```md
# Policies

A policy is a JSON document that defines permissions. Policies are stored in the `iam_policies` collection and attached to users, roles, or groups.

## Policy Format

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowReadPosts",
      "effect": "Allow",
      "action": ["collections:read", "collections:list"],
      "resource": ["posts"]
    }
  ]
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `version` | Yes | Policy version string. Must be `"2024-01-01"` |
| `statement` | Yes | Array of [statements](/concepts/statements) |

## Attaching Policies

Policies can be attached through three paths:

1. **Direct** — link a policy to a user via `iam_user_policies`
2. **Via role** — link a policy to a role via `iam_role_policies`, then assign the role to a user via `iam_user_roles`
3. **Via group** — link a policy to a group via `iam_group_policies`, then add a user to the group via `iam_group_users`

During [evaluation](/concepts/evaluation-flow), statements from all three paths are collected and evaluated together.

## Validation

Policies are validated when created or updated. The `version` field must be `"2024-01-01"`, and each statement must have a valid `effect`, at least one `action`, and at least one `resource`.

::: warning
An empty `statement` array is valid — it simply grants no permissions. The user will be implicitly denied everything.
:::
```

- [ ] **Step 2: Create statements.md**

Create `docs/concepts/statements.md`:

```md
# Statements

A statement is a single permission rule within a [policy](/concepts/policies). Each statement either allows or denies specific actions on specific resources.

## Statement Format

```json
{
  "sid": "AllowReadPosts",
  "effect": "Allow",
  "action": ["collections:read"],
  "resource": ["posts"]
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `sid` | No | Statement ID. A human-readable identifier for the statement |
| `effect` | Yes | `"Allow"` or `"Deny"` |
| `action` | Yes | Array of [action strings](/concepts/actions-and-resources) |
| `resource` | Yes | Array of [resource strings](/concepts/actions-and-resources) |

## Allow vs Deny

- **Allow** — grants permission to perform the listed actions on the listed resources
- **Deny** — explicitly blocks the actions, even if another statement allows them

::: warning
Deny always overrides Allow. If any statement from any policy (direct, role, or group) denies an action, it is denied — regardless of how many Allow statements exist.
:::

## Implicit Deny

If no statement matches a request (neither Allow nor Deny), the request is implicitly denied. You don't need Deny statements to block access — only to override existing Allows.

## Wildcard Matching

Both `action` and `resource` arrays support `*` wildcards:

```json
{ "action": ["collections:*"], "resource": ["*"] }
```

This matches any collection operation on any resource. See [Actions & Resources](/concepts/actions-and-resources) for details.
```

- [ ] **Step 3: Create actions-and-resources.md**

Create `docs/concepts/actions-and-resources.md`:

```md
# Actions & Resources

Actions describe **what** can be done. Resources describe **where** it applies.

## CRUD Actions

IAM automatically enforces these actions on [managed collections](/concepts/managed-collections):

| Action | PocketBase Operation |
|--------|---------------------|
| `collections:list` | List records |
| `collections:view` | View a single record |
| `collections:create` | Create a record |
| `collections:update` | Update a record |
| `collections:delete` | Delete a record |

::: info
`collections:read` is a convenience alias that matches both `collections:list` and `collections:view`.
:::

## Resources

For CRUD operations, the resource is the collection name:

```json
{ "action": ["collections:read"], "resource": ["posts"] }
```

A single statement can target multiple collections:

```json
{ "action": ["collections:read"], "resource": ["posts", "comments"] }
```

## Custom Actions

You can register custom actions for use outside of CRUD enforcement. Register them at startup:

```go
iam.RegisterAction(app, "custom:billing:refund", "Issue a billing refund")
```

Then check them via the [check endpoint](/api/check-endpoint):

```json
{ "action": "custom:billing:refund", "resource": "order:123" }
```

Custom actions appear in the `iam_actions` view alongside CRUD actions.

## Wildcards

Use `*` to match any value:

| Pattern | Matches |
|---------|---------|
| `collections:*` | Any collection operation |
| `*` | Any action or any resource |
| `custom:billing:*` | Any custom billing action |

Wildcards work in both `action` and `resource` arrays.
```

- [ ] **Step 4: Verify pages render**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM/docs
npx vitepress dev --port 5174
```

Expected: All three concepts pages render. Cross-links between them work. Kill server after verifying.

- [ ] **Step 5: Commit**

```bash
git add docs/concepts/policies.md docs/concepts/statements.md docs/concepts/actions-and-resources.md
git commit -m "add concepts: policies, statements, actions and resources"
```

---

### Task 5: Concepts Pages (Part 2 — Evaluation Flow, Managed Collections, Roles, Groups, Caching)

**Files:**
- Create: `docs/concepts/evaluation-flow.md`
- Create: `docs/concepts/managed-collections.md`
- Create: `docs/concepts/roles.md`
- Create: `docs/concepts/groups.md`
- Create: `docs/concepts/caching.md`

- [ ] **Step 1: Create evaluation-flow.md**

Create `docs/concepts/evaluation-flow.md`:

```md
# Evaluation Flow

When an authenticated user makes a CRUD request to a [managed collection](/concepts/managed-collections), IAM evaluates their permissions through this flow:

## Steps

1. **Skip check** — if the user is a superuser or the collection is not managed, the request passes through to PocketBase's native rules
2. **Collect statements** — IAM gathers all [statements](/concepts/statements) from the user's policies across all three attachment paths (direct, [role](/concepts/roles), [group](/concepts/groups))
3. **Check for Deny** — if any statement explicitly denies the action on the resource, the request is denied
4. **Check for Allow** — if any statement allows the action on the resource, the request proceeds
5. **Implicit deny** — if no statement matches, the request is denied

```
Request → Managed? → Superuser? → Collect Statements
                                        ↓
                                  Any Deny match? → Yes → DENIED
                                        ↓ No
                                  Any Allow match? → Yes → ALLOWED
                                        ↓ No
                                    DENIED (implicit)
```

## Deny Always Wins

This is the most important rule: an explicit Deny from **any** policy, attached through **any** path, overrides all Allow statements. This matches the AWS IAM evaluation model.

::: tip
Because Deny overrides Allow, you can use broad Allow statements (e.g., `"action": ["collections:*"], "resource": ["*"]`) and then add targeted Deny statements to restrict specific operations.
:::

## Denied Responses

Denied requests return **404 Not Found**, not 403 Forbidden. This prevents leaking whether a collection exists.

## Statement Collection

Statements are batch-fetched from all three paths in a single database round to avoid N+1 queries. Results are [cached](/concepts/caching) per user.
```

- [ ] **Step 2: Create managed-collections.md**

Create `docs/concepts/managed-collections.md`:

```md
# Managed Collections

IAM enforcement is opt-in. Only collections registered in `iam_managed_collections` are gated by IAM policies. All other collections use PocketBase's native rules.

## Registering a Collection

Add a collection name to `iam_managed_collections` through the [dashboard](/dashboard/overview) or directly via the PocketBase admin UI. This is a superuser-only operation.

When a collection is registered as managed, IAM automatically sets its PocketBase rules to `@request.auth.id != ''`. This ensures:

- **Unauthenticated requests** are blocked at the PocketBase layer (before IAM)
- **Authenticated requests** pass through to IAM for policy evaluation
- **Superusers** bypass both PocketBase rules and IAM

::: warning
The rule is `@request.auth.id != ''` (not empty string `""`). Setting it to `""` would make the collection fully public, bypassing IAM entirely.
:::

## Unregistering a Collection

When a collection is removed from `iam_managed_collections`, IAM resets its PocketBase rules to `nil` (superuser-only access). You'll need to manually set appropriate rules if you want non-superuser access.

## Non-Managed Collections

Collections not in `iam_managed_collections` are completely unaffected by IAM. Their PocketBase rules work as usual.
```

- [ ] **Step 3: Create roles.md**

Create `docs/concepts/roles.md`:

```md
# Roles

A role is a named collection of policies. Instead of attaching policies directly to each user, you define a role and assign it to users.

## How Roles Work

1. Create a role in `iam_roles` (e.g., "editor", "viewer")
2. Attach policies to the role via `iam_role_policies`
3. Assign the role to users via `iam_user_roles`

During [evaluation](/concepts/evaluation-flow), all policies from all of a user's roles are collected alongside their direct policies and [group](/concepts/groups) policies.

## When to Use Roles

Roles are useful when multiple users need the same set of permissions. Instead of attaching the same policies to each user individually, create a role and assign it once.

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "EditorAccess",
      "effect": "Allow",
      "action": ["collections:read", "collections:create", "collections:update"],
      "resource": ["posts", "comments"]
    }
  ]
}
```

Attach this policy to an "editor" role, then assign the role to any user who should be able to read, create, and update posts and comments.

## Roles vs Direct Policies

There is no difference in evaluation. Statements from roles are treated identically to direct statements. Roles are an organizational tool — they make it easier to manage permissions at scale.
```

- [ ] **Step 4: Create groups.md**

Create `docs/concepts/groups.md`:

```md
# Groups

A group is a collection of users that share policies. Groups provide a third path for attaching policies alongside [direct attachment](/concepts/policies#attaching-policies) and [roles](/concepts/roles).

## How Groups Work

1. Create a group in `iam_groups` (e.g., "engineering", "support")
2. Add users to the group via `iam_group_users`
3. Attach policies to the group via `iam_group_policies`

During [evaluation](/concepts/evaluation-flow), all policies from all of a user's groups are collected alongside direct and role policies.

## Groups vs Roles

| | Roles | Groups |
|-|-------|--------|
| **Attaches to** | Users | Users |
| **Contains** | Policies | Users + Policies |
| **Use case** | "What can this permission set do?" | "Who is in this team?" |

Both are organizational tools. Use roles when you're thinking about permissions ("editors can do X"). Use groups when you're thinking about people ("the engineering team gets Y").

## Combining Groups and Roles

A user can belong to multiple groups and have multiple roles simultaneously. All statements from all paths are collected and evaluated together. [Deny always overrides Allow](/concepts/evaluation-flow#deny-always-wins).
```

- [ ] **Step 5: Create caching.md**

Create `docs/concepts/caching.md`:

```md
# Caching

IAM uses a dual LRU+TTL cache to avoid redundant database queries during policy evaluation.

## How It Works

When IAM evaluates a request, it checks the cache for the user's collected statements. On a cache miss, it queries the database and stores the result.

The cache has two eviction strategies:

- **LRU** — least recently used entries are evicted when the cache reaches `CacheMaxSize`
- **TTL** — entries expire after `CacheTTL` regardless of usage

## Configuration

Configure cache behavior via `iam.Options`:

```go
iam.Setup(app, iam.Options{
    CacheMaxSize: 10_000,   // max entries (default: 10,000)
    CacheTTL:     60 * time.Second, // entry lifetime (default: 60s)
})
```

## Invalidation

The cache is invalidated automatically when IAM data changes. Changes to any of these collections trigger invalidation for affected users:

- `iam_policies` — invalidates all users with this policy (direct, via role, or via group)
- `iam_user_policies` — invalidates the affected user
- `iam_user_roles` — invalidates the affected user
- `iam_role_policies` — invalidates all users with this role
- `iam_group_users` — invalidates the affected user
- `iam_group_policies` — invalidates all users in this group

::: info
Cache invalidation is hook-driven — changes take effect immediately, not after TTL expiry.
:::
```

- [ ] **Step 6: Verify all concepts pages render**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM/docs
npx vitepress dev --port 5174
```

Expected: All 8 concepts pages render. Cross-links work. Kill server after verifying.

- [ ] **Step 7: Commit**

```bash
git add docs/concepts/
git commit -m "add concepts: evaluation flow, managed collections, roles, groups, caching"
```

---

### Task 6: Dashboard Pages

**Files:**
- Create: `docs/dashboard/overview.md`
- Create: `docs/dashboard/managing-policies.md`
- Create: `docs/dashboard/managing-roles-and-groups.md`
- Create: `docs/dashboard/policy-simulator.md`

- [ ] **Step 1: Create overview.md**

Create `docs/dashboard/overview.md`:

```md
# Dashboard Overview

pocketbase-IAM includes a built-in admin dashboard at `/_/iam/`. It provides a web UI for managing all IAM resources.

## Accessing the Dashboard

Start your PocketBase server and navigate to `http://localhost:8090/_/iam/`. You must be logged in as a superuser.

## Pages

| Page | Purpose |
|------|---------|
| **Policies** | Create and edit JSON policy documents |
| **Roles** | Define roles and attach policies to them |
| **Groups** | Manage user groups and their policies |
| **Users** | View users and their IAM assignments (direct policies, roles, groups) |
| **Managed Collections** | Register or unregister collections for IAM enforcement |
| **Simulator** | Test policy evaluation for any user/action/resource combination |

## Tech Stack

The dashboard is built with React, [Refine](https://refine.dev), Tailwind CSS, and [shadcn/ui](https://ui.shadcn.com). It is embedded in the Go binary — no separate frontend deployment is needed.

## Development

To develop the dashboard with hot reload:

```bash
cd ui
npm install
npm run dev
```

The Vite dev server proxies API requests to PocketBase at `:8090`.

To build the dashboard into the Go binary:

```bash
cd ui
npm run build
```

This outputs to `iam/dashboard/`, which is embedded via `//go:embed`.
```

- [ ] **Step 2: Create managing-policies.md**

Create `docs/dashboard/managing-policies.md`:

```md
# Managing Policies

The Policies page lets you create, edit, and delete [policy documents](/concepts/policies).

## Creating a Policy

1. Click **Create** on the Policies page
2. Enter a name for the policy
3. Write the policy JSON in the editor:

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowReadPosts",
      "effect": "Allow",
      "action": ["collections:read"],
      "resource": ["posts"]
    }
  ]
}
```

4. Save the policy

The editor validates the JSON on save. Invalid policies (wrong version, missing fields, bad effect values) are rejected with an error message.

## Attaching a Policy

After creating a policy, attach it to users, roles, or groups:

- **To a user** — go to the Users page, select a user, and add the policy
- **To a role** — go to the Roles page, select a role, and add the policy
- **To a group** — go to the Groups page, select a group, and add the policy

See [Policies — Attaching Policies](/concepts/policies#attaching-policies) for how the three attachment paths work.
```

- [ ] **Step 3: Create managing-roles-and-groups.md**

Create `docs/dashboard/managing-roles-and-groups.md`:

```md
# Managing Roles & Groups

The Roles and Groups pages let you organize permissions and users.

## Roles

### Creating a Role

1. Click **Create** on the Roles page
2. Enter a role name (e.g., "editor", "viewer", "admin")
3. Save the role

### Attaching Policies to a Role

1. Select a role
2. Add one or more [policies](/concepts/policies)
3. Save

### Assigning a Role to Users

1. Go to the Users page
2. Select a user
3. Add the role
4. Save

## Groups

### Creating a Group

1. Click **Create** on the Groups page
2. Enter a group name (e.g., "engineering", "support")
3. Save the group

### Adding Users to a Group

1. Select a group
2. Add users
3. Save

### Attaching Policies to a Group

1. Select a group
2. Add one or more [policies](/concepts/policies)
3. Save

## Duplicate Prevention

IAM prevents duplicate assignments. You cannot assign the same role to a user twice, add the same user to a group twice, or attach the same policy to a role/group twice. The dashboard shows an error if you try.
```

- [ ] **Step 4: Create policy-simulator.md**

Create `docs/dashboard/policy-simulator.md`:

```md
# Policy Simulator

The Simulator page lets you test policy evaluation without making real requests. It is available to superusers only.

## Using the Simulator

1. Select a **user** to simulate
2. Enter an **action** (e.g., `collections:read`)
3. Enter a **resource** (e.g., `posts`)
4. Click **Simulate**

The simulator runs the full [evaluation flow](/concepts/evaluation-flow) and shows:

- **Result** — whether the action would be allowed or denied
- **Matched statements** — which statements matched the action and resource, from which policies, and through which attachment path (direct, role, or group)

## API Equivalent

The simulator uses the [simulate endpoint](/api/simulate-endpoint) under the hood:

```bash
curl -X POST http://localhost:8090/api/iam/simulate \
  -H "Content-Type: application/json" \
  -d '{"user_id": "USER_ID", "action": "collections:read", "resource": "posts"}'
```

::: tip
Use the simulator to debug unexpected Deny results. The matched statements trace shows exactly which statement caused the denial and where it came from.
:::
```

- [ ] **Step 5: Verify dashboard pages render**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM/docs
npx vitepress dev --port 5174
```

Expected: All 4 dashboard pages render with working cross-links. Kill server after verifying.

- [ ] **Step 6: Commit**

```bash
git add docs/dashboard/
git commit -m "add dashboard pages"
```

---

### Task 7: API Reference Pages

**Files:**
- Create: `docs/api/setup-options.md`
- Create: `docs/api/check-endpoint.md`
- Create: `docs/api/simulate-endpoint.md`
- Create: `docs/api/custom-actions.md`

- [ ] **Step 1: Create setup-options.md**

Create `docs/api/setup-options.md`:

```md
# Setup Options

`iam.Setup()` initializes the IAM system on a PocketBase app. It registers routes, hooks, enforcement, migrations, and the admin dashboard.

## Signature

```go
func Setup(app core.App, opts Options) error
```

## Options

```go
type Options struct {
    CacheMaxSize int           // Max entries in the LRU cache (default: 10,000)
    CacheTTL     time.Duration // Cache entry lifetime (default: 60s)
    Logger       *slog.Logger  // Structured logger (default: app's logger)
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `CacheMaxSize` | `int` | `10000` | Maximum number of entries in the policy LRU cache |
| `CacheTTL` | `time.Duration` | `60s` | How long cached evaluations remain valid |
| `Logger` | `*slog.Logger` | `nil` | Custom logger. If nil, uses the PocketBase app's logger |

## Default Options

```go
opts := iam.DefaultOptions()
// CacheMaxSize: 10_000
// CacheTTL: 60s
// Logger: nil (uses app's logger)
```

## Example

```go
iam.Setup(app, iam.Options{
    CacheMaxSize: 50_000,
    CacheTTL:     5 * time.Minute,
})
```

## What Setup Does

1. Validates options
2. Creates the policy cache
3. Registers API routes (`/api/iam/check`, `/api/iam/simulate`)
4. Registers enforcement hooks on all CRUD operations
5. Registers policy validation, duplicate prevention, and cache invalidation hooks
6. Registers managed collection sync hooks
7. Serves the embedded dashboard at `/_/iam/`
8. Syncs managed collection rules on boot

Migrations are auto-registered via `init()` — no manual import needed.
```

- [ ] **Step 2: Create check-endpoint.md**

Create `docs/api/check-endpoint.md`:

```md
# Check Endpoint

`POST /api/iam/check` evaluates whether the authenticated user can perform an action on a resource.

## Request

**Auth required:** Yes (any authenticated user)

```bash
curl -X POST http://localhost:8090/api/iam/check \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action": "collections:read", "resource": "posts"}'
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | `string` | Yes | The action to check (e.g., `collections:read`) |
| `resource` | `string` | Yes | The resource to check (e.g., `posts`) |

## Response

```json
{ "allowed": true }
```

| Field | Type | Description |
|-------|------|-------------|
| `allowed` | `boolean` | Whether the action is permitted |

## Behavior

- Superusers always get `{"allowed": true}`
- Unauthenticated requests get 401
- The endpoint runs the full [evaluation flow](/concepts/evaluation-flow)

::: tip
Use this endpoint for [custom actions](/api/custom-actions) that aren't automatically enforced by CRUD hooks. For example, checking if a user can issue a refund before processing it in your application code.
:::
```

- [ ] **Step 3: Create simulate-endpoint.md**

Create `docs/api/simulate-endpoint.md`:

```md
# Simulate Endpoint

`POST /api/iam/simulate` runs a verbose policy evaluation for any user. Returns the full evaluation trace including matched statements.

## Request

**Auth required:** Superuser only

```bash
curl -X POST http://localhost:8090/api/iam/simulate \
  -H "Authorization: Bearer SUPERUSER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "USER_ID", "action": "collections:read", "resource": "posts"}'
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | `string` | Yes | The user ID to simulate |
| `action` | `string` | Yes | The action to check |
| `resource` | `string` | Yes | The resource to check |

## Response

The response includes the evaluation result and a trace of all matched statements:

```json
{
  "allowed": false,
  "reason": "explicit_deny",
  "matched_statements": [
    {
      "sid": "DenyDeleteAny",
      "effect": "Deny",
      "action": ["collections:delete"],
      "resource": ["*"],
      "policy_id": "abc123",
      "policy_name": "RestrictDeletes",
      "source": "role",
      "source_name": "editor"
    }
  ]
}
```

## Use Cases

- Debug why a user is being denied access
- Verify a policy change has the intended effect before deploying
- Audit which statements are granting or denying access

The [dashboard simulator](/dashboard/policy-simulator) provides a UI for this endpoint.
```

- [ ] **Step 4: Create custom-actions.md**

Create `docs/api/custom-actions.md`:

```md
# Custom Actions

IAM automatically enforces `collections:*` actions on managed collections. For application-specific permissions, you can register custom actions.

## Registering an Action

```go
iam.RegisterAction(app, "custom:billing:refund", "Issue a billing refund")
```

Call `RegisterAction` after `iam.Setup()` and before `app.Start()`. Registration is idempotent — calling it multiple times with the same action is safe.

## Using Custom Actions in Policies

Once registered, custom actions can be used in [policy statements](/concepts/statements):

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowRefunds",
      "effect": "Allow",
      "action": ["custom:billing:refund"],
      "resource": ["order:*"]
    }
  ]
}
```

## Checking Custom Actions

Custom actions are not automatically enforced. Use the [check endpoint](/api/check-endpoint) in your application code:

```go
app.OnServe().BindFunc(func(se *core.ServeEvent) error {
    se.Router.POST("/api/refund/{orderId}", func(e *core.RequestEvent) error {
        // Check IAM permission before proceeding
        // ... call POST /api/iam/check with action and resource
    })
    return se.Next()
})
```

## Action Registry

All registered custom actions appear in the `iam_actions` view alongside the built-in CRUD actions. This view is used by the [dashboard](/dashboard/overview) for autocomplete and discoverability.
```

- [ ] **Step 5: Verify API pages render**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM/docs
npx vitepress dev --port 5174
```

Expected: All 4 API reference pages render. Collapsed sidebar section expands on click. Kill server after verifying.

- [ ] **Step 6: Commit**

```bash
git add docs/api/
git commit -m "add API reference pages"
```

---

### Task 8: Examples Pages

**Files:**
- Create: `docs/examples/basic-crud.md`
- Create: `docs/examples/role-based-access.md`
- Create: `docs/examples/group-policies.md`
- Create: `docs/examples/wildcard-patterns.md`
- Create: `docs/examples/deny-overrides.md`

- [ ] **Step 1: Create basic-crud.md**

Create `docs/examples/basic-crud.md`:

```md
# Basic CRUD Policy

Allow a user to read and create posts.

## Policy

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "ReadAndCreatePosts",
      "effect": "Allow",
      "action": ["collections:read", "collections:create"],
      "resource": ["posts"]
    }
  ]
}
```

## What This Grants

| Action | Resource | Result |
|--------|----------|--------|
| `collections:list` | `posts` | Allowed (`read` matches `list` and `view`) |
| `collections:view` | `posts` | Allowed |
| `collections:create` | `posts` | Allowed |
| `collections:update` | `posts` | Denied (implicit) |
| `collections:delete` | `posts` | Denied (implicit) |
| `collections:read` | `comments` | Denied (resource doesn't match) |

## Setup

1. Register `posts` as a [managed collection](/concepts/managed-collections)
2. Create the policy above in `iam_policies`
3. Attach the policy to a user via `iam_user_policies`
```

- [ ] **Step 2: Create role-based-access.md**

Create `docs/examples/role-based-access.md`:

```md
# Role-Based Access

Define an "editor" role that can read, create, and update posts and comments.

## Policy

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "EditorAccess",
      "effect": "Allow",
      "action": ["collections:read", "collections:create", "collections:update"],
      "resource": ["posts", "comments"]
    }
  ]
}
```

## Setup

1. Register `posts` and `comments` as [managed collections](/concepts/managed-collections)
2. Create the policy above
3. Create a role called "editor" in `iam_roles`
4. Attach the policy to the "editor" role via `iam_role_policies`
5. Assign the "editor" role to users via `iam_user_roles`

Any user with the "editor" role can now read, create, and update posts and comments. Deleting is implicitly denied.
```

- [ ] **Step 3: Create group-policies.md**

Create `docs/examples/group-policies.md`:

```md
# Group Policies

Give the "support" team read access to all managed collections.

## Policy

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "SupportReadAll",
      "effect": "Allow",
      "action": ["collections:read"],
      "resource": ["*"]
    }
  ]
}
```

## Setup

1. Create a group called "support" in `iam_groups`
2. Create the policy above
3. Attach the policy to the "support" group via `iam_group_policies`
4. Add team members to the group via `iam_group_users`

All users in the "support" group can now list and view records in every managed collection. They cannot create, update, or delete.
```

- [ ] **Step 4: Create wildcard-patterns.md**

Create `docs/examples/wildcard-patterns.md`:

```md
# Wildcard Patterns

Use `*` to match any action or resource.

## Full Admin Access

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "FullAccess",
      "effect": "Allow",
      "action": ["*"],
      "resource": ["*"]
    }
  ]
}
```

This allows all operations on all managed collections.

## All Operations on One Collection

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllPostsOperations",
      "effect": "Allow",
      "action": ["collections:*"],
      "resource": ["posts"]
    }
  ]
}
```

## Read All Collections

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "ReadEverything",
      "effect": "Allow",
      "action": ["collections:read"],
      "resource": ["*"]
    }
  ]
}
```

## Custom Action Wildcards

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllBillingActions",
      "effect": "Allow",
      "action": ["custom:billing:*"],
      "resource": ["*"]
    }
  ]
}
```
```

- [ ] **Step 5: Create deny-overrides.md**

Create `docs/examples/deny-overrides.md`:

```md
# Deny Overrides

Deny statements always override Allow, regardless of which policy or attachment path they come from.

## Broad Allow + Targeted Deny

Allow full access but block deletes:

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowEverything",
      "effect": "Allow",
      "action": ["*"],
      "resource": ["*"]
    },
    {
      "sid": "BlockDeletes",
      "effect": "Deny",
      "action": ["collections:delete"],
      "resource": ["*"]
    }
  ]
}
```

Even though `AllowEverything` matches `collections:delete`, the explicit Deny wins.

## Cross-Policy Deny

Deny works across policies and attachment paths. If a user has:

- **Direct policy** → Allow `collections:*` on `posts`
- **Role policy** → Deny `collections:delete` on `*`

The user can read, create, and update posts, but cannot delete anything. The Deny from the role policy overrides the Allow from the direct policy.

::: warning
Be careful with broad Deny statements. A Deny on `*`/`*` from any policy attached through any path will block everything, and no Allow statement can override it.
:::
```

- [ ] **Step 6: Verify examples pages render**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM/docs
npx vitepress dev --port 5174
```

Expected: All 5 examples pages render. Collapsed sidebar section expands on click. Kill server after verifying.

- [ ] **Step 7: Commit**

```bash
git add docs/examples/
git commit -m "add example pages"
```

---

### Task 9: GitHub Actions Deployment

**Files:**
- Create: `.github/workflows/docs.yml`

- [ ] **Step 1: Create GitHub Actions workflow**

Create `.github/workflows/docs.yml`:

```yaml
name: Deploy Docs

on:
  push:
    branches: [main]
    paths:
      - "docs/**"
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm
          cache-dependency-path: docs/package-lock.json

      - name: Install dependencies
        run: npm ci
        working-directory: docs

      - name: Build
        run: npm run build
        working-directory: docs

      - uses: actions/upload-pages-artifact@v3
        with:
          path: docs/.vitepress/dist

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/docs.yml
git commit -m "add GitHub Actions workflow for docs deployment"
```

---

### Task 10: Final Verification

- [ ] **Step 1: Full build check**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM/docs
npm run build
```

Expected: Build completes without errors. Output in `.vitepress/dist/`.

- [ ] **Step 2: Preview built site**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM/docs
npm run preview -- --port 5174
```

Expected: All pages render at `http://localhost:5174/pocketbase-IAM/`. Navigation, sidebar, search, and cross-links all work. Kill server after verifying.

- [ ] **Step 3: Verify search works**

Open the site, use the search bar (Ctrl+K / Cmd+K). Search for "deny" — should find results in evaluation-flow, statements, and deny-overrides pages.

- [ ] **Step 4: Commit any fixes**

If any issues were found, fix and commit them.
