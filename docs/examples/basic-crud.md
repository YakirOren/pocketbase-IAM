# Basic CRUD Policy

Allow a user to read and create posts.

## Policy

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "ReadAndCreatePosts",
      "effect": "Allow",
      "action": ["collections:read", "collections:create"],
      "resource": ["posts"]
    }
  ]
}
```

## What This Grants

| Action | Resource | Result |
|--------|----------|--------|
| `collections:list` | `posts` | Allowed (`read` matches `list` and `view`) |
| `collections:view` | `posts` | Allowed |
| `collections:create` | `posts` | Allowed |
| `collections:update` | `posts` | Denied (implicit) |
| `collections:delete` | `posts` | Denied (implicit) |
| `collections:read` | `comments` | Denied (resource doesn't match) |

## Setup

1. Register `posts` as a [managed collection](/concepts/managed-collections)
2. Create the policy above in `iam_policies`
3. Attach the policy to a user via `iam_user_policies`
