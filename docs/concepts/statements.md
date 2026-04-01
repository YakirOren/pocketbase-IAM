# Statements

A statement is a single permission rule within a [policy](/concepts/policies). Each statement either allows or denies specific actions on specific resources.

## Statement Format

```json
{
  "sid": "AllowReadPosts",
  "effect": "Allow",
  "action": ["collections:read"],
  "resource": ["posts"]
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `sid` | No | Statement ID. A human-readable identifier for the statement |
| `effect` | Yes | `"Allow"` or `"Deny"` |
| `action` | Yes | Array of [action strings](/concepts/actions-and-resources) |
| `resource` | Yes | Array of [resource strings](/concepts/actions-and-resources) |

## Allow vs Deny

- **Allow** — grants permission to perform the listed actions on the listed resources
- **Deny** — explicitly blocks the actions, even if another statement allows them

::: warning
Deny always overrides Allow. If any statement from any policy (direct, role, or group) denies an action, it is denied — regardless of how many Allow statements exist.
:::

## Implicit Deny

If no statement matches a request (neither Allow nor Deny), the request is implicitly denied. You don't need Deny statements to block access — only to override existing Allows.

## Wildcard Matching

Both `action` and `resource` arrays support `*` wildcards:

```json
{ "action": ["collections:*"], "resource": ["*"] }
```

This matches any collection operation on any resource. See [Actions & Resources](/concepts/actions-and-resources) for details.
