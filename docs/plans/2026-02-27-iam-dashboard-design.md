# IAM Dashboard — Design Document

## Overview

Admin-only dashboard for managing the PocketBase IAM system. Full CRUD management of all IAM entities (policies, roles, groups, user assignments, managed collections) plus a Policy Simulator for testing permissions.

**Audience:** App admins / superusers only.
**Tech stack:** Undecided — this document defines features, UI components, and UX flow independent of framework choice.

---

## Global Layout

Fixed sidebar + scrollable main content area with breadcrumb navigation.

```
┌──────────┬──────────────────────────────────────────┐
│ SIDEBAR  │  MAIN CONTENT                            │
│          │                                          │
│ IAM      │  [Breadcrumb: Policies > ReadOnlyPosts]  │
│          │                                          │
│ Policies │  [Page content]                          │
│ Roles    │                                          │
│ Groups   │                                          │
│ Users    │                                          │
│ Managed  │                                          │
│  Colls   │                                          │
│          │                                          │
│ ──────── │                                          │
│ Simulator│                                          │
└──────────┴──────────────────────────────────────────┘
```

**Sidebar items:**
1. Policies — CRUD for `iam_policies`
2. Roles — CRUD for `iam_roles` + policy attachment
3. Groups — CRUD for `iam_groups` + member + policy management
4. Users — PB users with IAM summary view
5. Managed Collections — register/unregister collections for IAM enforcement
6. Policy Simulator — test permissions (below divider)

---

## Pages

### 1. Policies

#### List View

| Column | Source | Notes |
|--------|--------|-------|
| Name | `iam_policies.name` | Sortable, searchable |
| Description | `iam_policies.description` | Truncated |
| Statements | Computed from `document.statement.length` | Count |
| Attached To | Computed from 3 join tables | Summary badges: "3 users, 1 role, 2 groups" |
| Actions | — | Edit, Delete buttons |

**Features:**
- Search by name or action string
- Filter by effect (Allow/Deny), attached/unattached
- Paginated table
- [+ Create Policy] button top-right

#### Create/Edit View

**Top section:** Name (text input), Description (text input).

**Policy Document section — Dual-mode editor:**

**Form Builder mode (default):**
- Version field (text, defaults to "2024-01-01")
- Statement cards, each containing:
  - SID (text input)
  - Effect (radio: Allow / Deny)
  - Actions (tag list with [+ Add action] — each tag is an editable text input)
  - Resources (tag list with [+ Add resource])
  - [Remove statement] button
- [+ Add Statement] button

**Raw JSON mode (toggle):**
- Full JSON editor with syntax highlighting
- Validation on blur — inline errors matching backend rules:
  - `version` must be non-empty
  - `statement` must be non-empty array
  - Each statement: `effect` is "Allow" or "Deny", non-empty `action[]`, non-empty `resource[]`
  - Action strings must contain `:` or be lone `*`

**Attachments section (on edit only):**
- Direct Users: chip list with [+ Attach User] picker
- Roles: chip list with [+ Attach to Role] picker
- Groups: chip list with [+ Attach to Group] picker
- Each chip has [x] to detach

---

### 2. Roles

#### List View

| Column | Source |
|--------|--------|
| Name | `iam_roles.name` |
| Description | `iam_roles.description` |
| Policies | Count from `iam_role_policies` |
| Users | Count from `iam_user_roles` |
| Actions | Edit, Delete |

#### Detail/Edit View

- Name + Description (editable)
- **Attached Policies:** Chip list with [+ Attach Policy] picker, [x] to detach
- **Assigned Users:** Chip list with [+ Assign User] picker, [x] to unassign
- **Effective Permissions (read-only):** Table showing merged Allow/Deny per action from all attached policies. Computed client-side by evaluating all statements from attached policies.

| Action | Effect |
|--------|--------|
| collections:posts:read | Allow |
| collections:posts:delete | Deny |

---

### 3. Groups

Same pattern as Roles with two differences:

- **Members** section (instead of "Assigned Users") — manages `iam_group_users`
- **Attached Policies** — manages `iam_group_policies`

Includes the same **Effective Permissions** read-only table.

---

### 4. Users

#### List View

Users are read from PB's `users` collection. No create/edit/delete of users themselves — that's PB's responsibility.

