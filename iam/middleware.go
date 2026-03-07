package iam

import (
	"log/slog"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

const (
	colPolicies           = "iam_policies"
	colRolePolicies       = "iam_role_policies"
	colUserPolicies       = "iam_user_policies"
	colUserRoles          = "iam_user_roles"
	colGroupUsers         = "iam_group_users"
	colGroupPolicies      = "iam_group_policies"
	colManagedCollections = "iam_managed_collections"
)

// registerEnforcementHooks registers hooks for Create, Update, Delete, View, and List
// operations on all collections. Each hook checks if the collection is IAM-managed,
// skips superusers and unauthenticated requests, then evaluates IAM policies.
func registerEnforcementHooks(app core.App, cache *PolicyCache, logger *slog.Logger) {
	enforce := func(collectionName, operation string, hasSuperuserAuth func() bool, auth *core.Record, next func() error) error {
		managed, err := IsManagedCollection(app, cache, collectionName)
		if err != nil {
			return err
		}
		if !managed {
			return next()
		}

		if hasSuperuserAuth() {
			return next()
		}

		if auth == nil {
			return next()
		}

		action := ActionForOperation(collectionName, operation)
		allowed, reason, err := Evaluate(app, cache, auth.Id, action, "*")
		if err != nil {
			logger.Error("IAM evaluation error", "error", err, "user", auth.Id, "action", action)
			return apis.NewApiError(500, "internal error", nil)
		}
		if !allowed {
			logger.Warn("IAM access denied", "user", auth.Id, "action", action, "reason", reason)
			return apis.NewForbiddenError("access denied", nil)
		}
		return next()
	}

	app.OnRecordCreateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return enforce(e.Record.Collection().Name, "create", e.HasSuperuserAuth, e.Auth, e.Next)
	})

	app.OnRecordUpdateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return enforce(e.Record.Collection().Name, "update", e.HasSuperuserAuth, e.Auth, e.Next)
	})

	app.OnRecordDeleteRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return enforce(e.Record.Collection().Name, "delete", e.HasSuperuserAuth, e.Auth, e.Next)
	})

	app.OnRecordViewRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return enforce(e.Record.Collection().Name, "view", e.HasSuperuserAuth, e.Auth, e.Next)
	})

	app.OnRecordsListRequest().BindFunc(func(e *core.RecordsListRequestEvent) error {
		return enforce(e.Collection.Name, "list", e.HasSuperuserAuth, e.Auth, e.Next)
	})
}

// registerPolicyValidationHooks validates the policy document on create and update
// of iam_policies records.
func registerPolicyValidationHooks(app core.App) {
	validate := func(e *core.RecordRequestEvent) error {
		doc, err := ParsePolicy(e.Record.Get("document"))
		if err != nil {
			return e.BadRequestError("invalid policy document: "+err.Error(), nil)
		}
		if err := ValidatePolicy(doc); err != nil {
			return e.BadRequestError("invalid policy document: "+err.Error(), nil)
		}
		return e.Next()
	}

	app.OnRecordCreateRequest(colPolicies).BindFunc(validate)
	app.OnRecordUpdateRequest(colPolicies).BindFunc(validate)
}

// registerDuplicatePreventionHooks prevents duplicate assignments in join tables.
func registerDuplicatePreventionHooks(app core.App) {
	joinTables := []struct {
		collection string
		field1     string
		field2     string
	}{
		{colRolePolicies, "role", "policy"},
		{colUserPolicies, "user", "policy"},
		{colUserRoles, "user", "role"},
		{colGroupUsers, "group", "user"},
		{colGroupPolicies, "group", "policy"},
	}

	for _, jt := range joinTables {
		jt := jt
		app.OnRecordCreateRequest(jt.collection).BindFunc(func(e *core.RecordRequestEvent) error {
			val1 := e.Record.GetString(jt.field1)
			val2 := e.Record.GetString(jt.field2)

			existing, _ := app.FindFirstRecordByFilter(jt.collection,
				jt.field1+" = {:v1} && "+jt.field2+" = {:v2}",
				dbx.Params{"v1": val1, "v2": val2},
			)
			if existing != nil {
				return e.BadRequestError("duplicate assignment", nil)
			}
			return e.Next()
		})
	}
}

