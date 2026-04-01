# Custom Actions

IAM automatically enforces `collections:*` actions on managed collections. For application-specific permissions, you can register custom actions.

## Registering an Action

```go
iam.RegisterAction(app, "custom:billing:refund", "Issue a billing refund")
```

Call `RegisterAction` after `iam.Setup()` and before `app.Start()`. Registration is idempotent — calling it multiple times with the same action is safe.

## Using Custom Actions in Policies

Once registered, custom actions can be used in [policy statements](/concepts/statements):

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowRefunds",
      "effect": "Allow",
      "action": ["custom:billing:refund"],
      "resource": ["order:*"]
    }
  ]
}
```

## Checking Custom Actions

Custom actions are not automatically enforced. Use the [check endpoint](/api/check-endpoint) in your application code:

```go
app.OnServe().BindFunc(func(se *core.ServeEvent) error {
    se.Router.POST("/api/refund/{orderId}", func(e *core.RequestEvent) error {
        // Check IAM permission before proceeding
        // ... call POST /api/iam/check with action and resource
    })
    return se.Next()
})
```

## Action Registry

All registered custom actions appear in the `iam_actions` view alongside the built-in CRUD actions. This view is used by the [dashboard](/dashboard/overview) for autocomplete and discoverability.
