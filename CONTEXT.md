# Session Context — PocketBase IAM

## Overview
AWS IAM-inspired RBAC system built as a custom PocketBase Go application.
- Repo: `git@github.com:YakirOren/pocketbase-IAM.git`
- Full plan: `PLAN.md` in repo root
- Target: PocketBase v0.36.5

## Status
- Plan finalized and committed. No code written yet.
- Implementation order: policy.go → helpers.go → cache.go → collections.go → engine.go → routes.go → middleware.go

## Key Decisions (from planning session)
- **Go** custom PocketBase binary
- **Opt-in IAM** via `iam_managed_collections` registry — non-registered collections untouched
- **9 collections**: iam_managed_collections, iam_policies, iam_roles, iam_groups, + 5 join tables
- **No custom CRUD routes** — use PB's built-in collection API. Only 1 custom route: `POST /api/iam/check`
- **Unauthenticated requests bypass IAM** — fall through to PB native rules
- **LRU+TTL cache** using `github.com/jellydator/ttlcache/v3` (max 10k entries, 60s TTL)
- **Smart cache invalidation** — targeted per-user, never blow whole cache
- **CascadeDelete** on all join table relations
- **Duplicate prevention hooks** on join tables
- **Policy validation hook** on iam_policies create/update
- **Deny reasons logged server-side**, not returned to client
- **Wildcard matching**: lone `*` matches everything, `*` in segments matches one segment
- **AssumeRole deferred to v2**

## Issues Identified & Resolved During Planning
1. Custom CRUD routes unnecessary → use PB's built-in API
2. Auto-setting PB rules too aggressive → opt-in registry collection
3. Unauthenticated access broken → bypass IAM for unauthenticated
4. Instant lockout on existing apps → opt-in solves this
5. sync.Map unbounded memory → LRU+TTL with ttlcache/v3
6. InvalidateAll() thundering herd → targeted invalidation per affected users
7. No policy validation → hook on iam_policies create/update
8. Join table duplicates → duplicate prevention hooks
9. Orphaned join records → CascadeDelete on all relations
10. Resource ambiguity for CRUD → resource always `*` for hooks, specific only for custom actions

## PocketBase API Reference (v0.36.5)
- Hooks: `OnRecordCreateRequest`, `OnRecordUpdateRequest`, `OnRecordDeleteRequest`, `OnRecordViewRequest`, `OnRecordsListRequest`
- Auth in hooks: `e.Auth` or `e.RequestInfo().Auth`
- Superuser check: `e.HasSuperuserAuth()`
- Routes: registered via `app.OnServe()` → `se.Router.POST("/path", handler)`
- Collections: `core.NewBaseCollection("name")`, fields via `collection.Fields.Add(&core.TextField{...})`
- Records: `app.FindRecordsByFilter()`, `app.Save(record)`, `app.Delete(record)`
- Relation fields: `core.RelationField{Name, CollectionId, CascadeDelete}`

## User Preferences
- Prefers concise approach, no custom routes when PB handles it natively
- Wants thorough critique before implementation
- Short commit messages, no co-author tags
