# Managed Collections

IAM enforcement is opt-in. Only collections registered in `iam_managed_collections` are gated by IAM policies. All other collections use PocketBase's native rules.

## Registering a Collection

Add a collection name to `iam_managed_collections` through the [dashboard](/dashboard/overview) or directly via the PocketBase admin UI. This is a superuser-only operation.

When a collection is registered as managed, IAM automatically sets its PocketBase rules to `@request.auth.id != ''`. This ensures:

- **Unauthenticated requests** are blocked at the PocketBase layer (before IAM)
- **Authenticated requests** pass through to IAM for policy evaluation
- **Superusers** bypass both PocketBase rules and IAM

::: warning
The rule is `@request.auth.id != ''` (not empty string `""`). Setting it to `""` would make the collection fully public, bypassing IAM entirely.
:::

## Unregistering a Collection

When a collection is removed from `iam_managed_collections`, IAM resets its PocketBase rules to `nil` (superuser-only access). You'll need to manually set appropriate rules if you want non-superuser access.

## Non-Managed Collections

Collections not in `iam_managed_collections` are completely unaffected by IAM. Their PocketBase rules work as usual.
