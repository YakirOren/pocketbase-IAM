# Managing Policies

The Policies page lets you create, edit, and delete [policy documents](/concepts/policies).

## Creating a Policy

1. Click **Create** on the Policies page
2. Enter a name for the policy
3. Write the policy JSON in the editor:

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowReadPosts",
      "effect": "Allow",
      "action": ["collections:read"],
      "resource": ["posts"]
    }
  ]
}
```

4. Save the policy

The editor validates the JSON on save. Invalid policies (wrong version, missing fields, bad effect values) are rejected with an error message.

## Attaching a Policy

After creating a policy, attach it to users, roles, or groups:

- **To a user** — go to the Users page, select a user, and add the policy
- **To a role** — go to the Roles page, select a role, and add the policy
- **To a group** — go to the Groups page, select a group, and add the policy

See [Policies — Attaching Policies](/concepts/policies#attaching-policies) for how the three attachment paths work.
