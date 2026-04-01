# pocketbase-IAM

AWS IAM-inspired policy-based access control for PocketBase.

## Features

- **Policy-based RBAC** — JSON policy documents with Allow/Deny statements
- **Deny overrides Allow** — explicit Deny always wins, matching AWS IAM evaluation
- **Multiple attachment paths** — policies attach to users directly, via roles, or via groups
- **Wildcard matching** — `*` patterns in actions and resources (`collections:*`, `*`)
- **Opt-in enforcement** — only registered "managed collections" are gated by IAM
- **Action registry** — discoverable actions via `iam_actions` view
- **Policy simulator** — test permissions before deploying (superuser-only)
- **Admin dashboard** — React/Refine UI for managing policies, roles, groups, and users

## Quick Start

**Prerequisites:** Go 1.24+, Node 18+ (for the admin UI)

```bash
git clone https://github.com/YakirOren/pocketbase-IAM.git
cd pocketbase-IAM
go run . serve
```

Open `http://localhost:8090/_/` and create a superuser account.

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
    },
    {
      "sid": "DenyDeleteAny",
      "effect": "Deny",
      "action": ["collections:delete"],
      "resource": ["*"]
    }
  ]
}
```

Each statement requires `effect` (Allow/Deny), `action` (array), and `resource` (array). Actions describe **what** can be done, resources describe **what it applies to** — following the AWS IAM model. `sid` is optional but recommended. `version` is required.

## How It Works

1. **Register managed collections** — add collection names to `iam_managed_collections` (superuser-only). This auto-sets the collection's PB rules to require authentication.
2. **Create policies** — define JSON policy documents in `iam_policies`.
3. **Attach policies** — link policies to users, roles, or groups via the join tables.
4. **Evaluation flow:**
   - Superusers and unauthenticated requests bypass IAM entirely
   - Any matching **Deny** statement → request denied (403)
   - Any matching **Allow** statement → request proceeds
   - No match → implicit deny (403)
5. **Non-managed collections** are unaffected — PocketBase's native rules apply as usual.

## Action & Resource Format

Actions describe the operation. Resources describe the target.

**CRUD operations** on managed collections:
- Action: `collections:<operation>` (e.g., `collections:read`, `collections:list`)
- Resource: the collection name (e.g., `posts`, `users`, or `*` for all)

Where operation is: `list`, `view`, `create`, `update`, `delete`.

One statement can grant access to multiple collections:
```json
{ "action": ["collections:read"], "resource": ["posts", "comments"] }
```

**Custom actions** can be registered in the action registry and checked via the API:
```json
{ "action": ["custom:billing:refund"], "resource": ["order:*"] }
```

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/iam/check` | Authenticated | Check if the current user can perform an action |
| `POST` | `/api/iam/simulate` | Superuser | Verbose evaluation with full trace |

**`/api/iam/check`** request:
```json
{ "action": "collections:read", "resource": "posts" }
```
Response:
```json
{ "allowed": true }
```

**`/api/iam/simulate`** request:
```json
{ "user_id": "USER_ID", "action": "collections:read", "resource": "posts" }
```
Response includes the full evaluation trace with matched statements.

## Collections

| Collection | Type | Purpose |
|------------|------|---------|
| `iam_policies` | Base | Policy documents |
| `iam_roles` | Base | Named roles |
| `iam_groups` | Base | User groups |
| `iam_managed_collections` | Base | Registry of IAM-enforced collections |
| `iam_action_registry` | Base | Registered custom actions |
| `iam_actions` | View | Unified view of all available actions |
| `iam_user_policies` | Join | User ↔ Policy |
| `iam_user_roles` | Join | User ↔ Role |
| `iam_role_policies` | Join | Role ↔ Policy |
| `iam_group_users` | Join | Group ↔ User |
| `iam_group_policies` | Join | Group ↔ Policy |

All IAM collections are system collections with superuser-only write access.

## Admin Dashboard

The `ui/` directory contains a React admin dashboard built with Refine, Tailwind CSS, and shadcn/ui. Pages include:

- **Policies** — create and manage policy documents
- **Roles** — define roles and attach policies
- **Groups** — manage user groups and their policies
- **Users** — view users and their IAM assignments
- **Managed Collections** — register/unregister collections for IAM enforcement
- **Simulator** — test policy evaluation for any user/action/resource combination

## Development

**Backend:**
```bash
go run . serve
```

**Frontend:**
```bash
cd ui
npm install
npm run dev
```

**Tests:**
```bash
go test ./iam/...
```

## Tech Stack

- **Backend:** Go, [PocketBase](https://pocketbase.io) v0.36.5, [ttlcache](https://github.com/jellydator/ttlcache)
- **Frontend:** React, [Refine](https://refine.dev), Tailwind CSS, [shadcn/ui](https://ui.shadcn.com)
