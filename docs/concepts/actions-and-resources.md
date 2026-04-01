# Actions & Resources

Actions describe **what** can be done. Resources describe **where** it applies.

## CRUD Actions

IAM automatically enforces these actions on [managed collections](/concepts/managed-collections):

| Action | PocketBase Operation |
|--------|---------------------|
| `collections:list` | List records |
| `collections:view` | View a single record |
| `collections:create` | Create a record |
| `collections:update` | Update a record |
| `collections:delete` | Delete a record |

::: info
`collections:read` is a convenience alias that matches both `collections:list` and `collections:view`.
:::

## Resources

For CRUD operations, the resource is the collection name:

```json
{ "action": ["collections:read"], "resource": ["posts"] }
```

A single statement can target multiple collections:

```json
{ "action": ["collections:read"], "resource": ["posts", "comments"] }
```

## Custom Actions

You can register custom actions for use outside of CRUD enforcement. Register them at startup:

```go
iam.RegisterAction(app, "custom:billing:refund", "Issue a billing refund")
```

Then check them via the [check endpoint](/api/check-endpoint):

```json
{ "action": "custom:billing:refund", "resource": "order:123" }
```

Custom actions appear in the `iam_actions` view alongside CRUD actions.

## Wildcards

Use `*` to match any value:

| Pattern | Matches |
|---------|---------|
| `collections:*` | Any collection operation |
| `*` | Any action or any resource |
| `custom:billing:*` | Any custom billing action |

Wildcards work in both `action` and `resource` arrays.
