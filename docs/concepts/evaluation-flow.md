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
