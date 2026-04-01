# Getting Started

## Prerequisites

- Go 1.24+
- Node 18+ (for the admin dashboard, optional)

## Install

```bash
go get github.com/yakiroren/pocketbase-IAM/iam
```

## Quick Start

```go
package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/yakiroren/pocketbase-IAM/iam"
)

func main() {
	app := pocketbase.New()

	if err := iam.Setup(app, iam.DefaultOptions()); err != nil {
		log.Fatalf("Failed to setup IAM: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
```

::: info What happens on first launch
When PocketBase starts, IAM automatically:

1. Creates 11 collections it needs (`iam_policies`, `iam_roles`, `iam_groups`, etc.)
2. Syncs rules on any already-managed collections
3. Registers enforcement hooks for all CRUD operations
:::

## Your First Policy

1. Open `http://localhost:8090/_/` and create a superuser account
2. Navigate to `/_/iam/` to open the IAM dashboard
3. Register a collection as managed (e.g., `posts`)
4. Create a policy:

```json
{
  "version": "2024-01-01",
  "statement": [
    {
      "sid": "AllowReadPosts",
      "effect": "Allow",
      "action": ["collections:read", "collections:list"],
      "resource": ["posts"]
    }
  ]
}
```

5. Attach the policy to a user, role, or group

Now authenticated users with this policy can read and list posts, but all other operations are implicitly denied.

## What's Next?

- [Policies](/concepts/policies) — policy document format and validation
- [Statements](/concepts/statements) — Allow and Deny statement semantics
- [Actions & Resources](/concepts/actions-and-resources) — what actions and resources mean
- [Evaluation Flow](/concepts/evaluation-flow) — how the engine decides Allow or Deny
- [Managed Collections](/concepts/managed-collections) — how opt-in enforcement works
