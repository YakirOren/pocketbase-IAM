package iam

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PolicyDocument represents an IAM policy document containing permission statements.
type PolicyDocument struct {
	Version   string      `json:"version"`
	Statement []Statement `json:"statement"`
}

// Statement represents a single permission statement within a policy document.
type Statement struct {
	SID      string   `json:"sid"`
	Effect   string   `json:"effect"`
	Action   []string `json:"action"`
	Resource []string `json:"resource"`
}

// ParsePolicy parses a raw value (from record.Get("document")) into a PolicyDocument.
func ParsePolicy(raw any) (*PolicyDocument, error) {
	var data []byte

	switch v := raw.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		var err error
		data, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal policy document: %w", err)
		}
	}

	var doc PolicyDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse policy document: %w", err)
	}

	return &doc, nil
}

// ValidatePolicy validates a parsed policy document.
func ValidatePolicy(doc *PolicyDocument) error {
	if doc.Version == "" {
		return fmt.Errorf("version must be a non-empty string")
	}

	if len(doc.Statement) == 0 {
		return fmt.Errorf("statement must be a non-empty array")
	}

	for i, s := range doc.Statement {
		if s.Effect != "Allow" && s.Effect != "Deny" {
			return fmt.Errorf("statement[%d]: effect must be Allow or Deny, got %q", i, s.Effect)
		}

		if len(s.Action) == 0 {
			return fmt.Errorf("statement[%d]: action must be a non-empty array", i)
		}
		for j, a := range s.Action {
			if a != "*" && !strings.Contains(a, ":") {
				return fmt.Errorf("statement[%d]: action[%d] must contain ':' or be '*', got %q", i, j, a)
			}
		}

		if len(s.Resource) == 0 {
			return fmt.Errorf("statement[%d]: resource must be a non-empty array", i)
		}
	}

	return nil
}
