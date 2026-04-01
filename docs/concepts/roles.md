# Roles

A role is a named collection of policies. Instead of attaching policies directly to each user, you define a role and assign it to users.

## How Roles Work

1. Create a role in `iam_roles` (e.g., "editor", "viewer")
2. Attach policies to the role via `iam_role_policies`
3. Assign the role to users via `iam_user_roles`

During [evaluation](/concepts/evaluation-flow), all policies from all of a user's roles are collected alongside their direct policies and [group](/concepts/groups) policies.

## When to Use Roles

Roles are useful when multiple users need the same set of permissions. Instead of attaching the same policies to each user individually, create a role and assign it once.

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "EditorAccess",
      "effect": "Allow",
      "action": ["collections:read", "collections:create", "collections:update"],
      "resource": ["posts", "comments"]
    }
  ]
}
```

Attach this policy to an "editor" role, then assign the role to any user who should be able to read, create, and update posts and comments.

## Roles vs Direct Policies

There is no difference in evaluation. Statements from roles are treated identically to direct statements. Roles are an organizational tool — they make it easier to manage permissions at scale.
