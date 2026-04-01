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