// registerCacheInvalidationHooks invalidates cached policies when join tables
// or policy documents change.
func registerCacheInvalidationHooks(app core.App, cache *PolicyCache, logger *slog.Logger) {
	// Helper: invalidate a single user from a record's field.
	// On updates, also invalidates the old user if the field value changed.
	invalidateUserField := func(e *core.RecordEvent, field string) error {
		cache.InvalidateUser(e.Record.GetString(field))
		if orig := e.Record.Original(); orig != nil {
			if oldVal := orig.GetString(field); oldVal != "" && oldVal != e.Record.GetString(field) {
				cache.InvalidateUser(oldVal)
			}
		}
		return e.Next()
	}

	// Helper: find all users related to a field value and invalidate them.
	invalidateUsersForFieldValue := func(collection, field, value string) {
		if value == "" {
			return
		}
		records, err := app.FindRecordsByFilter(collection, field+" = {:val}", "", 0, 0, dbx.Params{"val": value})
		if err != nil {
			logger.Error("cache invalidation query failed", "collection", collection, "field", field, "error", err)
			return
		}
		ids := make([]string, len(records))
		for i, r := range records {
			ids[i] = r.GetString("user")
		}
		cache.InvalidateUsers(ids)
	}

	// Helper: find all users in a group and invalidate them.
	// On updates, also handles the old group if it changed.
	invalidateGroupUsers := func(e *core.RecordEvent) error {
		invalidateUsersForFieldValue(colGroupUsers, "group", e.Record.GetString("group"))
		if orig := e.Record.Original(); orig != nil {
			if oldVal := orig.GetString("group"); oldVal != "" && oldVal != e.Record.GetString("group") {
				invalidateUsersForFieldValue(colGroupUsers, "group", oldVal)
			}
		}
		return e.Next()
	}

	// Helper: find all users with a role and invalidate them.
	// On updates, also handles the old role if it changed.
	invalidateRoleUsers := func(e *core.RecordEvent) error {
		invalidateUsersForFieldValue(colUserRoles, "role", e.Record.GetString("role"))
		if orig := e.Record.Original(); orig != nil {
			if oldVal := orig.GetString("role"); oldVal != "" && oldVal != e.Record.GetString("role") {
				invalidateUsersForFieldValue(colUserRoles, "role", oldVal)
			}
		}
		return e.Next()
	}

	// --- User-level invalidation: iam_user_policies, iam_user_roles, iam_group_users ---
	userJoinTables := []string{colUserPolicies, colUserRoles, colGroupUsers}
	app.OnRecordAfterCreateSuccess(userJoinTables...).BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})
	app.OnRecordAfterUpdateSuccess(userJoinTables...).BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})
	app.OnRecordAfterDeleteSuccess(userJoinTables...).BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})

	// --- Group-level invalidation: iam_group_policies ---
	app.OnRecordAfterCreateSuccess(colGroupPolicies).BindFunc(invalidateGroupUsers)
	app.OnRecordAfterUpdateSuccess(colGroupPolicies).BindFunc(invalidateGroupUsers)
	app.OnRecordAfterDeleteSuccess(colGroupPolicies).BindFunc(invalidateGroupUsers)

	// --- Role-level invalidation: iam_role_policies ---
	app.OnRecordAfterCreateSuccess(colRolePolicies).BindFunc(invalidateRoleUsers)
	app.OnRecordAfterUpdateSuccess(colRolePolicies).BindFunc(invalidateRoleUsers)
	app.OnRecordAfterDeleteSuccess(colRolePolicies).BindFunc(invalidateRoleUsers)

	// --- Policy document changes: iam_policies ---
	invalidatePolicyUsers := func(e *core.RecordEvent) error {
		policyID := e.Record.Id

		seen := make(map[string]struct{})

		// 1. Direct user-policy assignments
		userPolicies, err := app.FindRecordsByFilter(colUserPolicies, "policy = {:pid}", "", 0, 0, dbx.Params{"pid": policyID})
		if err != nil {
			logger.Error("cache invalidation: failed to find user-policy assignments", "policy", policyID, "error", err)
		} else {
			for _, r := range userPolicies {
				seen[r.GetString("user")] = struct{}{}
			}
		}

		// 2. Group-policy -> group-users
		groupPolicies, err := app.FindRecordsByFilter(colGroupPolicies, "policy = {:pid}", "", 0, 0, dbx.Params{"pid": policyID})
		if err != nil {
			logger.Error("cache invalidation: failed to find group-policy assignments", "policy", policyID, "error", err)
		} else {
			for _, gp := range groupPolicies {
				groupUsers, err := app.FindRecordsByFilter(colGroupUsers, "group = {:gid}", "", 0, 0, dbx.Params{"gid": gp.GetString("group")})
				if err != nil {
					logger.Error("cache invalidation: failed to find group users", "group", gp.GetString("group"), "error", err)
				} else {
					for _, gu := range groupUsers {
						seen[gu.GetString("user")] = struct{}{}
					}
				}
			}
		}

		// 3. Role-policy -> user-roles
		rolePolicies, err := app.FindRecordsByFilter(colRolePolicies, "policy = {:pid}", "", 0, 0, dbx.Params{"pid": policyID})
		if err != nil {
			logger.Error("cache invalidation: failed to find role-policy assignments", "policy", policyID, "error", err)
		} else {
			for _, rp := range rolePolicies {
				userRoles, err := app.FindRecordsByFilter(colUserRoles, "role = {:rid}", "", 0, 0, dbx.Params{"rid": rp.GetString("role")})
				if err != nil {
					logger.Error("cache invalidation: failed to find role users", "role", rp.GetString("role"), "error", err)
				} else {
					for _, ur := range userRoles {
						seen[ur.GetString("user")] = struct{}{}
					}
				}
			}
		}

		ids := make([]string, 0, len(seen))
		for id := range seen {
			ids = append(ids, id)
		}
		cache.InvalidateUsers(ids)

		return e.Next()
	}

	app.OnRecordAfterUpdateSuccess(colPolicies).BindFunc(invalidatePolicyUsers)
	app.OnRecordAfterDeleteSuccess(colPolicies).BindFunc(invalidatePolicyUsers)
}

