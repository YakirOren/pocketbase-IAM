# Simulate Endpoint

`POST /api/iam/simulate` runs a verbose policy evaluation for any user. Returns the full evaluation trace including matched statements.

## Request

**Auth required:** Superuser only

```bash
curl -X POST http://localhost:8090/api/iam/simulate \
  -H "Authorization: Bearer SUPERUSER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id": "USER_ID", "action": "collections:read", "resource": "posts"}'
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_id` | `string` | Yes | The user ID to simulate |
| `action` | `string` | Yes | The action to check |
| `resource` | `string` | Yes | The resource to check |

## Response

The response includes the evaluation result and a trace of all matched statements:

```json
{
  "allowed": false,
  "reason": "explicit_deny",
  "matched_statements": [
    {
      "sid": "DenyDeleteAny",
      "effect": "Deny",
      "action": ["collections:delete"],
      "resource": ["*"],
      "policy_id": "abc123",
      "policy_name": "RestrictDeletes",
      "source": "role",
      "source_name": "editor"
    }
  ]
}
```

## Use Cases

- Debug why a user is being denied access
- Verify a policy change has the intended effect before deploying
- Audit which statements are granting or denying access

The [dashboard simulator](/dashboard/policy-simulator) provides a UI for this endpoint.
