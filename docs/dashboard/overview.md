# Dashboard Overview

pocketbase-IAM includes a built-in admin dashboard at `/_/iam/`. It provides a web UI for managing all IAM resources.

## Accessing the Dashboard

Start your PocketBase server and navigate to `http://localhost:8090/_/iam/`. You must be logged in as a superuser.

## Pages

| Page | Purpose |
|------|---------|
| **Policies** | Create and edit JSON policy documents |
| **Roles** | Define roles and attach policies to them |
| **Groups** | Manage user groups and their policies |
| **Users** | View users and their IAM assignments (direct policies, roles, groups) |
| **Managed Collections** | Register or unregister collections for IAM enforcement |
| **Simulator** | Test policy evaluation for any user/action/resource combination |

## Tech Stack

The dashboard is built with React, [Refine](https://refine.dev), Tailwind CSS, and [shadcn/ui](https://ui.shadcn.com). It is embedded in the Go binary — no separate frontend deployment is needed.

## Development

To develop the dashboard with hot reload:

```bash
cd ui
npm install
npm run dev
```

The Vite dev server proxies API requests to PocketBase at `:8090`.

To build the dashboard into the Go binary:

```bash
cd ui
npm run build
```

This outputs to `iam/dashboard/`, which is embedded via `//go:embed`.
