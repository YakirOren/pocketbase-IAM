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
		colManagedCollections,
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

	// 1. Direct policies
	directRecords, err := app.FindRecordsByFilter(
		colUserPolicies,
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

	// 2. Group policies — batched
	groupUserRecords, err := app.FindRecordsByFilter(
		colGroupUsers,
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
			colGroupPolicies,
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

	// 3. Role policies — batched
	roleRecords, err := app.FindRecordsByFilter(
		colUserRoles,
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
			colRolePolicies,
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
	policyRecords, err := app.FindRecordsByIds(colPolicies, idSlice)
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

// TraceStep represents a single step in the evaluation trace for the simulator.
type TraceStep struct {
	Message string `json:"message"`
}

// MatchedStatement describes which statement matched during evaluation.
type MatchedStatement struct {
	SID        string `json:"sid"`
	Effect     string `json:"effect"`
	PolicyName string `json:"policy_name"`
}

// AttachedVia describes how a policy reached the user.
type AttachedVia struct {
	Type string `json:"type"` // "direct", "role", or "group"
	Name string `json:"name"`
}

// SimulateResult is the full result of a verbose policy evaluation.
type SimulateResult struct {
	Allowed          bool              `json:"allowed"`
	Reason           string            `json:"reason"`
	MatchedStatement *MatchedStatement `json:"matched_statement,omitempty"`
	Trace            []string          `json:"trace"`
}

// EvaluateVerbose performs the same evaluation as Evaluate but returns a detailed trace.
// The cache parameter is accepted for API consistency but is intentionally not used;
// the simulator always queries the database directly to ensure fresh results and
// accurate source attribution.
func EvaluateVerbose(app core.App, cache *PolicyCache, userID, action, resource string) (*SimulateResult, error) {
	result := &SimulateResult{}

	// Collect statements with source tracking
	type trackedStatement struct {
		Statement
		PolicyName string
		Source     string // "direct", "role:<name>", "group:<name>"
	}

	var tracked []trackedStatement

	// 1. Direct policies
	directRecords, err := app.FindRecordsByFilter(
		colUserPolicies,
		"user = {:uid}",
		"", 0, 0,
		dbx.Params{"uid": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user policies: %w", err)
	}
	directPolicyIDs := make([]string, 0, len(directRecords))
	for _, r := range directRecords {
		directPolicyIDs = append(directPolicyIDs, r.GetString("policy"))
	}
	result.Trace = append(result.Trace, fmt.Sprintf("Found %d direct policy attachment(s)", len(directPolicyIDs)))

	// 2. Group policies
	groupUserRecords, err := app.FindRecordsByFilter(
		colGroupUsers,
		"user = {:uid}",
		"", 0, 0,
		dbx.Params{"uid": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group memberships: %w", err)
	}
	result.Trace = append(result.Trace, fmt.Sprintf("User belongs to %d group(s)", len(groupUserRecords)))

	type groupInfo struct {
		id   string
		name string
	}
	var groups []groupInfo
	if len(groupUserRecords) > 0 {
		groupIDs := make([]string, len(groupUserRecords))
		for i, gu := range groupUserRecords {
			groupIDs[i] = gu.GetString("group")
		}
		groupRecords, err := app.FindRecordsByIds(colGroups, groupIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch groups: %w", err)
		}
		for _, gr := range groupRecords {
			groups = append(groups, groupInfo{id: gr.Id, name: gr.GetString("name")})
		}

		filter, params := buildInFilter("group", groupIDs)
		groupPolicyRecords, err := app.FindRecordsByFilter(
			colGroupPolicies, filter, "", 0, 0, params,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch group policies: %w", err)
		}
		// Collect unique policy IDs and track which groups they came from.
		groupPolicyGroups := make(map[string][]string) // policy ID → group IDs
		for _, gp := range groupPolicyRecords {
			pid := gp.GetString("policy")
			gid := gp.GetString("group")
			groupPolicyGroups[pid] = append(groupPolicyGroups[pid], gid)
		}
		gpIDs := make([]string, 0, len(groupPolicyGroups))
		for pid := range groupPolicyGroups {
			gpIDs = append(gpIDs, pid)
		}
		if len(gpIDs) > 0 {
			pRecs, err := app.FindRecordsByIds(colPolicies, gpIDs)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch group-attached policies: %w", err)
			}
			for _, pr := range pRecs {
				doc, err := ParsePolicy(pr.Get("document"))
				if err != nil {
					continue
				}
				pName := pr.GetString("name")
				for _, gid := range groupPolicyGroups[pr.Id] {
					groupName := gid
					for _, g := range groups {
						if g.id == gid {
							groupName = g.name
							break
						}
					}
					for _, stmt := range doc.Statement {
						tracked = append(tracked, trackedStatement{
							Statement:  stmt,
							PolicyName: pName,
							Source:     "group:" + groupName,
						})
					}
				}
			}
		}
	}

	// 3. Role policies
	roleRecords, err := app.FindRecordsByFilter(
		colUserRoles,
		"user = {:uid}",
		"", 0, 0,
		dbx.Params{"uid": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user roles: %w", err)
	}
	result.Trace = append(result.Trace, fmt.Sprintf("User has %d role(s)", len(roleRecords)))

	type roleInfo struct {
		id   string
		name string
	}
	var roles []roleInfo
	if len(roleRecords) > 0 {
		roleIDs := make([]string, len(roleRecords))
		for i, ur := range roleRecords {
			roleIDs[i] = ur.GetString("role")
		}
		roleRecs, err := app.FindRecordsByIds(colRoles, roleIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch roles: %w", err)
		}
		for _, rr := range roleRecs {
			roles = append(roles, roleInfo{id: rr.Id, name: rr.GetString("name")})
		}

		filter, params := buildInFilter("role", roleIDs)
		rolePolicyRecords, err := app.FindRecordsByFilter(
			colRolePolicies, filter, "", 0, 0, params,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch role policies: %w", err)
		}
		// Collect unique policy IDs and track which roles they came from.
		rolePolicyRoles := make(map[string][]string) // policy ID → role IDs
		for _, rp := range rolePolicyRecords {
			pid := rp.GetString("policy")
			rid := rp.GetString("role")
			rolePolicyRoles[pid] = append(rolePolicyRoles[pid], rid)
		}
		rpIDs := make([]string, 0, len(rolePolicyRoles))
		for pid := range rolePolicyRoles {
			rpIDs = append(rpIDs, pid)
		}
		if len(rpIDs) > 0 {
			pRecs, err := app.FindRecordsByIds(colPolicies, rpIDs)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch role-attached policies: %w", err)
			}
			for _, pr := range pRecs {
				doc, err := ParsePolicy(pr.Get("document"))
				if err != nil {
					continue
				}
				pName := pr.GetString("name")
				for _, rid := range rolePolicyRoles[pr.Id] {
					roleName := rid
					for _, r := range roles {
						if r.id == rid {
							roleName = r.name
							break
						}
					}
					for _, stmt := range doc.Statement {
						tracked = append(tracked, trackedStatement{
							Statement:  stmt,
							PolicyName: pName,
							Source:     "role:" + roleName,
						})
					}
				}
			}
		}
	}

	// 4. Direct policies — fetch and parse
	if len(directPolicyIDs) > 0 {
		pRecs, err := app.FindRecordsByIds(colPolicies, directPolicyIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to batch-fetch direct policies: %w", err)
		}
		for _, pr := range pRecs {
			doc, err := ParsePolicy(pr.Get("document"))
			if err != nil {
				continue
			}
			pName := pr.GetString("name")
			for _, stmt := range doc.Statement {
				tracked = append(tracked, trackedStatement{
					Statement:  stmt,
					PolicyName: pName,
					Source:     "direct",
				})
			}
		}
	}

	result.Trace = append(result.Trace, fmt.Sprintf("Collected %d statement(s) total", len(tracked)))

	// Evaluate: deny first
	for _, ts := range tracked {
		if ts.Effect != "Deny" {
			continue
		}
		for _, a := range ts.Action {
			if !MatchPattern(a, action) {
				continue
			}
			for _, r := range ts.Resource {
				if MatchPattern(r, resource) {
					result.Allowed = false
					result.Reason = fmt.Sprintf("explicit deny by %s", ts.SID)
					result.MatchedStatement = &MatchedStatement{
						SID:        ts.SID,
						Effect:     "Deny",
						PolicyName: ts.PolicyName,
					}
					result.Trace = append(result.Trace, fmt.Sprintf("DENY matched: statement %q in policy %q (via %s)", ts.SID, ts.PolicyName, ts.Source))
					return result, nil
				}
			}
		}
	}

	// Evaluate: allow
	for _, ts := range tracked {
		if ts.Effect != "Allow" {
			continue
		}
		for _, a := range ts.Action {
			if !MatchPattern(a, action) {
				continue
			}
			for _, r := range ts.Resource {
				if MatchPattern(r, resource) {
					result.Allowed = true
					result.Reason = fmt.Sprintf("allowed by %s", ts.SID)
					result.MatchedStatement = &MatchedStatement{
						SID:        ts.SID,
						Effect:     "Allow",
						PolicyName: ts.PolicyName,
					}
					result.Trace = append(result.Trace, fmt.Sprintf("ALLOW matched: statement %q in policy %q (via %s)", ts.SID, ts.PolicyName, ts.Source))
					return result, nil
				}
			}
		}
	}

	// Implicit deny
	result.Allowed = false
	result.Reason = "implicit deny: no matching policy"
	result.Trace = append(result.Trace, "No matching statement found — implicit deny")
	return result, nil
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
