# Policy Simulator

The Simulator page lets you test policy evaluation without making real requests. It is available to superusers only.

## Using the Simulator

1. Select a **user** to simulate
2. Enter an **action** (e.g., `collections:read`)
3. Enter a **resource** (e.g., `posts`)
4. Click **Simulate**

The simulator runs the full [evaluation flow](/concepts/evaluation-flow) and shows:

- **Result** — whether the action would be allowed or denied
- **Matched statements** — which statements matched the action and resource, from which policies, and through which attachment path (direct, role, or group)

## API Equivalent

The simulator uses the [simulate endpoint](/api/simulate-endpoint) under the hood:

```bash
curl -X POST http://localhost:8090/api/iam/simulate \
  -H "Content-Type: application/json" \
  -d '{"user_id": "USER_ID", "action": "collections:read", "resource": "posts"}'
```

::: tip
Use the simulator to debug unexpected Deny results. The matched statements trace shows exactly which statement caused the denial and where it came from.
:::
