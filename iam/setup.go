package iam

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Options configures the IAM system.
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

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		CacheMaxSize: 10_000,
		CacheTTL:     60 * time.Second,
	}
}

func (o *Options) validate() error {
	if o.CacheMaxSize <= 0 {
		return fmt.Errorf("CacheMaxSize must be positive, got %d", o.CacheMaxSize)
	}
	if o.CacheTTL <= 0 {
		return fmt.Errorf("CacheTTL must be positive, got %s", o.CacheTTL)
	}
	return nil
}

// Setup initializes the IAM system on a PocketBase app.
// It registers all routes, hooks, enforcement, and the admin dashboard.
// Migrations are auto-registered via package init().
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
	registerEnforcementHooks(app, cache, logger)
	registerPolicyValidationHooks(app)
	registerDuplicatePreventionHooks(app)
	registerCacheInvalidationHooks(app, cache, logger)
	registerManagedCollectionHooks(app, cache, logger)
	registerDashboardRoutes(app)

	// Boot sync: set rules on already-managed collections.
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		if err := SyncManagedCollectionRules(app, logger); err != nil {
			logger.Error("failed to sync managed collection rules on boot", "error", err)
		}
		return se.Next()
	})

	// Graceful shutdown: stop cache cleanup goroutines.
	app.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		cache.Stop()
		return e.Next()
	})

	return nil
}

// SyncManagedCollectionRules reads all iam_managed_collections and sets
// their PocketBase rules to require authentication.
func SyncManagedCollectionRules(app core.App, logger *slog.Logger) error {
	records, err := app.FindRecordsByFilter("iam_managed_collections", "", "", 0, 0)
	if err != nil {
		return fmt.Errorf("failed to query managed collections: %w", err)
	}
	for _, r := range records {
		name := r.GetString("collection_name")
		if err := setCollectionRulesOpen(app, name); err != nil {
			logger.Error("failed to sync rules for collection", "collection", name, "error", err)
		}
	}
	return nil
}

func setCollectionRulesOpen(app core.App, collectionName string) error {
	col, err := app.FindCollectionByNameOrId(collectionName)
	if err != nil {
		return err
	}
	rule := "@request.auth.id != ''"
	col.ListRule = types.Pointer(rule)
	col.ViewRule = types.Pointer(rule)
	col.CreateRule = types.Pointer(rule)
	col.UpdateRule = types.Pointer(rule)
	col.DeleteRule = types.Pointer(rule)
	return app.Save(col)
}

func setCollectionRulesClosed(app core.App, collectionName string) error {
	col, err := app.FindCollectionByNameOrId(collectionName)
	if err != nil {
		return err
	}
	col.ListRule = nil
	col.ViewRule = nil
	col.CreateRule = nil
	col.UpdateRule = nil
	col.DeleteRule = nil
	return app.Save(col)
}
