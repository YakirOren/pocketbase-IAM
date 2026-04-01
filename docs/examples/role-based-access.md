# Role-Based Access

Define an "editor" role that can read, create, and update posts and comments.

## Policy

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

## Setup

1. Register `posts` and `comments` as [managed collections](/concepts/managed-collections)
2. Create the policy above
3. Create a role called "editor" in `iam_roles`
4. Attach the policy to the "editor" role via `iam_role_policies`
5. Assign the "editor" role to users via `iam_user_roles`

Any user with the "editor" role can now read, create, and update posts and comments. Deleting is implicitly denied.