| Column | Source |
|--------|--------|
| Email | `users.email` |
| Roles | Count from `iam_user_roles` |
| Groups | Count from `iam_group_users` |
| Direct Policies | Count from `iam_user_policies` |
| Actions | View |

Searchable by email or name.

#### User IAM Summary View

The most important view for auditing — answers "what can this user do?"

**Sections:**

1. **Roles** — chip list with [+ Assign Role], [x] to remove. Writes to `iam_user_roles`.
2. **Groups** — chip list with [+ Add to Group], [x] to remove. Writes to `iam_group_users`.
3. **Direct Policies** — chip list with [+ Attach Policy], [x] to detach. Writes to `iam_user_policies`.
4. **All Effective Permissions** — read-only table with **source attribution**:

| Action | Effect | Source |
|--------|--------|--------|
| collections:posts:read | Allow | Role: Editor |
| collections:posts:* | Allow | Direct: SpecialAccess |
| collections:users:delete | Deny | Group: QA Team |

Deny-overrides-allow logic applied. Source shows which policy and how it's attached (direct, via which role, via which group).

5. **[Test in Simulator]** — link to open Policy Simulator pre-filled with this user.

---

### 5. Managed Collections

Two-section layout:

**Currently Managed (top):**

| Collection | Status | PB Rule Override | Actions |
|------------|--------|------------------|---------|
| posts | Enforcing | `@request.auth.id != ''` | [Unregister] |

**Available Collections (bottom):**
PB collections NOT yet registered for IAM enforcement.

| Collection | Current PB Rules | Actions |
|------------|------------------|---------|
| products | `""` (public) | [Register] |
| orders | `@request.auth.id = user` | [Register] |

**Confirmation dialogs:**
- **Register:** "This will set PB rules to `@request.auth.id != ''` for all CRUD operations. Unauthenticated access will be blocked. IAM will gate all authenticated access. Continue?"
- **Unregister:** "This will remove IAM enforcement and reset PB rules to open (nil). You may need to manually set PB rules afterwards. Continue?"

---

### 6. Policy Simulator

Test whether a user is allowed to perform an action — like AWS IAM Policy Simulator.

**Inputs:**
- **User** — dropdown/search to select any PB user
- **Action** — text input with autocomplete from actions found in existing policy documents
- **Resource** — text input, defaults to `*`

**[Simulate] button**

**Result panel:**

```
┌────────────────────────────────────────────────────┐
│  ✓ ALLOWED  /  ✗ DENIED                           │
│                                                    │
│  Matched statement: "AllowReadPosts"               │
│  From policy: ReadOnlyPosts                        │
│  Attached via: Role "Editor"                       │
│                                                    │
│  ── Evaluation Trace ────────────────────────────  │
│  1. Collected 5 policies (2 direct, 1 role,        │
│     2 group)                                       │
│  2. Checked 8 statements                           │
│  3. No Deny match found                            │
│  4. Allow match: "AllowReadPosts" in ReadOnlyPosts │
│     via Role "Editor"                              │
└────────────────────────────────────────────────────┘
```

**Result details:**
- **ALLOWED** (green) or **DENIED** (red) status
- Which statement matched (SID + policy name)
- How that policy reaches the user (direct, via which role, via which group)
- **Evaluation trace** — step-by-step:
  1. How many policies collected and from where
  2. How many statements checked
  3. Whether any Deny matched (explicit deny)
  4. Whether any Allow matched, or implicit deny (no match)

**Implementation note:** The simulator needs a server-side endpoint that returns the full evaluation trace — unlike `/api/iam/check` which only returns `{allowed: bool}`. This could be a separate admin-only endpoint (e.g., `POST /api/iam/simulate`) or the existing check endpoint with an admin-only `verbose=true` parameter.

**Action autocomplete:** Scans all `iam_policies` documents to extract unique action strings. Provides suggestions as the admin types.

---

## Shared UI Components

