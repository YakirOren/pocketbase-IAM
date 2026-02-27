package iam

import "testing"

func TestParsePolicy(t *testing.T) {
	validJSON := `{"version":"2024-01-01","statement":[{"sid":"s1","effect":"Allow","action":["collections:posts:read"],"resource":["*"]}]}`

	t.Run("string input", func(t *testing.T) {
		doc, err := ParsePolicy(validJSON)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc.Version != "2024-01-01" {
			t.Errorf("version = %q, want %q", doc.Version, "2024-01-01")
		}
		if len(doc.Statement) != 1 {
			t.Fatalf("got %d statements, want 1", len(doc.Statement))
		}
		if doc.Statement[0].SID != "s1" {
			t.Errorf("sid = %q, want %q", doc.Statement[0].SID, "s1")
		}
	})

	t.Run("byte slice input", func(t *testing.T) {
		doc, err := ParsePolicy([]byte(validJSON))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc.Version != "2024-01-01" {
			t.Errorf("version = %q, want %q", doc.Version, "2024-01-01")
		}
	})

	t.Run("map input", func(t *testing.T) {
		input := map[string]any{
			"version": "2024-01-01",
			"statement": []any{
				map[string]any{
					"sid":      "s1",
					"effect":   "Allow",
					"action":   []any{"collections:posts:read"},
					"resource": []any{"*"},
				},
			},
		}
		doc, err := ParsePolicy(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc.Version != "2024-01-01" {
			t.Errorf("version = %q, want %q", doc.Version, "2024-01-01")
		}
		if len(doc.Statement) != 1 {
			t.Fatalf("got %d statements, want 1", len(doc.Statement))
		}
	})

	t.Run("invalid JSON string", func(t *testing.T) {
		_, err := ParsePolicy("{invalid")
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})

	t.Run("nil input parses to empty doc", func(t *testing.T) {
		// nil marshals to "null" JSON, which unmarshals to zero-value PolicyDocument.
		// ValidatePolicy catches this (empty version, empty statements).
		doc, err := ParsePolicy(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := ValidatePolicy(doc); err == nil {
			t.Fatal("expected validation to fail for nil-parsed doc")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := ParsePolicy("")
		if err == nil {
			t.Fatal("expected error for empty string")
		}
	})
}

func TestValidatePolicy(t *testing.T) {
	validDoc := &PolicyDocument{
		Version: "2024-01-01",
		Statement: []Statement{
			{SID: "s1", Effect: "Allow", Action: []string{"collections:posts:read"}, Resource: []string{"*"}},
		},
	}

	t.Run("valid document", func(t *testing.T) {
		if err := ValidatePolicy(validDoc); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("valid deny effect", func(t *testing.T) {
		doc := &PolicyDocument{
			Version: "v1",
			Statement: []Statement{
				{SID: "s1", Effect: "Deny", Action: []string{"collections:posts:read"}, Resource: []string{"*"}},
			},
		}
		if err := ValidatePolicy(doc); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("wildcard action", func(t *testing.T) {
		doc := &PolicyDocument{
			Version: "v1",
			Statement: []Statement{
				{SID: "s1", Effect: "Allow", Action: []string{"*"}, Resource: []string{"*"}},
			},
		}
		if err := ValidatePolicy(doc); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty version", func(t *testing.T) {
		doc := &PolicyDocument{Version: "", Statement: []Statement{{SID: "s1", Effect: "Allow", Action: []string{"a:b"}, Resource: []string{"*"}}}}
		if err := ValidatePolicy(doc); err == nil {
			t.Fatal("expected error for empty version")
		}
	})

	t.Run("empty statements", func(t *testing.T) {
		doc := &PolicyDocument{Version: "v1", Statement: nil}
		if err := ValidatePolicy(doc); err == nil {
			t.Fatal("expected error for empty statements")
		}
	})

	t.Run("bad effect", func(t *testing.T) {
		doc := &PolicyDocument{
			Version: "v1",
			Statement: []Statement{
				{SID: "s1", Effect: "allow", Action: []string{"a:b"}, Resource: []string{"*"}},
			},
		}
		if err := ValidatePolicy(doc); err == nil {
			t.Fatal("expected error for bad effect")
		}
	})

	t.Run("empty actions", func(t *testing.T) {
		doc := &PolicyDocument{
			Version: "v1",
			Statement: []Statement{
				{SID: "s1", Effect: "Allow", Action: nil, Resource: []string{"*"}},
			},
		}
		if err := ValidatePolicy(doc); err == nil {
			t.Fatal("expected error for empty actions")
		}
	})

	t.Run("action without colon", func(t *testing.T) {
		doc := &PolicyDocument{
			Version: "v1",
			Statement: []Statement{
				{SID: "s1", Effect: "Allow", Action: []string{"invalid"}, Resource: []string{"*"}},
			},
		}
		if err := ValidatePolicy(doc); err == nil {
			t.Fatal("expected error for action without colon")
		}
	})

	t.Run("empty resources", func(t *testing.T) {
		doc := &PolicyDocument{
			Version: "v1",
			Statement: []Statement{
				{SID: "s1", Effect: "Allow", Action: []string{"a:b"}, Resource: nil},
			},
		}
		if err := ValidatePolicy(doc); err == nil {
			t.Fatal("expected error for empty resources")
		}
	})

	t.Run("empty string in resources", func(t *testing.T) {
		doc := &PolicyDocument{
			Version: "v1",
			Statement: []Statement{
				{SID: "s1", Effect: "Allow", Action: []string{"a:b"}, Resource: []string{""}},
			},
		}
		if err := ValidatePolicy(doc); err == nil {
			t.Fatal("expected error for empty string in resources")
		}
	})

	t.Run("multiple statements mixed", func(t *testing.T) {
		doc := &PolicyDocument{
			Version: "v1",
			Statement: []Statement{
				{SID: "s1", Effect: "Allow", Action: []string{"collections:posts:read"}, Resource: []string{"*"}},
				{SID: "s2", Effect: "Deny", Action: []string{"collections:users:delete"}, Resource: []string{"*"}},
			},
		}
		if err := ValidatePolicy(doc); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
