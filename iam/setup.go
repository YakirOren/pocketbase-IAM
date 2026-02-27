package iam

import (
	"sync"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

var (
	sharedCache     *PolicyCache
	sharedCacheOnce sync.Once
)

func getOrCreateCache() *PolicyCache {
	sharedCacheOnce.Do(func() {
		sharedCache = NewPolicyCache()
	})
	return sharedCache
}

// RegisterRoutes registers the IAM custom API routes.
func RegisterRoutes(app core.App) {
	cache := getOrCreateCache()
	registerRoutes(app, cache)
}

// RegisterHooks registers all IAM hooks: enforcement, validation,
// duplicate prevention, cache invalidation, and managed-collection sync.
func RegisterHooks(app core.App) {
	cache := getOrCreateCache()
	registerEnforcementHooks(app, cache)
	registerPolicyValidationHooks(app)
	registerDuplicatePreventionHooks(app)
	registerCacheInvalidationHooks(app, cache)
	registerManagedCollectionHooks(app, cache)

	// Boot sync: set rules on already-managed collections.
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		if err := SyncManagedCollectionRules(app); err != nil {
			app.Logger().Error("failed to sync managed collection rules on boot", "error", err)
		}
		return se.Next()
	})
}

// SyncManagedCollectionRules reads all iam_managed_collections and sets
// their PocketBase rules to require authentication.
func SyncManagedCollectionRules(app core.App) error {
	records, err := app.FindRecordsByFilter("iam_managed_collections", "", "", 0, 0)
	if err != nil {
		return nil // no records or collection doesn't exist yet
	}
	for _, r := range records {
		name := r.GetString("collection_name")
		if err := setCollectionRulesOpen(app, name); err != nil {
			app.Logger().Error("failed to sync rules for collection", "collection", name, "error", err)
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
