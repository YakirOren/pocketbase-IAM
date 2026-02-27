# PocketBase IAM - AWS-style RBAC System

## Context

Build an AWS IAM-inspired RBAC system as a custom PocketBase Go application. The system supports **users, user groups, roles, and policies** — mirroring [AWS IAM Identities](https://docs.aws.amazon.com/IAM/latest/UserGuide/id.html). Policies can be attached to users, groups, and roles. Permission evaluation follows AWS's model: Explicit Deny > Explicit Allow > Implicit Deny.

AssumeRole (temporary role sessions) is deferred to v2.

This is a greenfield project. Target: PocketBase v0.36.5 (latest stable).

---

## Key Design Decisions

### Use PocketBase's built-in collection API — no custom CRUD routes

All IAM entities are PocketBase collections. CRUD is handled by PB's standard `/api/collections/{name}/records` endpoints. We only add **1 custom route**:

- `POST /api/iam/check` — check if an action is allowed for the current user

### IAM-managed collections via registry (opt-in)

IAM enforcement only applies to collections explicitly registered in the `iam_managed_collections` collection. Non-registered collections are untouched — PB's native rules handle them as before. This means:
- No surprise lockouts when adding IAM to an existing app
- Internal/system collections stay unaffected
- Admin explicitly controls which collections IAM protects

When a collection is registered as IAM-managed, its PB rules are auto-set to `@request.auth.id != ''` (authenticated-only) so IAM becomes the sole action-level gatekeeper for authenticated users, while unauthenticated requests are blocked at the PB rule layer. Admins can still add PB rule expressions for **row-level filtering** (e.g., `@request.auth.id = user`).

### Unauthenticated requests bypass IAM

Unauthenticated requests skip IAM entirely and fall through to PB's native collection rules. This preserves public access patterns (public blog posts, etc.). IAM only evaluates **authenticated** users.

### Multiple roles per user

A user can be assigned many roles via `iam_user_roles`. The engine collects policies from **all sources** (direct + groups + roles), unions them, then applies deny-overrides-allow.

### User groups (like AWS IAM groups)

Groups contain users. Policies attach to groups. Users inherit all group policies. Multiple group membership allowed. No nesting.

### Custom resources are freeform strings

No registry needed. App code and policies agree on naming conventions:

```go
allowed, _ := iam.Evaluate(app, userID, "custom:billing:refund", "order:12345")
```
```json
{ "action": ["custom:billing:refund"], "resource": ["order:*"] }
```

For PocketBase CRUD hooks, resource is always `*` (IAM operates at action-level for collections). Specific resource matching only applies to custom actions via `/api/iam/check` or Go code.

### Wildcard matching rules

| Pattern | Matches | Does NOT match |
|---------|---------|----------------|
| `*` | **Everything** (any string, any depth) | — |
| `collections:*:read` | `collections:posts:read` | `collections:posts:create` |
| `collections:posts:*` | `collections:posts:read`, `collections:posts:delete` | `collections:users:read` |
| `order:*` | `order:123`, `order:abc` | `invoice:123` |

Rules:
- A lone `*` matches **any string** regardless of segments
- `*` within a colon-delimited pattern matches **exactly one segment**
- Patterns and values are split on `:` and compared segment by segment
- Segment counts must match (unless pattern is lone `*`)

### In-memory policy cache (LRU + TTL)

Resolved policies cached per user using `github.com/jellydator/ttlcache/v3`:

- **Key:** user ID → **Value:** resolved `[]Statement`
- **TTL:** 60 seconds per entry
- **Max size:** configurable (default 10,000 entries). LRU eviction when full — inactive users evicted first, active users stay hot.
- **Smart invalidation** via hooks on join table changes:
  - `iam_user_policies` / `iam_user_roles` / `iam_group_users` change → invalidate that specific user
  - `iam_group_policies` change → invalidate all users in that group (1 query to find them)
  - `iam_role_policies` change → invalidate all users with that role (1 query to find them)
  - `iam_policies` document update → find affected users via join tables, invalidate only them
- **Thundering herd mitigation:** After bulk invalidation, only users who actually make requests trigger DB re-population. LRU + natural request arrival distributes the load.

### Duplicate prevention on join tables

A hook on create for each join table checks if a record with the same combination already exists (e.g., same user+role in `iam_user_roles`). Rejects with 400 if duplicate.

### Deny context is logged, not returned

`/api/iam/check` returns `{"allowed": true/false}` — no reason exposed to client. The denial reason (implicit deny vs. explicit deny + which policy) is logged server-side via PocketBase's logger, matching how PB handles HTTP error details.

---

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│                    PocketBase App                        │
│                                                          │
│  ┌─────────────┐  ┌──────────┐  ┌─────────────────┐     │
│  │ 1 Custom    │  │ Policy   │  │ Enforcement     │     │
│  │ Route       │  │ Engine   │  │ Hooks           │     │
│  │ /iam/check  │  │ + Cache  │  │ (only on        │     │
│  │             │  │          │  │  managed colls)  │     │
│  └─────────────┘  └──────────┘  └─────────────────┘     │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │      PocketBase Built-in Collection CRUD API       │  │
│  │   (standard /api/collections/*/records routes)     │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │              IAM Collections (DB)                  │  │
│  │  iam_managed_collections  iam_policies             │  │
│  │  iam_roles    iam_groups                           │  │
│  │  iam_role_policies   iam_user_policies             │  │
│  │  iam_user_roles      iam_group_users               │  │
│  │  iam_group_policies                                │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

---

## 1. Project Structure

```
pocketbase-IAM/
├── main.go                 # Entry point: boot PocketBase, register IAM
├── go.mod / go.sum
├── iam/
│   ├── policy.go           # Policy document types + parsing + validation
│   ├── helpers.go          # Wildcard matching, action string builder
│   ├── cache.go            # LRU+TTL cache wrapping ttlcache/v3
│   ├── engine.go           # Permission evaluation (collect + evaluate)
│   ├── routes.go           # POST /api/iam/check
│   ├── middleware.go       # All hooks (~32 bindings)
│   └── setup.go            # RegisterRoutes + RegisterHooks entry points
├── migrations/
│   └── 1_create_iam_collections.go   # Creates all 9 IAM collections
└── pb_data/                # Runtime data (gitignored)
```

---

## 2. Collection Schemas (9 collections)

### `iam_managed_collections` (base)
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| collection_name | Text | yes | Name of the PB collection to enforce IAM on (unique) |

**Rules:** Superuser-only CRUD. When a record is created here, IAM auto-sets the target collection's PB rules to `@request.auth.id != ''` (authenticated-only) for all operations, blocking unauthenticated access while letting IAM gate authenticated users. When removed, PB rules are restored to `nil`.

### `iam_policies` (base)
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| name | Text | yes | Unique policy name (e.g. "ReadOnlyPosts") |
| description | Text | no | Human-readable description |
| document | JSON | yes | Policy JSON document (validated on create/update) |

**Rules:** Superuser write. Authenticated read.

### `iam_roles` (base)
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| name | Text | yes | Unique role name (e.g. "Editor") |
| description | Text | no | Human-readable description |

### `iam_groups` (base)
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| name | Text | yes | Unique group name (e.g. "Developers") |
| description | Text | no | Human-readable description |

### `iam_role_policies` (base — join)
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| role | Relation → iam_roles | yes | CascadeDelete: true |
| policy | Relation → iam_policies | yes | CascadeDelete: true |

### `iam_user_policies` (base — join)
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| user | Relation → users | yes | CascadeDelete: true |
| policy | Relation → iam_policies | yes | CascadeDelete: true |

### `iam_user_roles` (base — join)
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| user | Relation → users | yes | CascadeDelete: true |
| role | Relation → iam_roles | yes | CascadeDelete: true |

### `iam_group_users` (base — join)
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| group | Relation → iam_groups | yes | CascadeDelete: true |
| user | Relation → users | yes | CascadeDelete: true |

### `iam_group_policies` (base — join)
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| group | Relation → iam_groups | yes | CascadeDelete: true |
| policy | Relation → iam_policies | yes | CascadeDelete: true |

All join tables: superuser-only CRUD. CascadeDelete on all relations. Duplicate-prevention hook on create. Unique composite indexes on both relation fields (defense-in-depth against race conditions).

---

## 3. Policy Document Format

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowReadPosts",
      "effect": "Allow",
      "action": ["collections:posts:read", "collections:posts:list"],
      "resource": ["*"]
    },
    {
      "sid": "DenyDeleteUsers",
      "effect": "Deny",
      "action": ["collections:users:delete"],
      "resource": ["*"]
    },
    {
      "sid": "AllowBillingRefund",
      "effect": "Allow",
      "action": ["custom:billing:refund"],
      "resource": ["*"]
    }
  ]
}
```

**Validation** (enforced via hook on `iam_policies` create/update):
- `version` must be non-empty string
- `statement` must be non-empty array
- Each statement: `effect` is "Allow" or "Deny", non-empty `action[]`, non-empty `resource[]`
- Action strings must contain at least one `:` separator, or be a lone `*` (match-all)
- Invalid → 400 Bad Request

---

## 4. Permission Evaluation

```
Request arrives → IAM middleware hook
    │
    ├─ Collection not in iam_managed_collections? → SKIP (PB rules handle it)
    ├─ Is superuser? → SKIP IAM
    ├─ Is unauthenticated? → SKIP IAM (PB rules handle public access)
    │
    ├─ Build action: "collections:<name>:<operation>"
    │
    ├─ Check cache for user's resolved policies
    │   ├─ Hit → use cached statements
    │   └─ Miss → collect from DB + cache with 60s TTL:
    │       ├─ Direct policies (iam_user_policies)
    │       ├─ Group policies (iam_group_users → iam_group_policies)
    │       └─ Role policies (iam_user_roles → iam_role_policies)
    │
    ├─ Any DENY matching action+resource? → 403 (log: explicit deny + policy name)
    ├─ Any ALLOW matching action+resource? → e.Next() (proceed to PB rules)
    └─ No match → 403 (log: implicit deny, no matching policy)
```

After IAM passes, PB's collection rules run for row-level filtering.

---

## 5. File-by-File Implementation

### `main.go`
- `pocketbase.New()`, register migrate command
- `_ "pocketbase-iam/migrations"` — auto-run migrations on boot
- `iam.RegisterRoutes(app)` — 1 custom endpoint
- `iam.RegisterHooks(app)` — all hooks (enforcement, validation, cache, duplicate prevention, managed-collection rules sync + boot sync)
- `app.Start()`

### `migrations/1_create_iam_collections.go`
- Single migration creating all 9 IAM collections in order (entity tables first, then join tables)
- Uses `core.NewBaseCollection()` with typed fields
- `core.RelationField` with `CascadeDelete: true` on all join tables
- Unique composite indexes on all 5 join tables
- Sets collection rules: superuser-only write, authenticated read on iam_policies/iam_roles/iam_groups
- `downFunc` deletes all 9 in reverse order

### `iam/setup.go`
- `RegisterRoutes(app)` — creates shared cache, registers custom route
- `RegisterHooks(app)` — creates shared cache, registers all hooks + boot sync
- `SyncManagedCollectionRules(app)` — reads `iam_managed_collections`, sets target collections' PB rules to `@request.auth.id != ''`. Called on boot.

### `iam/policy.go`
- `PolicyDocument`, `Statement` structs
- `ParsePolicy(raw any) (*PolicyDocument, error)`
- `ValidatePolicy(doc *PolicyDocument) error`

### `iam/cache.go`
- `PolicyCache` struct wrapping `ttlcache.Cache[string, []Statement]`
- `NewPolicyCache(maxSize int, ttl time.Duration) *PolicyCache`
- `Get(userID) ([]Statement, bool)` — returns if present and not expired
- `Set(userID, []Statement)` — stores with TTL
- `Invalidate(userID)` — removes entry
- `InvalidateUsers(userIDs []string)` — batch remove

### `iam/engine.go`
- `Evaluate(app core.App, cache *PolicyCache, userID, action, resource string) (bool, error)`
- `collectStatements(app, userID) ([]Statement, error)` — collects all unique policy IDs from 3 sources, then batch-fetches in one query
- `evaluateStatements(statements, action, resource) (allowed bool, reason string)` — returns reason for logging
- `IsManagedCollection(app, cache, collectionName) bool` — checks registry (also cached)

### `iam/routes.go`
- `POST /api/iam/check` (requires auth):
  - Body: `{"action": "...", "resource": "..."}`
  - Returns: `{"allowed": true/false}`
  - Logs denial reason server-side

### `iam/middleware.go`
- **Enforcement hooks:** `OnRecordCreateRequest`, `OnRecordUpdateRequest`, `OnRecordDeleteRequest`, `OnRecordViewRequest`, `OnRecordsListRequest`
  - First checks if collection is IAM-managed (via cached registry lookup)
  - Skips superusers and unauthenticated (both bypass IAM)
  - Calls `engine.Evaluate()`, returns 403 or `e.Next()`
  - Logs denial reason via `app.Logger()`
- **Policy validation hook:** `OnRecordCreate/Update("iam_policies")`
  - Validates `document` field, rejects 400 if invalid
- **Duplicate prevention hooks:** `OnRecordCreate` on all join tables
  - Queries for existing record with same field combination, rejects 400 if exists
- **Cache invalidation hooks:** `OnRecordCreate/Update/Delete` on all join tables
  - Smart invalidation: only affected users, not the whole cache
  - `iam_policies` update → find users with that policy via join tables, invalidate them
- **Managed-collection sync hooks:** `OnRecordCreate/Delete("iam_managed_collections")`
  - On create: set target collection's PB rules to `@request.auth.id != ''`, invalidate managed-collection cache
  - On delete: set target collection's PB rules back to `nil`, invalidate cache

### `iam/helpers.go`
- `MatchPattern(pattern, value string) bool` — wildcard matching
- `ActionForOperation(collection, operation string) string`

---

## 6. Implementation Order

1. `go.mod` + `main.go` — scaffold project, PocketBase dependency
2. `migrations/1_create_iam_collections.go` — 9 collections + indexes
3. `iam/policy.go` — types and validation (no deps)
4. `iam/helpers.go` — wildcard matching (no deps)
5. `iam/cache.go` — LRU+TTL cache (no deps)
6. `iam/engine.go` — permission evaluation engine
7. `iam/routes.go` — 1 custom endpoint
8. `iam/middleware.go` — all hooks (enforcement, validation, cache, duplicates, managed-collection sync)
9. `iam/setup.go` — wiring (RegisterRoutes + RegisterHooks)

---

## 7. Verification Plan

1. **Boot**: `go run . serve` — all 9 IAM collections appear in PB admin UI
2. **Opt-in**: Create "posts" collection, verify IAM does NOT enforce on it yet
3. **Register collection**: Add "posts" to `iam_managed_collections` → PB rules auto-set to `@request.auth.id != ''`
4. **Policy validation**: Create iam_policies with invalid JSON → 400
5. **Duplicate prevention**: Try attaching same policy to same user twice → 400
6. **Permission enforcement**: Create Allow policy for `collections:posts:read`, attach to user → user CAN read, CANNOT write
7. **Public access**: Unauthenticated request to IAM-managed collection → blocked by PB rules (`@request.auth.id != ''`). Unregister collection + set PB rule to `""` manually → unauthenticated users can read
8. **Deny overrides Allow**: Attach Deny policy → read blocked
9. **Groups**: Add user to group with policy → user inherits
10. **Multiple roles**: 2 roles with different policies → union applies
11. **Cascade deletes**: Delete role → join records cleaned up
12. **Cache invalidation**: Change attachment → verify immediate effect
13. **Unregister**: Remove "posts" from `iam_managed_collections` → PB rules restored, IAM no longer enforces
14. **Custom actions**: `/api/iam/check` with freeform action+resource
15. **Wildcards**: `*`, `collections:*:read`, `collections:posts:*`, `order:*`
16. **Logs**: Verify denial reasons appear in PB server logs

---

## Future (v2)

- **AssumeRole**: Temporary role sessions with session tokens and expiry
- **Conditions**: Time-based, IP-based, MFA conditions on policy statements
- **Permission boundaries**: Cap maximum permissions for a user regardless of policies
- **Policy versioning**: Track policy document changes over time
- **Audit logging**: Structured audit trail for all IAM decisions and admin actions
- **IAM self-management**: Allow policies to control who can manage IAM entities
