package iam

import "testing"

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		value   string
		want    bool
	}{
		// Lone wildcard matches everything
		{"lone wildcard matches simple", "*", "anything", true},
		{"lone wildcard matches colon-delimited", "*", "collections:posts:read", true},
		{"lone wildcard matches empty", "*", "", true},

		// Exact matches
		{"exact match", "collections:posts:read", "collections:posts:read", true},
		{"exact mismatch operation", "collections:posts:read", "collections:posts:create", false},
		{"exact mismatch collection", "collections:posts:read", "collections:users:read", false},

		// Wildcard segment matching
		{"wildcard middle segment", "collections:*:read", "collections:posts:read", true},
		{"wildcard middle segment other", "collections:*:read", "collections:users:read", true},
		{"wildcard middle segment mismatch op", "collections:*:read", "collections:posts:create", false},
		{"wildcard last segment", "collections:posts:*", "collections:posts:read", true},
		{"wildcard last segment delete", "collections:posts:*", "collections:posts:delete", true},
		{"wildcard last segment mismatch coll", "collections:posts:*", "collections:users:read", false},
		{"wildcard first segment", "*:posts:read", "collections:posts:read", true},
		{"wildcard first segment mismatch", "*:posts:read", "collections:users:read", false},

		// Two-segment patterns
		{"two segment wildcard", "order:*", "order:123", true},
		{"two segment wildcard other", "order:*", "order:abc", true},
		{"two segment mismatch", "order:*", "invoice:123", false},

		// Segment count mismatch
		{"fewer pattern segments", "a:b", "a:b:c", false},
		{"more pattern segments", "a:b:c", "a:b", false},

		// No wildcards
		{"single segment exact", "read", "read", true},
		{"single segment mismatch", "read", "write", false},

		// Empty strings
		{"empty pattern non-wildcard", "", "", true},
		{"empty pattern vs non-empty", "", "something", false},
		{"non-empty vs empty", "something", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchPattern(tt.pattern, tt.value)
			if got != tt.want {
				t.Errorf("MatchPattern(%q, %q) = %v, want %v", tt.pattern, tt.value, got, tt.want)
			}
		})
	}
}

func TestActionForOperation(t *testing.T) {
	tests := []struct {
		collection string
		operation  string
		want       string
	}{
		{"posts", "read", "collections:posts:read"},
		{"posts", "create", "collections:posts:create"},
		{"posts", "update", "collections:posts:update"},
		{"posts", "delete", "collections:posts:delete"},
		{"posts", "list", "collections:posts:list"},
		{"users", "read", "collections:users:read"},
	}

	for _, tt := range tests {
		t.Run(tt.collection+":"+tt.operation, func(t *testing.T) {
			got := ActionForOperation(tt.collection, tt.operation)
			if got != tt.want {
				t.Errorf("ActionForOperation(%q, %q) = %q, want %q", tt.collection, tt.operation, got, tt.want)
			}
		})
	}
}
