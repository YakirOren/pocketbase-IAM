package iam

import (
	"testing"

	"github.com/pocketbase/dbx"
)

func TestEvaluateStatements(t *testing.T) {
	allowRead := Statement{
		SID:      "AllowRead",
		Effect:   "Allow",
		Action:   []string{"collections:read"},
		Resource: []string{"posts"},
	}
	allowWrite := Statement{
		SID:      "AllowWrite",
		Effect:   "Allow",
		Action:   []string{"collections:create", "collections:update"},
		Resource: []string{"posts"},
	}
	denyRead := Statement{
		SID:      "DenyRead",
		Effect:   "Deny",
		Action:   []string{"collections:read"},
		Resource: []string{"posts"},
	}
	allowAll := Statement{
		SID:      "AllowAll",
		Effect:   "Allow",
		Action:   []string{"*"},
		Resource: []string{"*"},
	}
	denyDelete := Statement{
		SID:      "DenyDelete",
		Effect:   "Deny",
		Action:   []string{"collections:delete"},
		Resource: []string{"*"},
	}
	allowSpecificResource := Statement{
		SID:      "AllowOrder",
		Effect:   "Allow",
		Action:   []string{"custom:billing:refund"},
		Resource: []string{"order:*"},
	}

	tests := []struct {
		name       string
		stmts      []Statement
		action     string
		resource   string
		wantAllow  bool
		wantReason string
	}{
		{
			name:       "allow matches",
			stmts:      []Statement{allowRead},
			action:     "collections:read",
			resource:   "posts",
			wantAllow:  true,
			wantReason: "allowed by AllowRead",
		},
		{
			name:       "allow does not match different action",
			stmts:      []Statement{allowRead},
			action:     "collections:create",
			resource:   "posts",
			wantAllow:  false,
			wantReason: "implicit deny: no matching policy",
		},
		{
			name:       "allow does not match different resource",
			stmts:      []Statement{allowRead},
			action:     "collections:read",
			resource:   "users",
			wantAllow:  false,
			wantReason: "implicit deny: no matching policy",
		},
		{
			name:       "deny matches",
			stmts:      []Statement{denyRead},
			action:     "collections:read",
			resource:   "posts",
			wantAllow:  false,
			wantReason: "explicit deny by DenyRead",
		},
		{
			name:       "deny overrides allow",
			stmts:      []Statement{allowRead, denyRead},
			action:     "collections:read",
			resource:   "posts",
			wantAllow:  false,
			wantReason: "explicit deny by DenyRead",
		},
		{
			name:       "deny overrides allow regardless of order",
			stmts:      []Statement{denyRead, allowRead},
			action:     "collections:read",
			resource:   "posts",
			wantAllow:  false,
			wantReason: "explicit deny by DenyRead",
		},
		{
			name:       "wildcard allow matches any action and resource",
			stmts:      []Statement{allowAll},
			action:     "collections:delete",
			resource:   "posts",
			wantAllow:  true,
			wantReason: "allowed by AllowAll",
		},
		{
			name:       "wildcard deny matches delete on any resource",
			stmts:      []Statement{allowAll, denyDelete},
			action:     "collections:delete",
			resource:   "posts",
			wantAllow:  false,
			wantReason: "explicit deny by DenyDelete",
		},
		{
			name:       "wildcard deny does not affect non-delete",
			stmts:      []Statement{allowAll, denyDelete},
			action:     "collections:read",
			resource:   "posts",
			wantAllow:  true,
			wantReason: "allowed by AllowAll",
		},
		{
			name:       "empty statements implicit deny",
			stmts:      nil,
			action:     "collections:read",
			resource:   "posts",
			wantAllow:  false,
			wantReason: "implicit deny: no matching policy",
		},
		{
			name:       "multiple allow statements first match wins",
			stmts:      []Statement{allowRead, allowWrite},
			action:     "collections:create",
			resource:   "posts",
			wantAllow:  true,
			wantReason: "allowed by AllowWrite",
		},
		{
			name:       "one statement grants action on multiple resources",
			stmts: []Statement{{
				SID:      "AllowReadMultiple",
				Effect:   "Allow",
				Action:   []string{"collections:read"},
				Resource: []string{"posts", "comments"},
			}},
			action:     "collections:read",
			resource:   "comments",
			wantAllow:  true,
			wantReason: "allowed by AllowReadMultiple",
		},
		{
			name:       "wildcard resource grants action on any collection",
			stmts: []Statement{{
				SID:      "AllowReadAll",
				Effect:   "Allow",
				Action:   []string{"collections:read"},
				Resource: []string{"*"},
			}},
			action:     "collections:read",
			resource:   "users",
			wantAllow:  true,
			wantReason: "allowed by AllowReadAll",
		},
		{
			name:       "resource matching",
			stmts:      []Statement{allowSpecificResource},
			action:     "custom:billing:refund",
			resource:   "order:123",
			wantAllow:  true,
			wantReason: "allowed by AllowOrder",
		},
		{
			name:       "resource mismatch",
			stmts:      []Statement{allowSpecificResource},
			action:     "custom:billing:refund",
			resource:   "invoice:123",
			wantAllow:  false,
			wantReason: "implicit deny: no matching policy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, reason := evaluateStatements(tt.stmts, tt.action, tt.resource)
			if allowed != tt.wantAllow {
				t.Errorf("allowed = %v, want %v", allowed, tt.wantAllow)
			}
			if reason != tt.wantReason {
				t.Errorf("reason = %q, want %q", reason, tt.wantReason)
			}
		})
	}
}

func TestBuildInFilter(t *testing.T) {
	t.Run("single ID", func(t *testing.T) {
		filter, params := buildInFilter("user", []string{"abc123"})
		if filter != "user = {:p0}" {
			t.Errorf("filter = %q, want %q", filter, "user = {:p0}")
		}
		if params["p0"] != "abc123" {
			t.Errorf("params[p0] = %v, want %q", params["p0"], "abc123")
		}
	})

	t.Run("multiple IDs", func(t *testing.T) {
		filter, params := buildInFilter("role", []string{"r1", "r2", "r3"})
		// Filter should contain all three conditions joined by ||
		expected := "role = {:p0} || role = {:p1} || role = {:p2}"
		if filter != expected {
			t.Errorf("filter = %q, want %q", filter, expected)
		}
		if len(params) != 3 {
			t.Fatalf("got %d params, want 3", len(params))
		}
		assertParamsContainValues(t, params, []string{"r1", "r2", "r3"})
	})

	t.Run("empty slice", func(t *testing.T) {
		filter, params := buildInFilter("user", nil)
		if filter != "1=0" {
			t.Errorf("filter = %q, want %q", filter, "1=0")
		}
		if len(params) != 0 {
			t.Errorf("got %d params, want 0", len(params))
		}
	})
}

func assertParamsContainValues(t *testing.T, params dbx.Params, values []string) {
	t.Helper()
	seen := make(map[string]bool)
	for _, v := range params {
		if s, ok := v.(string); ok {
			seen[s] = true
		}
	}
	for _, v := range values {
		if !seen[v] {
			t.Errorf("params missing value %q", v)
		}
	}
}