| Component | Used In | Description |
|-----------|---------|-------------|
| DataTable | All list views | Sortable columns, search, filters, pagination |
| EntityChipList | Policy edit, Role detail, Group detail, User summary | Chips with [x] remove + [+ Add] picker |
| EntityPicker | Chip lists | Modal/dropdown to search and select entities (users, roles, groups, policies) |
| EffectivePermissionsTable | Role detail, Group detail, User summary | Read-only merged Allow/Deny table with source |
| PolicyFormBuilder | Policy create/edit | Statement card editor (SID, Effect, Actions, Resources) |
| JSONEditor | Policy create/edit (raw mode) | Syntax-highlighted JSON editor with validation |
| ConfirmDialog | Managed collections, all deletes | Warning dialog with description of consequences |
| Breadcrumb | Global layout | Navigation trail |
| Sidebar | Global layout | Fixed navigation menu |

---

## UX Flows

### Flow 1: Create a policy and attach it to a role

1. Sidebar > **Policies** > [+ Create Policy]
2. Fill name, description, add statements via form builder
3. Save — redirects to policy detail view
4. In Attachments section, click [+ Attach to Role]
5. Pick role from Entity Picker modal > confirm

### Flow 2: See what a user can do

1. Sidebar > **Users** > click user
2. User IAM Summary shows all roles, groups, direct policies
3. Effective Permissions table shows merged Allow/Deny with source
4. Click [Test in Simulator] to try specific actions

### Flow 3: Register a collection for IAM enforcement

1. Sidebar > **Managed Collections**
2. In "Available Collections" section, click [Register] on target collection
3. Confirm in warning dialog
4. Collection moves to "Currently Managed" section
5. PB rules auto-set to `@request.auth.id != ''`

### Flow 4: Debug why a user can't access something

1. Sidebar > **Policy Simulator**
2. Select user, type the action (autocomplete helps), set resource
3. Click [Simulate]
4. Result shows DENIED with trace: "No Allow match found" or "Explicit Deny from policy X attached via Group Y"
5. Navigate to the relevant policy/group/role to fix

### Flow 5: Create a role with policies and assign users

1. Sidebar > **Roles** > [+ Create Role]
2. Fill name, description > Save
3. On role detail, click [+ Attach Policy] — select policies
4. Click [+ Assign User] — select users
5. Effective Permissions table updates to show what this role grants

---

## Data Dependencies & API Calls

Every page's data maps directly to PocketBase's standard collection CRUD API (`/api/collections/{name}/records`):

| Page | Collections Read | Collections Write |
|------|-----------------|-------------------|
| Policies list | `iam_policies`, `iam_user_policies`, `iam_role_policies`, `iam_group_policies` | — |
| Policy edit | `iam_policies` | `iam_policies`, `iam_user_policies`, `iam_role_policies`, `iam_group_policies` |
| Roles list | `iam_roles`, `iam_role_policies`, `iam_user_roles` | — |
| Role detail | `iam_roles`, `iam_role_policies`, `iam_user_roles`, `iam_policies` | `iam_roles`, `iam_role_policies`, `iam_user_roles` |
| Groups list | `iam_groups`, `iam_group_policies`, `iam_group_users` | — |
| Group detail | `iam_groups`, `iam_group_policies`, `iam_group_users`, `iam_policies` | `iam_groups`, `iam_group_policies`, `iam_group_users` |
| Users list | `users`, `iam_user_roles`, `iam_group_users`, `iam_user_policies` | — |
| User summary | `users`, `iam_user_roles`, `iam_user_policies`, `iam_group_users`, `iam_roles`, `iam_groups`, `iam_policies`, `iam_role_policies`, `iam_group_policies` | `iam_user_roles`, `iam_user_policies`, `iam_group_users` |
| Managed Colls | `iam_managed_collections`, all PB collections (via API) | `iam_managed_collections` |
| Simulator | `users`, all policy-related collections | — (read-only, calls `/api/iam/simulate`) |

## New Backend Requirement

The Policy Simulator requires an admin-only endpoint that returns evaluation trace details:

```
POST /api/iam/simulate  (superuser-only)
Body: { "user_id": "...", "action": "...", "resource": "..." }
Response: {
  "allowed": true,
  "matched_statement": { "sid": "...", "effect": "Allow", "policy_name": "..." },
  "attached_via": { "type": "role", "name": "Editor" },
  "trace": [
    "Collected 5 policies (2 direct, 1 role, 2 group)",
    "Checked 8 statements",
    "No Deny match found",
    "Allow match: AllowReadPosts in ReadOnlyPosts via Role Editor"
  ]
}
```
