package iam

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	// The migration name is kept as the original filename for backwards compatibility
	// with existing databases that have already applied this migration.
	core.SystemMigrations.Register(upCreateIAMCollections, downCreateIAMCollections, "1_create_iam_collections.go")
}

func upCreateIAMCollections(app core.App) error {
	// --- 1. iam_managed_collections ---
	managedCols := core.NewBaseCollection("iam_managed_collections")
	managedCols.System = true
	managedCols.Fields.Add(&core.TextField{
		Name:     "collection_name",
		Required: true,
	})
	managedCols.AddIndex("idx_managed_cols_name", true, "collection_name", "")
	if err := app.Save(managedCols); err != nil {
		return err
	}

	// --- 2. iam_policies ---
	policies := core.NewBaseCollection("iam_policies")
	policies.System = true
	policies.Fields.Add(&core.TextField{
		Name:     "name",
		Required: true,
	})
	policies.Fields.Add(&core.TextField{
		Name: "description",
	})
	policies.Fields.Add(&core.JSONField{
		Name:     "document",
		Required: true,
	})
	policies.ListRule = types.Pointer("@request.auth.id != ''")
	policies.ViewRule = types.Pointer("@request.auth.id != ''")
	policies.AddIndex("idx_policies_name", true, "name", "")
	if err := app.Save(policies); err != nil {
		return err
	}

	// --- 3. iam_roles ---
	roles := core.NewBaseCollection("iam_roles")
	roles.System = true
	roles.Fields.Add(&core.TextField{
		Name:     "name",
		Required: true,
	})
	roles.Fields.Add(&core.TextField{
		Name: "description",
	})
	roles.ListRule = types.Pointer("@request.auth.id != ''")
	roles.ViewRule = types.Pointer("@request.auth.id != ''")
	roles.AddIndex("idx_roles_name", true, "name", "")
	if err := app.Save(roles); err != nil {
		return err
	}

	// --- 4. iam_groups ---
	groups := core.NewBaseCollection("iam_groups")
	groups.System = true
	groups.Fields.Add(&core.TextField{
		Name:     "name",
		Required: true,
	})
	groups.Fields.Add(&core.TextField{
		Name: "description",
	})
	groups.ListRule = types.Pointer("@request.auth.id != ''")
	groups.ViewRule = types.Pointer("@request.auth.id != ''")
	groups.AddIndex("idx_groups_name", true, "name", "")
	if err := app.Save(groups); err != nil {
		return err
	}

	// Look up collection IDs for relations.
	usersCol, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}

	// --- 5. iam_role_policies ---
	rolePolicies := core.NewBaseCollection("iam_role_policies")
	rolePolicies.System = true
	rolePolicies.Fields.Add(&core.RelationField{
		Name:          "role",
		CollectionId:  roles.Id,
		CascadeDelete: true,
		Required:      true,
	})
	rolePolicies.Fields.Add(&core.RelationField{
		Name:          "policy",
		CollectionId:  policies.Id,
		CascadeDelete: true,
		Required:      true,
	})
	rolePolicies.AddIndex("idx_role_policies_unique", true, "role, policy", "")
	if err := app.Save(rolePolicies); err != nil {
		return err
	}

	// --- 6. iam_user_policies ---
	userPolicies := core.NewBaseCollection("iam_user_policies")
	userPolicies.System = true
	userPolicies.Fields.Add(&core.RelationField{
		Name:          "user",
		CollectionId:  usersCol.Id,
		CascadeDelete: true,
		Required:      true,
	})
	userPolicies.Fields.Add(&core.RelationField{
		Name:          "policy",
		CollectionId:  policies.Id,
		CascadeDelete: true,
		Required:      true,
	})
	userPolicies.AddIndex("idx_user_policies_unique", true, "user, policy", "")
	if err := app.Save(userPolicies); err != nil {
		return err
	}

	// --- 7. iam_user_roles ---
	userRoles := core.NewBaseCollection("iam_user_roles")
	userRoles.System = true
	userRoles.Fields.Add(&core.RelationField{
		Name:          "user",
		CollectionId:  usersCol.Id,
		CascadeDelete: true,
		Required:      true,
	})
	userRoles.Fields.Add(&core.RelationField{
		Name:          "role",
		CollectionId:  roles.Id,
		CascadeDelete: true,
		Required:      true,
	})
	userRoles.AddIndex("idx_user_roles_unique", true, "user, role", "")
	if err := app.Save(userRoles); err != nil {
		return err
	}

	// --- 8. iam_group_users ---
	groupUsers := core.NewBaseCollection("iam_group_users")
	groupUsers.System = true
	groupUsers.Fields.Add(&core.RelationField{
		Name:          "group",
		CollectionId:  groups.Id,
		CascadeDelete: true,
		Required:      true,
	})
	groupUsers.Fields.Add(&core.RelationField{
		Name:          "user",
		CollectionId:  usersCol.Id,
		CascadeDelete: true,
		Required:      true,
	})
	groupUsers.AddIndex("idx_group_users_unique", true, "`group`, user", "")
	if err := app.Save(groupUsers); err != nil {
		return err
	}

	// --- 9. iam_group_policies ---
	groupPolicies := core.NewBaseCollection("iam_group_policies")
	groupPolicies.System = true
	groupPolicies.Fields.Add(&core.RelationField{
		Name:          "group",
		CollectionId:  groups.Id,
		CascadeDelete: true,
		Required:      true,
	})
	groupPolicies.Fields.Add(&core.RelationField{
		Name:          "policy",
		CollectionId:  policies.Id,
		CascadeDelete: true,
		Required:      true,
	})
	groupPolicies.AddIndex("idx_group_policies_unique", true, "`group`, policy", "")
	if err := app.Save(groupPolicies); err != nil {
		return err
	}

	return nil
}

func downCreateIAMCollections(app core.App) error {
	// Delete in reverse order: join tables first, then entity tables.
	names := []string{
		"iam_group_policies",
		"iam_group_users",
		"iam_user_roles",
		"iam_user_policies",
		"iam_role_policies",
		"iam_groups",
		"iam_roles",
		"iam_policies",
		"iam_managed_collections",
	}
	for _, name := range names {
		col, err := app.FindCollectionByNameOrId(name)
		if err != nil {
			return err
		}
		if err := app.Delete(col); err != nil {
			return err
		}
	}
	return nil
}
