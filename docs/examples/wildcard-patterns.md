# Wildcard Patterns

Use `*` to match any action or resource.

## Full Admin Access

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "FullAccess",
      "effect": "Allow",
      "action": ["*"],
      "resource": ["*"]
    }
  ]
}
```

This allows all operations on all managed collections.

## All Operations on One Collection

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllPostsOperations",
      "effect": "Allow",
      "action": ["collections:*"],
      "resource": ["posts"]
    }
  ]
}
```

## Read All Collections

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "ReadEverything",
      "effect": "Allow",
      "action": ["collections:read"],
      "resource": ["*"]
    }
  ]
}
```

## Custom Action Wildcards

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllBillingActions",
      "effect": "Allow",
      "action": ["custom:billing:*"],
      "resource": ["*"]
    }
  ]
}
```
