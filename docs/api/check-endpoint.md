# Check Endpoint

`POST /api/iam/check` evaluates whether the authenticated user can perform an action on a resource.

## Request

**Auth required:** Yes (any authenticated user)

```bash
curl -X POST http://localhost:8090/api/iam/check \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action": "collections:read", "resource": "posts"}'
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | `string` | Yes | The action to check (e.g., `collections:read`) |
| `resource` | `string` | Yes | The resource to check (e.g., `posts`) |

## Response

```json
{ "allowed": true }
```

| Field | Type | Description |
|-------|------|-------------|
| `allowed` | `boolean` | Whether the action is permitted |

## Behavior

- Superusers always get `{"allowed": true}`
- Unauthenticated requests get 401
- The endpoint runs the full [evaluation flow](/concepts/evaluation-flow)

::: tip
Use this endpoint for [custom actions](/api/custom-actions) that aren't automatically enforced by CRUD hooks. For example, checking if a user can issue a refund before processing it in your application code.
:::
