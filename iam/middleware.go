package iam

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// registerEnforcementHooks registers hooks for Create, Update, Delete, View, and List
// operations on all collections. Each hook checks if the collection is IAM-managed,
// skips superusers and unauthenticated requests, then evaluates IAM policies.
func registerEnforcementHooks(app core.App, cache *PolicyCache) {
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
		allowed, err := Evaluate(app, cache, auth.Id, action, "*")
		if err != nil {
			app.Logger().Error("IAM evaluation error", "error", err, "user", auth.Id, "action", action)
			return apis.NewApiError(500, "internal error", nil)
		}
		if !allowed {
			app.Logger().Warn("IAM access denied", "user", auth.Id, "action", action)
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

	app.OnRecordCreateRequest("iam_policies").BindFunc(validate)
	app.OnRecordUpdateRequest("iam_policies").BindFunc(validate)
}

// registerDuplicatePreventionHooks prevents duplicate assignments in join tables.
func registerDuplicatePreventionHooks(app core.App) {
	joinTables := []struct {
		collection string
		field1     string
		field2     string
	}{
		{"iam_role_policies", "role", "policy"},
		{"iam_user_policies", "user", "policy"},
		{"iam_user_roles", "user", "role"},
		{"iam_group_users", "group", "user"},
		{"iam_group_policies", "group", "policy"},
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
func registerCacheInvalidationHooks(app core.App, cache *PolicyCache) {
	// Helper: invalidate a single user from a record's field.
	invalidateUserField := func(e *core.RecordEvent, field string) error {
		cache.InvalidateUser(e.Record.GetString(field))
		return e.Next()
	}

	// Helper: find all users in a group and invalidate them.
	invalidateGroupUsers := func(e *core.RecordEvent) error {
		groupID := e.Record.GetString("group")
		records, err := app.FindRecordsByFilter("iam_group_users", "group = {:gid}", "", 0, 0, dbx.Params{"gid": groupID})
		if err == nil {
			ids := make([]string, len(records))
			for i, r := range records {
				ids[i] = r.GetString("user")
			}
			cache.InvalidateUsers(ids)
		}
		return e.Next()
	}

	// Helper: find all users with a role and invalidate them.
	invalidateRoleUsers := func(e *core.RecordEvent) error {
		roleID := e.Record.GetString("role")
		records, err := app.FindRecordsByFilter("iam_user_roles", "role = {:rid}", "", 0, 0, dbx.Params{"rid": roleID})
		if err == nil {
			ids := make([]string, len(records))
			for i, r := range records {
				ids[i] = r.GetString("user")
			}
			cache.InvalidateUsers(ids)
		}
		return e.Next()
	}

	// --- User-level invalidation: iam_user_policies ---
	app.OnRecordAfterCreateSuccess("iam_user_policies").BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})
	app.OnRecordAfterUpdateSuccess("iam_user_policies").BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})
	app.OnRecordAfterDeleteSuccess("iam_user_policies").BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})

	// --- User-level invalidation: iam_user_roles ---
	app.OnRecordAfterCreateSuccess("iam_user_roles").BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})
	app.OnRecordAfterUpdateSuccess("iam_user_roles").BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})
	app.OnRecordAfterDeleteSuccess("iam_user_roles").BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})

	// --- User-level invalidation: iam_group_users ---
	app.OnRecordAfterCreateSuccess("iam_group_users").BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})
	app.OnRecordAfterUpdateSuccess("iam_group_users").BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})
	app.OnRecordAfterDeleteSuccess("iam_group_users").BindFunc(func(e *core.RecordEvent) error {
		return invalidateUserField(e, "user")
	})

	// --- Group-level invalidation: iam_group_policies ---
	app.OnRecordAfterCreateSuccess("iam_group_policies").BindFunc(func(e *core.RecordEvent) error {
		return invalidateGroupUsers(e)
	})
	app.OnRecordAfterUpdateSuccess("iam_group_policies").BindFunc(func(e *core.RecordEvent) error {
		return invalidateGroupUsers(e)
	})
	app.OnRecordAfterDeleteSuccess("iam_group_policies").BindFunc(func(e *core.RecordEvent) error {
		return invalidateGroupUsers(e)
	})

	// --- Role-level invalidation: iam_role_policies ---
	app.OnRecordAfterCreateSuccess("iam_role_policies").BindFunc(func(e *core.RecordEvent) error {
		return invalidateRoleUsers(e)
	})
	app.OnRecordAfterUpdateSuccess("iam_role_policies").BindFunc(func(e *core.RecordEvent) error {
		return invalidateRoleUsers(e)
	})
	app.OnRecordAfterDeleteSuccess("iam_role_policies").BindFunc(func(e *core.RecordEvent) error {
		return invalidateRoleUsers(e)
	})

	// --- Policy document changes: iam_policies ---
	invalidatePolicyUsers := func(e *core.RecordEvent) error {
		policyID := e.Record.Id

		seen := make(map[string]struct{})

		// 1. Direct user-policy assignments
		userPolicies, err := app.FindRecordsByFilter("iam_user_policies", "policy = {:pid}", "", 0, 0, dbx.Params{"pid": policyID})
		if err == nil {
			for _, r := range userPolicies {
				seen[r.GetString("user")] = struct{}{}
			}
		}

		// 2. Group-policy → group-users
		groupPolicies, err := app.FindRecordsByFilter("iam_group_policies", "policy = {:pid}", "", 0, 0, dbx.Params{"pid": policyID})
		if err == nil {
			for _, gp := range groupPolicies {
				groupUsers, err := app.FindRecordsByFilter("iam_group_users", "group = {:gid}", "", 0, 0, dbx.Params{"gid": gp.GetString("group")})
				if err == nil {
					for _, gu := range groupUsers {
						seen[gu.GetString("user")] = struct{}{}
					}
				}
			}
		}

		// 3. Role-policy → user-roles
		rolePolicies, err := app.FindRecordsByFilter("iam_role_policies", "policy = {:pid}", "", 0, 0, dbx.Params{"pid": policyID})
		if err == nil {
			for _, rp := range rolePolicies {
				userRoles, err := app.FindRecordsByFilter("iam_user_roles", "role = {:rid}", "", 0, 0, dbx.Params{"rid": rp.GetString("role")})
				if err == nil {
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

	app.OnRecordAfterUpdateSuccess("iam_policies").BindFunc(invalidatePolicyUsers)
	app.OnRecordAfterDeleteSuccess("iam_policies").BindFunc(invalidatePolicyUsers)
}

// registerManagedCollectionHooks syncs PocketBase collection rules when collections
// are added to or removed from iam_managed_collections.
func registerManagedCollectionHooks(app core.App, cache *PolicyCache) {
	app.OnRecordAfterCreateSuccess("iam_managed_collections").BindFunc(func(e *core.RecordEvent) error {
		name := e.Record.GetString("collection_name")
		if err := setCollectionRulesOpen(app, name); err != nil {
			app.Logger().Error("failed to set managed collection rules", "collection", name, "error", err)
		}
		cache.InvalidateManagedCollection(name)
		return e.Next()
	})

	app.OnRecordAfterDeleteSuccess("iam_managed_collections").BindFunc(func(e *core.RecordEvent) error {
		name := e.Record.GetString("collection_name")
		if err := setCollectionRulesClosed(app, name); err != nil {
			app.Logger().Error("failed to restore collection rules", "collection", name, "error", err)
		}
		cache.InvalidateManagedCollection(name)
		return e.Next()
	})
}
