# PocketBase IAM — Implementation Plan

## Context

Implement the AWS IAM-inspired RBAC system described in `PLAN.md`. This is a greenfield Go project — no code exists yet. The 9 IAM collections will be created via **PocketBase migrations** (user's choice). The full spec lives in `PLAN.md`; this plan covers execution details.

---

## Project Structure

```
pocketbase-IAM/
├── main.go
├── go.mod / go.sum
├── iam/
│   ├── policy.go         # Policy document types + parsing + validation
│   ├── helpers.go        # Wildcard matching, action string builder
│   ├── cache.go          # LRU+TTL cache wrapping ttlcache/v3
│   ├── engine.go         # Permission evaluation (collect + evaluate)
│   ├── routes.go         # POST /api/iam/check
│   ├── middleware.go     # All hooks (~32 bindings)
│   └── setup.go          # RegisterRoutes + RegisterHooks entry points
├── migrations/
│   └── 1_create_iam_collections.go
└── pb_data/              # Runtime (gitignored)
```

---

## Implementation Steps

### Step 1: Scaffold — `go.mod` + `main.go`

**`go.mod`**: module `pocketbase-iam`, require `github.com/pocketbase/pocketbase` v0.36.5 + `github.com/jellydator/ttlcache/v3`.

**`main.go`**:
```go
package main

import (
    "log"
    "os"
    "strings"

    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/plugins/migratecmd"

    _ "pocketbase-iam/migrations"
    "pocketbase-iam/iam"
)

func main() {
    app := pocketbase.New()

    isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
    migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
        Automigrate: isGoRun,
    })

    iam.RegisterRoutes(app)
    iam.RegisterHooks(app)

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### Step 2: Migration — `migrations/1_create_iam_collections.go`

Single migration file creating all 9 collections. Order matters — entity tables first, then join tables.

**Creation order:**
1. `iam_managed_collections` (Text field: `collection_name`, unique)
2. `iam_policies` (Text: `name` unique, Text: `description`, JSON: `document`)
3. `iam_roles` (Text: `name` unique, Text: `description`)
4. `iam_groups` (Text: `name` unique, Text: `description`)
5. `iam_role_policies` (Relation: `role` → iam_roles, Relation: `policy` → iam_policies, both CascadeDelete)
6. `iam_user_policies` (Relation: `user` → users, Relation: `policy` → iam_policies, both CascadeDelete)
7. `iam_user_roles` (Relation: `user` → users, Relation: `role` → iam_roles, both CascadeDelete)
8. `iam_group_users` (Relation: `group` → iam_groups, Relation: `user` → users, both CascadeDelete)
9. `iam_group_policies` (Relation: `group` → iam_groups, Relation: `policy` → iam_policies, both CascadeDelete)

**Rules:**
- All 9 collections: all rules = `nil` (superuser-only by default)
- `iam_policies`, `iam_roles`, `iam_groups`: `ListRule` and `ViewRule` = `types.Pointer("@request.auth.id != ''")`  (authenticated read)

**Key API calls:**
- `core.NewBaseCollection("name")` to create
- `collection.Fields.Add(&core.TextField{Name: "...", Required: true})` for fields
- `app.FindCollectionByNameOrId("users")` to get collection ID for relations
- `&core.RelationField{Name: "user", CollectionId: usersCol.Id, CascadeDelete: true, Required: true}`
- `collection.AddIndex("idx_...", true, "field1, field2", "")` for unique constraints
- `app.Save(collection)` to persist
- `downFunc` deletes all 9 in reverse order via `app.Delete()`

### Step 3: `iam/policy.go` — Types + Validation

```go
type PolicyDocument struct {
    Version   string      `json:"version"`
    Statement []Statement `json:"statement"`
}

type Statement struct {
    SID      string   `json:"sid"`
    Effect   string   `json:"effect"`
    Action   []string `json:"action"`
    Resource []string `json:"resource"`
}

func ParsePolicy(raw any) (*PolicyDocument, error)    // JSON unmarshal from interface{}
func ValidatePolicy(doc *PolicyDocument) error         // version non-empty, statements valid
```

Validation rules:
- `version` non-empty string
- `statement` non-empty array
- Each: `effect` = "Allow" or "Deny", `action` non-empty (each contains `:` or is `*`), `resource` non-empty

### Step 4: `iam/helpers.go` — Wildcard Matching

```go
func MatchPattern(pattern, value string) bool  // lone "*" = match all; otherwise segment-by-segment
func ActionForOperation(collectionName, operation string) string  // "collections:<name>:<op>"
```

### Step 5: `iam/cache.go` — LRU+TTL Cache

Wraps `ttlcache.Cache` with two sub-caches:
- `policies`: `string → []Statement` (10k entries, 60s TTL)
- `managed`: `string → bool` (1k entries, 5min TTL)

Methods: `GetPolicies`, `SetPolicies`, `InvalidateUser`, `InvalidateUsers`, `IsManagedCollection`, `SetManagedCollection`, `InvalidateManagedCollection`.

### Step 6: `iam/engine.go` — Permission Evaluation

```go
func Evaluate(app core.App, cache *PolicyCache, userID, action, resource string) (bool, error)
func IsManagedCollection(app core.App, cache *PolicyCache, collectionName string) (bool, error)
```

Internal:
- `collectStatements(app, userID)` — collect all unique policy IDs from 3 sources, then batch-fetch policies in one query
- `evaluateStatements(stmts, action, resource)` — deny-overrides-allow, returns `(allowed, reason)`

### Step 7: `iam/routes.go` — Custom Endpoint

`POST /api/iam/check` (requires auth via `apis.RequireAuth()`):
- Body: `{"action": "...", "resource": "..."}`  (resource defaults to `"*"`)
- Response: `{"allowed": true/false}`

### Step 8: `iam/middleware.go` — All Hooks (~32 bindings)

| Category | Hooks | Details |
|---|---|---|
| **Enforcement** (5) | `OnRecord{Create,Update,Delete,View}Request` + `OnRecordsListRequest` | Global, check managed → superuser → unauth → Evaluate |
| **Policy validation** (2) | `OnRecordCreateRequest("iam_policies")` + Update | Parse + validate `document` field |
| **Duplicate prevention** (5) | `OnRecordCreateRequest` on each join table | Query for existing combo, reject 400 |
| **Cache invalidation** (17) | `OnRecordAfter{Create,Update,Delete}Success` on join tables + `iam_policies` | Smart per-user invalidation (fires for cascade deletes too) |
| **Managed-collection sync** (2+1) | `OnRecordAfter{Create,Delete}Success("iam_managed_collections")` + boot sync | Set rules to `@request.auth.id != ''` / nil, invalidate cache |

Enforcement flow: not managed → skip; superuser → skip; unauthenticated → skip; else Evaluate → 403 or proceed.

### Step 9: `iam/setup.go` — Wiring

```go
func RegisterRoutes(app core.App)  // creates shared cache, calls registerRoutes
func RegisterHooks(app core.App)   // creates shared cache, calls all register*Hooks + boot sync
```

Package-level `sharedCache` singleton (single-process PocketBase).

---

## Issues Found & Fixes

### 1. Unauthenticated Gets Full Access to Managed Collections (Critical)

**Problem:** Original PLAN.md sets managed collection rules to `""` (open to everyone). PB evaluates rules *before* hooks run. Since hooks skip unauthenticated requests, unauthenticated users get full unrestricted CRUD on all managed collections.

**Fix:** Set managed collection rules to `"@request.auth.id != ''"` instead of `""`. This blocks unauthenticated access at the PB rule layer while letting IAM be the sole gatekeeper for authenticated users. Apply this in both `setCollectionRulesOpen()` and `SyncManagedCollectionRules()`.

### 2. Join Tables Need Unique Composite Indexes (Medium)

**Problem:** Duplicate-prevention hooks have a race condition — two simultaneous requests can both pass the "does it exist?" check before either writes. Application-level checks alone are insufficient.

**Fix:** Add unique composite indexes on all 5 join tables in the migration:
```
iam_role_policies:  unique(role, policy)
iam_user_policies:  unique(user, policy)
iam_user_roles:     unique(user, role)
iam_group_users:    unique(group, user)
iam_group_policies: unique(group, policy)
```
Keep the duplicate-prevention hooks for friendly 400 error messages; the DB index is defense-in-depth.

### 3. collectStatements N+1 Queries (Medium)

**Problem:** Each join table record triggers an individual policy fetch. ~35 DB queries on cache miss for a user with many roles/groups.

**Fix:** Collect all unique policy IDs across all 3 sources first (using a `map[string]struct{}`), then batch-fetch all policies in a single `FindRecordsByFilter` call. This also naturally deduplicates policies attached via multiple paths.

### 4. Engine Must Handle `sql.ErrNoRows` (Low)

**Problem:** `FindFirstRecordByFilter` returns `sql.ErrNoRows` when no record matches. `IsManagedCollection` must distinguish "not found" (cache as false) from actual DB errors (return error).

### 5. Cascade Delete Hooks — Confirmed Working

PocketBase **does** fire delete hooks for cascade-deleted records (recursive invocation). No extra workaround needed. Caveat from PB docs: "avoid global mutex locks inside hook handlers because they can be invoked recursively via cascade delete." Our hooks use ttlcache (internally thread-safe) so this is fine.

---

## Verification

1. `go build && ./pocketbase-iam serve` — 9 collections in admin UI
2. Create "posts" collection, verify IAM does NOT enforce
3. Register "posts" in `iam_managed_collections` → rules become `"@request.auth.id != ''"`
4. Authenticated request without policy → 403 (implicit deny)
5. Invalid policy document → 400
6. Attach Allow policy for `collections:posts:read` → user can read, not write
7. Duplicate join record → 400
8. Attach Deny for same action → 403 (explicit deny overrides)
9. Groups + roles → user inherits policies
10. Detach policy → next request reflects change (cache invalidated)
11. Unregister collection → rules restored to `nil`
12. `/api/iam/check` with custom action → `{"allowed": true/false}`
13. Unauthenticated + superuser bypass IAM
14. Delete role → cascade cleans join records
