package iam

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Evaluate checks if a user is allowed to perform an action on a resource.
// It collects all policy statements from direct, group, and role attachments,
// then applies deny-overrides-allow evaluation.
func Evaluate(app core.App, cache *PolicyCache, userID, action, resource string) (bool, error) {
	stmts, ok := cache.GetPolicies(userID)
	if !ok {
		var err error
		stmts, err = collectStatements(app, userID)
		if err != nil {
			return false, fmt.Errorf("failed to collect statements: %w", err)
		}
		cache.SetPolicies(userID, stmts)
	}

	allowed, _ := evaluateStatements(stmts, action, resource)
	return allowed, nil
}

// IsManagedCollection checks if a collection is IAM-managed
// (registered in iam_managed_collections).
func IsManagedCollection(app core.App, cache *PolicyCache, collectionName string) (bool, error) {
	managed, found := cache.GetManagedCollection(collectionName)
	if found {
		return managed, nil
	}

	_, err := app.FindFirstRecordByFilter(
		"iam_managed_collections",
		"collection_name = {:name}",
		dbx.Params{"name": collectionName},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cache.SetManagedCollection(collectionName, false)
			return false, nil
		}
		return false, fmt.Errorf("failed to check managed collection %q: %w", collectionName, err)
	}

	cache.SetManagedCollection(collectionName, true)
	return true, nil
}

// collectStatements gathers all policy statements applicable to a user
// from direct policies, group policies, and role policies. It batch-fetches
// all unique policy IDs in a single query to avoid N+1.
func collectStatements(app core.App, userID string) ([]Statement, error) {
	policyIDs := make(map[string]struct{})

	// 1. Direct policies (iam_user_policies)
	directRecords, err := app.FindRecordsByFilter(
		"iam_user_policies",
		"user = {:uid}",
		"", 0, 0,
		dbx.Params{"uid": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user policies: %w", err)
	}
	for _, r := range directRecords {
		policyIDs[r.GetString("policy")] = struct{}{}
	}

	// 2. Group policies (iam_group_users → iam_group_policies)
	groupUserRecords, err := app.FindRecordsByFilter(
		"iam_group_users",
		"user = {:uid}",
		"", 0, 0,
		dbx.Params{"uid": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group memberships: %w", err)
	}
	for _, gu := range groupUserRecords {
		groupID := gu.GetString("group")
		groupPolicyRecords, err := app.FindRecordsByFilter(
			"iam_group_policies",
			"group = {:gid}",
			"", 0, 0,
			dbx.Params{"gid": groupID},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch group policies for group %s: %w", groupID, err)
		}
		for _, gp := range groupPolicyRecords {
			policyIDs[gp.GetString("policy")] = struct{}{}
		}
	}

	// 3. Role policies (iam_user_roles → iam_role_policies)
	roleRecords, err := app.FindRecordsByFilter(
		"iam_user_roles",
		"user = {:uid}",
		"", 0, 0,
		dbx.Params{"uid": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user roles: %w", err)
	}
	for _, ur := range roleRecords {
		roleID := ur.GetString("role")
		rolePolicyRecords, err := app.FindRecordsByFilter(
			"iam_role_policies",
			"role = {:rid}",
			"", 0, 0,
			dbx.Params{"rid": roleID},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch role policies for role %s: %w", roleID, err)
		}
		for _, rp := range rolePolicyRecords {
			policyIDs[rp.GetString("policy")] = struct{}{}
		}
	}

	// No policies found — return empty
	if len(policyIDs) == 0 {
		return nil, nil
	}

	// Batch-fetch all unique policies in one query
	filter := buildIDFilter(policyIDs)
	policyRecords, err := app.FindRecordsByFilter(
		"iam_policies",
		filter,
		"", 0, 0,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to batch-fetch policies: %w", err)
	}

	// Parse all policy documents and collect statements
	var allStatements []Statement
	for _, pr := range policyRecords {
		doc, err := ParsePolicy(pr.Get("document"))
		if err != nil {
			app.Logger().Warn("skipping invalid policy document",
				"policyID", pr.Id,
				"error", err,
			)
			continue
		}
		allStatements = append(allStatements, doc.Statement...)
	}

	return allStatements, nil
}

// buildIDFilter constructs a PocketBase filter expression to match multiple IDs.
// e.g. "id = 'abc' || id = 'def'"
func buildIDFilter(ids map[string]struct{}) string {
	parts := make([]string, 0, len(ids))
	for id := range ids {
		parts = append(parts, fmt.Sprintf("id = '%s'", id))
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += " || " + parts[i]
	}
	return result
}

// evaluateStatements applies deny-overrides-allow evaluation to a set of statements.
// It returns whether the action is allowed and a human-readable reason for logging.
func evaluateStatements(stmts []Statement, action, resource string) (bool, string) {
	// Phase 1: Check all Deny statements first
	for _, stmt := range stmts {
		if stmt.Effect != "Deny" {
			continue
		}
		for _, a := range stmt.Action {
			if !MatchPattern(a, action) {
				continue
			}
			for _, r := range stmt.Resource {
				if MatchPattern(r, resource) {
					return false, fmt.Sprintf("explicit deny by %s", stmt.SID)
				}
			}
		}
	}

	// Phase 2: Check all Allow statements
	for _, stmt := range stmts {
		if stmt.Effect != "Allow" {
			continue
		}
		for _, a := range stmt.Action {
			if !MatchPattern(a, action) {
				continue
			}
			for _, r := range stmt.Resource {
				if MatchPattern(r, resource) {
					return true, fmt.Sprintf("allowed by %s", stmt.SID)
				}
			}
		}
	}

	// Phase 3: Implicit deny
	return false, "implicit deny: no matching policy"
}
