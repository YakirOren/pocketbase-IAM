# Library Packaging Plan

Turn pocketbase-IAM into a reusable Go library installable via `go get`.

## Consumer API

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

    // Or with custom options:
    // if err := iam.Setup(app, iam.Options{
    //     CacheMaxSize: 5000,
    //     CacheTTL:     30 * time.Second,
    // }); err != nil {
    //     log.Fatalf("Failed to setup IAM: %v", err)
    // }

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

## Options struct

```go
type Options struct {
    // CacheMaxSize is the maximum number of entries in the policy LRU cache.
    // Default: 10000
    CacheMaxSize int

    // CacheTTL is how long cached policy evaluations remain valid.
    // Default: 60s
    CacheTTL time.Duration

    // Logger is an optional structured logger for IAM events.
    // If nil, the PocketBase app's default logger is used.
    Logger *slog.Logger
}

func DefaultOptions() Options {
    return Options{
        CacheMaxSize: 10000,
        CacheTTL:     60 * time.Second,
    }
}
```

## Setup function

```go
func Setup(app core.App, opts Options) error {
    if err := opts.validate(); err != nil {
        return fmt.Errorf("invalid IAM options: %w", err)
    }
    cache := NewPolicyCache(opts.CacheMaxSize, opts.CacheTTL)
    logger := opts.Logger
    if logger == nil {
        logger = app.Logger()
    }
    registerRoutes(app, cache, logger)
    registerHooks(app, cache, logger)
    registerDashboardRoutes(app)
    return nil
}
```

## Implementation Steps

### Step 1: Change Go module path
- `go.mod`: `module pocketbase-iam` → `module github.com/yakiroren/pocketbase-IAM`
- Update all internal imports accordingly

### Step 2: Move migrations into `iam` package
- Move `migrations/1_create_iam_collections.go` → `iam/migration_collections.go`
- Move `migrations/2_create_iam_actions_view.go` → `iam/migration_actions_view.go`
- Keep `init()` functions registering into `core.SystemMigrations`
- Importing `iam` auto-registers migrations — zero config for consumers
- Delete `migrations/` directory entirely
- Remove test-only `posts` migration (not shipped with library)

### Step 3: Embed UI assets with `go:embed`
New file `iam/dashboard.go`:
```go
package iam

import "embed"

//go:embed dashboard/*
var dashboardFS embed.FS
```

Pre-build the UI and output to `iam/dashboard/`. Commit the built assets — consumers get them via `go get` with zero Node.js dependency.

### Step 4: Update UI for subpath serving
- `vite.config.ts`: set `base: "/_/iam/"`, change `outDir: "../iam/dashboard"`
- `App.tsx`: React Router basename = `/_/iam`
- PocketBase client URL stays relative (same host)

### Step 5: Serve dashboard at `/_/iam/`
New file `iam/dashboard_routes.go`:
- Serve embedded SPA via `http.FileServer` on `/_/iam/{path...}`
- SPA fallback: non-existent paths serve `index.html`
- Protected by `apis.RequireSuperuserAuth()`

### Step 6: Refactor setup.go
- Replace `sharedCache`/`sharedCacheOnce` singleton with explicit creation in `Setup()`
- `NewPolicyCache()` takes `maxSize` and `ttl` parameters
- All internal functions receive logger instead of calling `app.Logger()` directly
- Keep `RegisterRoutes`/`RegisterHooks` exported for power users

### Step 7: Update main.go as reference example
Minimal example using `iam.Setup(app, iam.DefaultOptions())`

### Step 8: Clean up
- Delete `migrations/` directory
- Delete `pb_public/` from tracked files
- Update `.gitignore`: remove `pb_public/`, don't ignore `iam/dashboard/`

### Step 9: Update README
Document installation, usage, and dashboard access at `/_/iam/`

### Step 10: Build & verify
```bash
cd ui && npm run build
cd .. && go build ./...
go test ./iam/...
```

## Final file structure

```
github.com/yakiroren/pocketbase-IAM/
├── go.mod                          # module github.com/yakiroren/pocketbase-IAM
├── main.go                         # reference example
├── iam/
│   ├── doc.go
│   ├── setup.go                    # Setup(), DefaultOptions(), RegisterRoutes(), RegisterHooks()
│   ├── policy.go
│   ├── helpers.go
│   ├── cache.go                    # NewPolicyCache(maxSize, ttl)
│   ├── engine.go
│   ├── routes.go
│   ├── middleware.go
│   ├── actions.go
│   ├── migration_collections.go    # moved from migrations/
│   ├── migration_actions_view.go   # moved from migrations/
│   ├── dashboard.go                # //go:embed dashboard/*
│   ├── dashboard_routes.go         # serves /_/iam/ SPA
│   ├── dashboard/                  # pre-built UI assets (committed)
│   │   ├── index.html
│   │   └── assets/...
│   ├── *_test.go
├── ui/                             # source code (development only)
│   └── ...
└── docs/
    └── plans/...
```

## Key decisions

| Decision | Rationale |
|----------|-----------|
| `Setup(app, opts) error` with `DefaultOptions()` | Explicit config, error handling, extensible |
| Migrations via `init()` in `iam` package | Importing the package auto-registers — zero config |
| `go:embed` for UI | No Node.js dependency for consumers |
| Dashboard at `/_/iam/` | Under PB admin namespace, no conflicts |
| Superuser-only dashboard | Consistent with PB admin security model |
| No global singleton cache | Explicit creation in `Setup()`, clean for library use |
| `Logger` option with app fallback | Consumers can plug in their own logger |
