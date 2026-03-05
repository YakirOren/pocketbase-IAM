package iam

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

// PolicyCache provides LRU+TTL caching for resolved user policy statements
// and managed collection lookups.
type PolicyCache struct {
	policies *ttlcache.Cache[string, []Statement]
	managed  *ttlcache.Cache[string, bool]
}

// NewPolicyCache creates a new cache with two sub-caches:
//   - policies: 10,000 max entries, 60s TTL (resolved user statements)
//   - managed: 1,000 max entries, 5min TTL (collection managed status)
func NewPolicyCache() *PolicyCache {
	policies := ttlcache.New[string, []Statement](
		ttlcache.WithCapacity[string, []Statement](10_000),
		ttlcache.WithTTL[string, []Statement](60*time.Second),
	)
	go policies.Start()

	managed := ttlcache.New[string, bool](
		ttlcache.WithCapacity[string, bool](1_000),
		ttlcache.WithTTL[string, bool](5*time.Minute),
	)
	go managed.Start()

	return &PolicyCache{
		policies: policies,
		managed:  managed,
	}
}

// GetPolicies returns the cached resolved statements for a user.
func (c *PolicyCache) GetPolicies(userID string) ([]Statement, bool) {
	item := c.policies.Get(userID)
	if item == nil {
		return nil, false
	}
	return item.Value(), true
}

// SetPolicies caches the resolved statements for a user.
func (c *PolicyCache) SetPolicies(userID string, stmts []Statement) {
	c.policies.Set(userID, stmts, ttlcache.DefaultTTL)
}

// InvalidateUser removes a single user's cached policies.
func (c *PolicyCache) InvalidateUser(userID string) {
	c.policies.Delete(userID)
}

// InvalidateUsers removes cached policies for multiple users.
func (c *PolicyCache) InvalidateUsers(userIDs []string) {
	for _, id := range userIDs {
		c.policies.Delete(id)
	}
}

// GetManagedCollection returns whether a collection is IAM-managed from cache.
func (c *PolicyCache) GetManagedCollection(name string) (bool, bool) {
	item := c.managed.Get(name)
	if item == nil {
		return false, false
	}
	return item.Value(), true
}

// SetManagedCollection caches the managed status of a collection.
func (c *PolicyCache) SetManagedCollection(name string, managed bool) {
	c.managed.Set(name, managed, ttlcache.DefaultTTL)
}

// InvalidateManagedCollection removes a collection's cached managed status.
func (c *PolicyCache) InvalidateManagedCollection(name string) {
	c.managed.Delete(name)
}

// Stop cleanly shuts down the cache's background cleanup goroutines.
func (c *PolicyCache) Stop() {
	c.policies.Stop()
	c.managed.Stop()
}