// registerManagedCollectionHooks syncs PocketBase collection rules when collections
// are added to or removed from iam_managed_collections.
func registerManagedCollectionHooks(app core.App, cache *PolicyCache, logger *slog.Logger) {
	// Guard: prevent IAM system collections from being managed (self-lock prevention).
	app.OnRecordCreateRequest(colManagedCollections).BindFunc(func(e *core.RecordRequestEvent) error {
		name := e.Record.GetString("collection_name")
		if strings.HasPrefix(name, "iam_") {
			return e.BadRequestError("cannot manage IAM system collections", nil)
		}
		return e.Next()
	})

	app.OnRecordAfterCreateSuccess(colManagedCollections).BindFunc(func(e *core.RecordEvent) error {
		name := e.Record.GetString("collection_name")
		if err := setCollectionRulesOpen(app, name); err != nil {
			logger.Error("failed to set managed collection rules", "collection", name, "error", err)
		}
		cache.InvalidateManagedCollection(name)
		return e.Next()
	})

	app.OnRecordAfterDeleteSuccess(colManagedCollections).BindFunc(func(e *core.RecordEvent) error {
		name := e.Record.GetString("collection_name")
		if err := setCollectionRulesClosed(app, name); err != nil {
			logger.Error("failed to restore collection rules", "collection", name, "error", err)
		}
		cache.InvalidateManagedCollection(name)
		return e.Next()
	})
}
