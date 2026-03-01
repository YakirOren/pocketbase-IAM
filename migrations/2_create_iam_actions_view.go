package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	core.SystemMigrations.Register(upCreateIAMActionsView, downCreateIAMActionsView, "2_create_iam_actions_view.go")
}

func upCreateIAMActionsView(app core.App) error {
	// --- 1. iam_action_registry (base collection — must exist before the view) ---
	registry := core.NewBaseCollection("iam_action_registry")
	registry.System = true
	registry.Fields.Add(&core.TextField{
		Name:     "action",
		Required: true,
	})
	registry.Fields.Add(&core.TextField{
		Name: "description",
	})
	registry.ListRule = types.Pointer("@request.auth.id != ''")
	registry.ViewRule = types.Pointer("@request.auth.id != ''")
	registry.AddIndex("idx_action_registry_action", true, "action", "")
	if err := app.Save(registry); err != nil {
		return err
	}

	// --- 2. iam_actions (view collection) ---
	actions := core.NewViewCollection("iam_actions")
	actions.System = true
	actions.ListRule = types.Pointer("@request.auth.id != ''")
	actions.ViewRule = types.Pointer("@request.auth.id != ''")
	actions.ViewQuery = `SELECT * FROM (
  SELECT
    mc.id || op.code as id,
    'collections:' || mc.collection_name || ':' || op.v as action,
    'managed' as source,
    '' as description
  FROM iam_managed_collections mc
  CROSS JOIN (
    SELECT 'list' as v, 'l' as code
    UNION ALL SELECT 'view', 'v'
    UNION ALL SELECT 'create', 'c'
    UNION ALL SELECT 'update', 'u'
    UNION ALL SELECT 'delete', 'd'
  ) op
  UNION ALL
  SELECT ca.id, ca.action, 'registered', ca.description
  FROM iam_action_registry ca
)`
	if err := app.Save(actions); err != nil {
		return err
	}

	return nil
}

func downCreateIAMActionsView(app core.App) error {
	// Delete view first (depends on iam_action_registry), then base.
	names := []string{
		"iam_actions",
		"iam_action_registry",
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
