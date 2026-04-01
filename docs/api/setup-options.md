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
