# Group Policies

Give the "support" team read access to all managed collections.

## Policy

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "SupportReadAll",
      "effect": "Allow",
      "action": ["collections:read"],
      "resource": ["*"]
    }
  ]
}
```

## Setup

1. Create a group called "support" in `iam_groups`
2. Create the policy above
3. Attach the policy to the "support" group via `iam_group_policies`
4. Add team members to the group via `iam_group_users`

All users in the "support" group can now list and view records in every managed collection. They cannot create, update, or delete.
