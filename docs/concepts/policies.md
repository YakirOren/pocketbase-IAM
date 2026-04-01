# Policies

A policy is a JSON document that defines permissions. Policies are stored in the `iam_policies` collection and attached to users, roles, or groups.

## Policy Format

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowReadPosts",
      "effect": "Allow",
      "action": ["collections:read", "collections:list"],
      "resource": ["posts"]
    }
  ]
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `version` | Yes | Policy version string. Must be `"2024-01-01"` |
| `statement` | Yes | Array of [statements](/concepts/statements) |

## Attaching Policies

Policies can be attached through three paths:

1. **Direct** — link a policy to a user via `iam_user_policies`
2. **Via role** — link a policy to a role via `iam_role_policies`, then assign the role to a user via `iam_user_roles`
3. **Via group** — link a policy to a group via `iam_group_policies`, then add a user to the group via `iam_group_users`

During [evaluation](/concepts/evaluation-flow), statements from all three paths are collected and evaluated together.

## Validation

Policies are validated when created or updated. The `version` field must be `"2024-01-01"`, and each statement must have a valid `effect`, at least one `action`, and at least one `resource`.

::: warning
An empty `statement` array is valid — it simply grants no permissions. The user will be implicitly denied everything.
:::
