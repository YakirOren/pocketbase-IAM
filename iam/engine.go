package iam

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Evaluate checks if a user is allowed to perform an action on a resource.
// It collects all policy statements from direct, group, and role attachments,
// then applies deny-overrides-allow evaluation.
// Returns (allowed, reason, error) where reason describes the evaluation outcome.
func Evaluate(app core.App, cache *PolicyCache, userID, action, resource string) (bool, string, error) {
	stmts, ok := cache.GetPolicies(userID)
	if !ok {
		var err error
		stmts, err = collectStatements(app, userID)
		if err != nil {
			return false, "", fmt.Errorf("failed to collect statements: %w", err)
		}
		cache.SetPolicies(userID, stmts)
	}

	allowed, reason := evaluateStatements(stmts, action, resource)
	return allowed, reason, nil
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
// intermediate and final queries to avoid N+1.
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

	// 2. Group policies (iam_group_users → iam_group_policies) — batched
	groupUserRecords, err := app.FindRecordsByFilter(
		"iam_group_users",
		"user = {:uid}",
		"", 0, 0,
		dbx.Params{"uid": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group memberships: %w", err)
	}
	if len(groupUserRecords) > 0 {
		groupIDs := make([]string, len(groupUserRecords))
		for i, gu := range groupUserRecords {
			groupIDs[i] = gu.GetString("group")
		}
		filter, params := buildInFilter("group", groupIDs)
		groupPolicyRecords, err := app.FindRecordsByFilter(
			"iam_group_policies",
			filter,
			"", 0, 0,
			params,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch group policies: %w", err)
		}
		for _, gp := range groupPolicyRecords {
			policyIDs[gp.GetString("policy")] = struct{}{}
		}
	}

	// 3. Role policies (iam_user_roles → iam_role_policies) — batched
	roleRecords, err := app.FindRecordsByFilter(
		"iam_user_roles",
		"user = {:uid}",
		"", 0, 0,
		dbx.Params{"uid": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user roles: %w", err)
	}
	if len(roleRecords) > 0 {
		roleIDs := make([]string, len(roleRecords))
		for i, ur := range roleRecords {
			roleIDs[i] = ur.GetString("role")
		}
		filter, params := buildInFilter("role", roleIDs)
		rolePolicyRecords, err := app.FindRecordsByFilter(
			"iam_role_policies",
			filter,
			"", 0, 0,
			params,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch role policies: %w", err)
		}
		for _, rp := range rolePolicyRecords {
			policyIDs[rp.GetString("policy")] = struct{}{}
		}
	}

	// No policies found — return empty
	if len(policyIDs) == 0 {
		return nil, nil
	}

	// Batch-fetch all unique policies using FindRecordsByIds (safe, no SQL injection)
	idSlice := make([]string, 0, len(policyIDs))
	for id := range policyIDs {
		idSlice = append(idSlice, id)
	}
	policyRecords, err := app.FindRecordsByIds("iam_policies", idSlice)
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

// buildInFilter constructs a parameterized PocketBase filter for matching
// multiple values on a field. e.g. "field = {:p0} || field = {:p1}"
func buildInFilter(field string, ids []string) (string, dbx.Params) {
	if len(ids) == 0 {
		return "1=0", dbx.Params{} // match nothing
	}
	params := dbx.Params{}
	parts := make([]string, len(ids))
	for i, id := range ids {
		key := fmt.Sprintf("p%d", i)
		parts[i] = fmt.Sprintf("%s = {:%s}", field, key)
		params[key] = id
	}
	return strings.Join(parts, " || "), params
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
