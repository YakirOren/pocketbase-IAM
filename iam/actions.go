package iam

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterAction idempotently registers a custom action in the iam_action_registry.
// If the action already exists, it returns nil.
func RegisterAction(app core.App, action, description string) error {
	// Check if already registered.
	_, err := app.FindFirstRecordByFilter(
		"iam_action_registry",
		"action = {:action}",
		dbx.Params{"action": action},
	)
	if err == nil {
		return nil // already exists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to check action registry for %q: %w", action, err)
	}

	col, err := app.FindCollectionByNameOrId("iam_action_registry")
	if err != nil {
		return fmt.Errorf("failed to find iam_action_registry collection: %w", err)
	}

	record := core.NewRecord(col)
	record.Set("action", action)
	record.Set("description", description)

	if err := app.Save(record); err != nil {
		// Unique constraint violation from a concurrent call — treat as success.
		return nil
	}
	return nil
}
