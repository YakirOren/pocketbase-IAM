# Deny Overrides

Deny statements always override Allow, regardless of which policy or attachment path they come from.

## Broad Allow + Targeted Deny

Allow full access but block deletes:

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowEverything",
      "effect": "Allow",
      "action": ["*"],
      "resource": ["*"]
    },
    {
      "sid": "BlockDeletes",
      "effect": "Deny",
      "action": ["collections:delete"],
      "resource": ["*"]
    }
  ]
}
```

Even though `AllowEverything` matches `collections:delete`, the explicit Deny wins.

## Cross-Policy Deny

Deny works across policies and attachment paths. If a user has:

- **Direct policy** → Allow `collections:*` on `posts`
- **Role policy** → Deny `collections:delete` on `*`

The user can read, create, and update posts, but cannot delete anything. The Deny from the role policy overrides the Allow from the direct policy.

::: warning
Be careful with broad Deny statements. A Deny on `*`/`*` from any policy attached through any path will block everything, and no Allow statement can override it.
:::
