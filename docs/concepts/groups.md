# Groups

A group is a collection of users that share policies. Groups provide a third path for attaching policies alongside [direct attachment](/concepts/policies#attaching-policies) and [roles](/concepts/roles).

## How Groups Work

1. Create a group in `iam_groups` (e.g., "engineering", "support")
2. Add users to the group via `iam_group_users`
3. Attach policies to the group via `iam_group_policies`

During [evaluation](/concepts/evaluation-flow), all policies from all of a user's groups are collected alongside direct and role policies.

## Groups vs Roles

| | Roles | Groups |
|-|-------|--------|
| **Attaches to** | Users | Users |
| **Contains** | Policies | Users + Policies |
| **Use case** | "What can this permission set do?" | "Who is in this team?" |

Both are organizational tools. Use roles when you're thinking about permissions ("editors can do X"). Use groups when you're thinking about people ("the engineering team gets Y").

## Combining Groups and Roles

A user can belong to multiple groups and have multiple roles simultaneously. All statements from all paths are collected and evaluated together. [Deny always overrides Allow](/concepts/evaluation-flow#deny-always-wins).
