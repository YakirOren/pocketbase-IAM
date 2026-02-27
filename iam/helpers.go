package iam

import "strings"

// MatchPattern checks if value matches the given pattern with wildcard support.
// A lone "*" matches any string. Otherwise, both pattern and value are split on ":"
// and compared segment by segment, where "*" in a segment matches any single segment.
func MatchPattern(pattern, value string) bool {
	if pattern == "*" {
		return true
	}

	patternParts := strings.Split(pattern, ":")
	valueParts := strings.Split(value, ":")

	if len(patternParts) != len(valueParts) {
		return false
	}

	for i, pp := range patternParts {
		if pp == "*" {
			continue
		}
		if pp != valueParts[i] {
			return false
		}
	}

	return true
}

// ActionForOperation builds the IAM action string for a PocketBase collection operation.
func ActionForOperation(collectionName, operation string) string {
	return "collections:" + collectionName + ":" + operation
}
