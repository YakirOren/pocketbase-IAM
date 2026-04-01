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
