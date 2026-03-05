# IAM Dashboard — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a React admin dashboard for managing PocketBase IAM entities (policies, roles, groups, users, managed collections) with a policy simulator.

**Architecture:** Refine (headless CRUD framework) + shadcn/ui + React Router v7, using the `refine-pocketbase` community provider for data/auth/live. The frontend lives in `ui/` and builds to `pb_public/` for PocketBase to serve.

**Tech Stack:** Refine v4, React 19, TypeScript, Vite, Tailwind CSS v4, shadcn/ui, TanStack Table, CodeMirror 6, Lucide React, refine-pocketbase

**Design doc:** `docs/plans/2026-02-27-iam-dashboard-design.md`

---

## Task 1: Scaffold Vite + React + TypeScript project

**Files:**
- Create: `ui/package.json`, `ui/tsconfig.json`, `ui/tsconfig.app.json`, `ui/tsconfig.node.json`, `ui/vite.config.ts`, `ui/index.html`, `ui/src/main.tsx`, `ui/src/index.css`

**Step 1: Create the Vite project**

```bash
cd /Users/yakiroren/Documents/Projects/pocketbase-IAM
npm create vite@latest ui -- --template react-ts
```

**Step 2: Install dependencies**

```bash
cd ui && npm install
```

**Step 3: Verify it runs**

```bash
npm run dev
```

Expected: Vite dev server starts on http://localhost:5173 with default React template.

**Step 4: Commit**

```bash
git add ui/
git commit -m "scaffold Vite React TypeScript project in ui/"
```

---

## Task 2: Configure Tailwind CSS v4 + shadcn/ui

**Files:**
- Modify: `ui/vite.config.ts`, `ui/tsconfig.json`, `ui/tsconfig.app.json`, `ui/src/index.css`
- Create: `ui/components.json`, `ui/src/lib/utils.ts`

**Step 1: Install Tailwind CSS v4**

```bash
cd ui
npm install tailwindcss @tailwindcss/vite
```

**Step 2: Configure vite.config.ts**

```ts
// ui/vite.config.ts
import path from "path";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://127.0.0.1:8090",
        changeOrigin: true,
      },
      "/_": {
        target: "http://127.0.0.1:8090",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "../pb_public",
    emptyOutDir: true,
  },
});
```

**Step 3: Add path aliases to tsconfig files**

In `ui/tsconfig.json`:
```json
{
  "files": [],
  "references": [
    { "path": "./tsconfig.app.json" },
    { "path": "./tsconfig.node.json" }
  ],
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

In `ui/tsconfig.app.json`, add to `compilerOptions`:
```json
"baseUrl": ".",
"paths": {
  "@/*": ["./src/*"]
}
```

**Step 4: Replace src/index.css contents**

```css
@import "tailwindcss";
```

**Step 5: Install @types/node**

```bash
npm install -D @types/node
```

**Step 6: Initialize shadcn/ui**

```bash
npx shadcn@latest init
```

Choose: New York style, Neutral base color, CSS variables yes, `src/index.css` for CSS.

**Step 7: Add core shadcn components**

```bash
npx shadcn@latest add button input label table card badge dialog separator tabs select textarea dropdown-menu form sheet tooltip
```

**Step 8: Verify Tailwind works**

Replace `ui/src/App.tsx` with:
```tsx
export default function App() {
  return <h1 className="text-2xl font-bold p-8">IAM Dashboard</h1>;
}
```

Run `npm run dev` — should show styled heading.

**Step 9: Commit**

```bash
git add ui/
git commit -m "configure Tailwind CSS v4 and shadcn/ui"
```

---

## Task 3: Install Refine + refine-pocketbase

**Files:**
- Modify: `ui/package.json`

**Step 1: Install Refine packages**

```bash
cd ui
npm install @refinedev/core @refinedev/react-router @refinedev/react-table @refinedev/react-hook-form react-router refine-pocketbase pocketbase
```

**Step 2: Verify install**

```bash
npm ls @refinedev/core refine-pocketbase pocketbase
```

Expected: All three packages listed with versions.

**Step 3: Commit**

```bash
git add ui/package.json ui/package-lock.json
git commit -m "add Refine and refine-pocketbase dependencies"
```

---

## Task 4: Set up PocketBase client and Refine providers

**Files:**
- Create: `ui/src/providers/pocketbase.ts`

**Step 1: Create the provider config**

```ts
// ui/src/providers/pocketbase.ts
import PocketBase from "pocketbase";
import { dataProvider, authProvider, liveProvider } from "refine-pocketbase";

// In dev, Vite proxies /api to PB. In prod, same origin.
const pb = new PocketBase("/");

export const pbClient = pb;
export const pbDataProvider = dataProvider(pb);
export const pbAuthProvider = authProvider(pb, {
  collection: "superusers",
});
export const pbLiveProvider = liveProvider(pb);
```

**Step 2: Commit**

```bash
git add ui/src/providers/
git commit -m "add PocketBase client and Refine providers"
```

---

## Task 5: Build the sidebar layout

**Files:**
- Create: `ui/src/components/layout/sidebar.tsx`, `ui/src/components/layout/layout.tsx`, `ui/src/components/layout/index.ts`

**Step 1: Create the sidebar component**

```tsx
// ui/src/components/layout/sidebar.tsx
import { Link, useLocation } from "react-router";
import {
  Shield,
  UserCheck,
  Users,
  Contact,
  Database,
  FlaskConical,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { Separator } from "@/components/ui/separator";

const navItems = [
  { label: "Policies", path: "/policies", icon: Shield },
  { label: "Roles", path: "/roles", icon: UserCheck },
  { label: "Groups", path: "/groups", icon: Users },
  { label: "Users", path: "/users", icon: Contact },
  { label: "Managed Collections", path: "/managed-collections", icon: Database },
];

const toolItems = [
  { label: "Policy Simulator", path: "/simulator", icon: FlaskConical },
];

export function Sidebar() {
  const location = useLocation();

  const renderItem = (item: (typeof navItems)[0]) => {
    const isActive = location.pathname.startsWith(item.path);
    return (
      <Link
        key={item.path}
        to={item.path}
        className={cn(
          "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
          isActive
            ? "bg-accent text-accent-foreground"
            : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
        )}
      >
        <item.icon className="h-4 w-4" />
        {item.label}
      </Link>
    );
  };

  return (
    <aside className="flex h-screen w-60 flex-col border-r bg-background">
      <div className="flex h-14 items-center border-b px-4">
        <Shield className="mr-2 h-5 w-5" />
        <span className="text-lg font-semibold">IAM</span>
      </div>
      <nav className="flex-1 space-y-1 p-3">
        {navItems.map(renderItem)}
        <Separator className="my-3" />
        {toolItems.map(renderItem)}
      </nav>
    </aside>
  );
}
```

**Step 2: Create the layout wrapper**

```tsx
// ui/src/components/layout/layout.tsx
import { Sidebar } from "./sidebar";
import { Outlet } from "react-router";

export function Layout() {
  return (
    <div className="flex h-screen">
      <Sidebar />
      <main className="flex-1 overflow-y-auto p-6">
        <Outlet />
      </main>
    </div>
  );
}
```

**Step 3: Create barrel export**

```ts
// ui/src/components/layout/index.ts
export { Layout } from "./layout";
export { Sidebar } from "./sidebar";
```

**Step 4: Commit**

```bash
git add ui/src/components/layout/
git commit -m "add sidebar and layout components"
```

---

## Task 6: Wire up App.tsx with Refine routing

**Files:**
- Modify: `ui/src/App.tsx`
- Create: `ui/src/pages/login.tsx`

**Step 1: Create a login page**

```tsx
// ui/src/pages/login.tsx
import { useState } from "react";
import { useLogin } from "@refinedev/core";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Shield } from "lucide-react";

export function LoginPage() {
  const { mutate: login, isLoading } = useLogin();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    login({ email, password });
  };

  return (
    <div className="flex h-screen items-center justify-center bg-muted">
      <Card className="w-full max-w-sm">
        <CardHeader className="text-center">
          <Shield className="mx-auto mb-2 h-8 w-8" />
          <CardTitle>IAM Dashboard</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Signing in..." : "Sign in"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
```

**Step 2: Wire up App.tsx**

```tsx
// ui/src/App.tsx
import { Refine, Authenticated } from "@refinedev/core";
import routerProvider, {
  NavigateToResource,
  CatchAllNavigate,
  UnsavedChangesNotifier,
  DocumentTitleHandler,
} from "@refinedev/react-router";
import { BrowserRouter, Routes, Route, Outlet } from "react-router";
import {
  Shield,
  UserCheck,
  Users,
  Contact,
  Database,
  FlaskConical,
} from "lucide-react";

import { pbDataProvider, pbAuthProvider, pbLiveProvider } from "@/providers/pocketbase";
import { Layout } from "@/components/layout";
import { LoginPage } from "@/pages/login";

function App() {
  return (
    <BrowserRouter>
      <Refine
        routerProvider={routerProvider}
        dataProvider={pbDataProvider}
        authProvider={pbAuthProvider}
        liveProvider={pbLiveProvider}
        resources={[
          {
            name: "iam_policies",
            list: "/policies",
            create: "/policies/create",
            edit: "/policies/edit/:id",
            show: "/policies/show/:id",
            meta: { canDelete: true, label: "Policies", icon: <Shield className="h-4 w-4" /> },
          },
          {
            name: "iam_roles",
            list: "/roles",
            create: "/roles/create",
            edit: "/roles/edit/:id",
            show: "/roles/show/:id",
            meta: { canDelete: true, label: "Roles", icon: <UserCheck className="h-4 w-4" /> },
          },
          {
            name: "iam_groups",
            list: "/groups",
            create: "/groups/create",
            edit: "/groups/edit/:id",
            show: "/groups/show/:id",
            meta: { canDelete: true, label: "Groups", icon: <Users className="h-4 w-4" /> },
          },
          {
            name: "users",
            list: "/users",
            show: "/users/show/:id",
            meta: { label: "Users", icon: <Contact className="h-4 w-4" /> },
          },
          {
            name: "iam_managed_collections",
            list: "/managed-collections",
            meta: { label: "Managed Collections", icon: <Database className="h-4 w-4" /> },
          },
        ]}
        options={{
          syncWithLocation: true,
          warnWhenUnsavedChanges: true,
          liveMode: "auto",
        }}
      >
        <Routes>
          {/* Authenticated routes */}
          <Route
            element={
              <Authenticated
                key="authenticated-routes"
                fallback={<CatchAllNavigate to="/login" />}
              >
                <Layout />
              </Authenticated>
            }
          >
            <Route index element={<NavigateToResource resource="iam_policies" />} />

            <Route path="/policies">
              <Route index element={<div>Policies list (TODO)</div>} />
              <Route path="create" element={<div>Create policy (TODO)</div>} />
              <Route path="edit/:id" element={<div>Edit policy (TODO)</div>} />
              <Route path="show/:id" element={<div>Show policy (TODO)</div>} />
            </Route>

            <Route path="/roles">
              <Route index element={<div>Roles list (TODO)</div>} />
              <Route path="create" element={<div>Create role (TODO)</div>} />
              <Route path="edit/:id" element={<div>Edit role (TODO)</div>} />
              <Route path="show/:id" element={<div>Show role (TODO)</div>} />
            </Route>

            <Route path="/groups">
              <Route index element={<div>Groups list (TODO)</div>} />
              <Route path="create" element={<div>Create group (TODO)</div>} />
              <Route path="edit/:id" element={<div>Edit group (TODO)</div>} />
              <Route path="show/:id" element={<div>Show group (TODO)</div>} />
            </Route>

            <Route path="/users">
              <Route index element={<div>Users list (TODO)</div>} />
              <Route path="show/:id" element={<div>User summary (TODO)</div>} />
            </Route>

            <Route path="/managed-collections">
              <Route index element={<div>Managed collections (TODO)</div>} />
            </Route>

            <Route path="/simulator" element={<div>Policy Simulator (TODO)</div>} />

            <Route path="*" element={<div>Page not found</div>} />
          </Route>

          {/* Public routes */}
          <Route
            element={
              <Authenticated key="auth-pages" fallback={<Outlet />}>
                <NavigateToResource resource="iam_policies" />
              </Authenticated>
            }
          >
            <Route path="/login" element={<LoginPage />} />
          </Route>
        </Routes>

        <UnsavedChangesNotifier />
        <DocumentTitleHandler />
      </Refine>
    </BrowserRouter>
  );
}

export default App;
```

**Step 3: Verify the app loads**

Start PocketBase backend: `go run . serve` (in project root).
Start frontend: `cd ui && npm run dev`.
Navigate to http://localhost:5173 — should redirect to `/login`.
Log in with superuser credentials — should show sidebar + "Policies list (TODO)".

**Step 4: Commit**

```bash
git add ui/src/
git commit -m "wire up Refine app shell with routing and login"
```

---

## Task 7: Update .gitignore

**Files:**
- Modify: `.gitignore`

**Step 1: Add node_modules and pb_public to .gitignore**

Append to `.gitignore`:
```
pb_public/
ui/node_modules/
ui/dist/
```

**Step 2: Commit**

```bash
git add .gitignore
git commit -m "update gitignore for frontend build output"
```

---

## Task 8: Build EntityPicker component

A reusable modal dialog for searching and selecting entities (users, roles, groups, policies).

**Files:**
- Create: `ui/src/components/entity-picker/entity-picker.tsx`, `ui/src/components/entity-picker/index.ts`

**Step 1: Create the component**

```tsx
// ui/src/components/entity-picker/entity-picker.tsx
import { useState } from "react";
import { useList } from "@refinedev/core";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus } from "lucide-react";

interface EntityPickerProps {
  /** PocketBase collection name (e.g., "iam_policies", "users") */
  resource: string;
  /** Field to display as label (e.g., "name", "email") */
  labelField: string;
  /** Field to use as secondary text (optional) */
  secondaryField?: string;
  /** Trigger button text */
  triggerLabel: string;
  /** IDs to exclude from the list (already attached) */
  excludeIds?: string[];
  /** Called when an entity is selected */
  onSelect: (id: string, record: Record<string, any>) => void;
}

export function EntityPicker({
  resource,
  labelField,
  secondaryField,
  triggerLabel,
  excludeIds = [],
  onSelect,
}: EntityPickerProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");

  const { data, isLoading } = useList({
    resource,
    filters: search
      ? [{ field: labelField, operator: "contains", value: search }]
      : [],
    pagination: { pageSize: 20 },
    queryOptions: { enabled: open },
  });

  const records = (data?.data ?? []).filter(
    (r) => !excludeIds.includes(r.id as string)
  );

  const handleSelect = (record: Record<string, any>) => {
    onSelect(record.id as string, record);
    setOpen(false);
    setSearch("");
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          <Plus className="mr-1 h-3 w-3" />
          {triggerLabel}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Select {triggerLabel.replace(/^\+ ?/, "")}</DialogTitle>
        </DialogHeader>
        <Input
          placeholder="Search..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          autoFocus
        />
        <div className="max-h-64 overflow-y-auto">
          {isLoading ? (
            <p className="py-4 text-center text-sm text-muted-foreground">Loading...</p>
          ) : records.length === 0 ? (
            <p className="py-4 text-center text-sm text-muted-foreground">No results</p>
          ) : (
            <ul className="space-y-1">
              {records.map((record) => (
                <li key={record.id as string}>
                  <button
                    type="button"
                    className="w-full rounded-md px-3 py-2 text-left text-sm hover:bg-accent"
                    onClick={() => handleSelect(record)}
                  >
                    <span className="font-medium">
                      {record[labelField] as string}
                    </span>
                    {secondaryField && record[secondaryField] && (
                      <span className="ml-2 text-muted-foreground">
                        {record[secondaryField] as string}
                      </span>
                    )}
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
```

**Step 2: Create barrel export**

```ts
// ui/src/components/entity-picker/index.ts
export { EntityPicker } from "./entity-picker";
```

**Step 3: Commit**

```bash
git add ui/src/components/entity-picker/
git commit -m "add EntityPicker component"
```

---

## Task 9: Build EntityChipList component

Displays attached entities as chips with remove buttons and an add picker.

**Files:**
- Create: `ui/src/components/entity-chip-list/entity-chip-list.tsx`, `ui/src/components/entity-chip-list/index.ts`

**Step 1: Create the component**

```tsx
// ui/src/components/entity-chip-list/entity-chip-list.tsx
import { Badge } from "@/components/ui/badge";
import { X } from "lucide-react";
import { EntityPicker } from "@/components/entity-picker";

interface ChipItem {
  id: string;
  label: string;
}

interface EntityChipListProps {
  /** Section title (e.g., "Attached Policies") */
  title: string;
  /** Currently attached items */
  items: ChipItem[];
  /** PocketBase collection to pick from */
  resource: string;
  /** Field to display as label in picker */
  labelField: string;
  /** Secondary display field in picker */
  secondaryField?: string;
  /** Trigger button text (e.g., "Attach Policy") */
  addLabel: string;
  /** Called when an entity is added */
  onAdd: (id: string, record: Record<string, any>) => void;
  /** Called when a chip is removed */
  onRemove: (id: string) => void;
  /** Loading state */
  isLoading?: boolean;
}

export function EntityChipList({
  title,
  items,
  resource,
  labelField,
  secondaryField,
  addLabel,
  onAdd,
  onRemove,
  isLoading = false,
}: EntityChipListProps) {
  return (
    <div className="space-y-2">
      <h3 className="text-sm font-medium text-muted-foreground">{title}</h3>
      <div className="flex flex-wrap items-center gap-2">
        {isLoading ? (
          <span className="text-sm text-muted-foreground">Loading...</span>
        ) : items.length === 0 ? (
          <span className="text-sm text-muted-foreground">None</span>
        ) : (
          items.map((item) => (
            <Badge key={item.id} variant="secondary" className="gap-1 pr-1">
              {item.label}
              <button
                type="button"
                onClick={() => onRemove(item.id)}
                className="ml-1 rounded-full p-0.5 hover:bg-muted"
              >
                <X className="h-3 w-3" />
              </button>
            </Badge>
          ))
        )}
        <EntityPicker
          resource={resource}
          labelField={labelField}
          secondaryField={secondaryField}
          triggerLabel={addLabel}
          excludeIds={items.map((i) => i.id)}
          onSelect={onAdd}
        />
      </div>
    </div>
  );
}
```

**Step 2: Create barrel export**

```ts
// ui/src/components/entity-chip-list/index.ts
export { EntityChipList } from "./entity-chip-list";
```

**Step 3: Commit**

```bash
git add ui/src/components/entity-chip-list/
git commit -m "add EntityChipList component"
```

---

## Task 10: Build JSON editor component (CodeMirror 6)

**Files:**
- Create: `ui/src/components/json-editor/json-editor.tsx`, `ui/src/components/json-editor/index.ts`

**Step 1: Install CodeMirror packages**

```bash
cd ui
npm install codemirror @codemirror/lang-json @codemirror/lint @codemirror/view @codemirror/state @codemirror/commands @codemirror/language @codemirror/autocomplete @codemirror/search
```

**Step 2: Create the component**

```tsx
// ui/src/components/json-editor/json-editor.tsx
import { useRef, useEffect, forwardRef, useImperativeHandle } from "react";
import { EditorState, type Extension } from "@codemirror/state";
import {
  EditorView,
  keymap,
  lineNumbers,
  highlightActiveLine,
  highlightSpecialChars,
  drawSelection,
} from "@codemirror/view";
import {
  defaultKeymap,
  history,
  historyKeymap,
  indentWithTab,
} from "@codemirror/commands";
import { json } from "@codemirror/lang-json";
import { linter, lintGutter, type Diagnostic } from "@codemirror/lint";
import {
  syntaxHighlighting,
  defaultHighlightStyle,
  indentOnInput,
  bracketMatching,
  foldGutter,
  foldKeymap,
} from "@codemirror/language";
import { closeBrackets, closeBracketsKeymap } from "@codemirror/autocomplete";
import { searchKeymap, highlightSelectionMatches } from "@codemirror/search";
import { cn } from "@/lib/utils";

const jsonLinter = linter((view) => {
  const diagnostics: Diagnostic[] = [];
  const doc = view.state.doc.toString();
  if (doc.trim().length === 0) return diagnostics;

  try {
    JSON.parse(doc);
  } catch (e) {
    if (e instanceof SyntaxError) {
      const posMatch = e.message.match(/position\s+(\d+)/i);
      const pos = posMatch ? Number(posMatch[1]) : 0;
      diagnostics.push({
        from: Math.min(pos, doc.length),
        to: Math.min(pos, doc.length),
        severity: "error",
        message: e.message,
      });
    }
  }
  return diagnostics;
});

export interface JsonEditorHandle {
  getValue: () => string;
  setValue: (value: string) => void;
}

interface JsonEditorProps {
  defaultValue?: string;
  onChange?: (value: string) => void;
  readOnly?: boolean;
  className?: string;
  resetKey?: string | number;
}

export const JsonEditor = forwardRef<JsonEditorHandle, JsonEditorProps>(
  function JsonEditor(
    { defaultValue = "", onChange, readOnly = false, className, resetKey },
    ref
  ) {
    const containerRef = useRef<HTMLDivElement>(null);
    const viewRef = useRef<EditorView | null>(null);
    const onChangeRef = useRef(onChange);
    onChangeRef.current = onChange;

    useImperativeHandle(ref, () => ({
      getValue() {
        return viewRef.current?.state.doc.toString() ?? "";
      },
      setValue(value: string) {
        const view = viewRef.current;
        if (!view) return;
        view.dispatch({
          changes: { from: 0, to: view.state.doc.length, insert: value },
        });
      },
    }));

    useEffect(() => {
      if (!containerRef.current) return;
      viewRef.current?.destroy();

      const extensions: Extension[] = [
        lineNumbers(),
        highlightActiveLine(),
        highlightSpecialChars(),
        drawSelection(),
        indentOnInput(),
        bracketMatching(),
        closeBrackets(),
        foldGutter(),
        history(),
        highlightSelectionMatches(),
        syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
        keymap.of([
          ...closeBracketsKeymap,
          ...defaultKeymap,
          ...historyKeymap,
          ...foldKeymap,
          ...searchKeymap,
          indentWithTab,
        ]),
        json(),
        jsonLinter,
        lintGutter(),
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            onChangeRef.current?.(update.state.doc.toString());
          }
        }),
        EditorView.theme({
          "&": { minHeight: "200px", maxHeight: "500px" },
          ".cm-scroller": { overflow: "auto" },
        }),
      ];

      if (readOnly) {
        extensions.push(EditorState.readOnly.of(true));
        extensions.push(EditorView.editable.of(false));
      }

      const view = new EditorView({
        state: EditorState.create({ doc: defaultValue, extensions }),
        parent: containerRef.current,
      });
      viewRef.current = view;

      return () => {
        view.destroy();
        viewRef.current = null;
      };
    }, [resetKey, readOnly]);

    return (
      <div
        ref={containerRef}
        className={cn(
          "overflow-hidden rounded-md border",
          className
        )}
      />
    );
  }
);
```

**Step 3: Create barrel export**

```ts
// ui/src/components/json-editor/index.ts
export { JsonEditor, type JsonEditorHandle } from "./json-editor";
```

**Step 4: Commit**

```bash
git add ui/src/components/json-editor/
git commit -m "add CodeMirror 6 JSON editor component"
```

---

## Task 11: Build PolicyFormBuilder component

Visual editor for policy document statements — the form builder mode.

**Files:**
- Create: `ui/src/components/policy-form-builder/statement-card.tsx`, `ui/src/components/policy-form-builder/tag-list-input.tsx`, `ui/src/components/policy-form-builder/policy-form-builder.tsx`, `ui/src/components/policy-form-builder/index.ts`
- Create: `ui/src/types/policy.ts`

**Step 1: Create shared policy types**

```ts
// ui/src/types/policy.ts
export interface PolicyStatement {
  sid: string;
  effect: "Allow" | "Deny";
  action: string[];
  resource: string[];
}

export interface PolicyDocument {
  version: string;
  statement: PolicyStatement[];
}

export const DEFAULT_STATEMENT: PolicyStatement = {
  sid: "",
  effect: "Allow",
  action: [""],
  resource: ["*"],
};

export const DEFAULT_POLICY_DOCUMENT: PolicyDocument = {
  version: "2024-01-01",
  statement: [{ ...DEFAULT_STATEMENT }],
};
```

**Step 2: Create TagListInput**

```tsx
// ui/src/components/policy-form-builder/tag-list-input.tsx
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus, X } from "lucide-react";

interface TagListInputProps {
  label: string;
  values: string[];
  onChange: (values: string[]) => void;
  placeholder?: string;
}

export function TagListInput({
  label,
  values,
  onChange,
  placeholder = "Enter value...",
}: TagListInputProps) {
  const updateValue = (index: number, value: string) => {
    const next = [...values];
    next[index] = value;
    onChange(next);
  };

  const removeValue = (index: number) => {
    onChange(values.filter((_, i) => i !== index));
  };

  const addValue = () => {
    onChange([...values, ""]);
  };

  return (
    <div className="space-y-1">
      <span className="text-xs font-medium text-muted-foreground">{label}</span>
      <div className="flex flex-wrap items-center gap-1.5">
        {values.map((v, i) => (
          <div key={i} className="flex items-center gap-1">
            <Input
              value={v}
              onChange={(e) => updateValue(i, e.target.value)}
              placeholder={placeholder}
              className="h-7 w-56 text-xs"
            />
            {values.length > 1 && (
              <button
                type="button"
                onClick={() => removeValue(i)}
                className="rounded p-0.5 hover:bg-muted"
              >
                <X className="h-3 w-3" />
              </button>
            )}
          </div>
        ))}
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="h-7 text-xs"
          onClick={addValue}
        >
          <Plus className="mr-1 h-3 w-3" />
          Add
        </Button>
      </div>
    </div>
  );
}
```

**Step 3: Create StatementCard**

```tsx
// ui/src/components/policy-form-builder/statement-card.tsx
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";
import { TagListInput } from "./tag-list-input";
import type { PolicyStatement } from "@/types/policy";

interface StatementCardProps {
  index: number;
  statement: PolicyStatement;
  onChange: (statement: PolicyStatement) => void;
  onRemove: () => void;
  canRemove: boolean;
}

export function StatementCard({
  index,
  statement,
  onChange,
  onRemove,
  canRemove,
}: StatementCardProps) {
  return (
    <Card>
      <CardContent className="space-y-3 pt-4">
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium">Statement {index + 1}</span>
          {canRemove && (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={onRemove}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          )}
        </div>

        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-1">
            <Label className="text-xs">SID</Label>
            <Input
              value={statement.sid}
              onChange={(e) =>
                onChange({ ...statement, sid: e.target.value })
              }
              placeholder="e.g. AllowReadPosts"
              className="h-8 text-sm"
            />
          </div>

          <div className="space-y-1">
            <Label className="text-xs">Effect</Label>
            <div className="flex gap-4 pt-1">
              <label className="flex items-center gap-1.5 text-sm">
                <input
                  type="radio"
                  name={`effect-${index}`}
                  checked={statement.effect === "Allow"}
                  onChange={() =>
                    onChange({ ...statement, effect: "Allow" })
                  }
                />
                Allow
              </label>
              <label className="flex items-center gap-1.5 text-sm">
                <input
                  type="radio"
                  name={`effect-${index}`}
                  checked={statement.effect === "Deny"}
                  onChange={() =>
                    onChange({ ...statement, effect: "Deny" })
                  }
                />
                Deny
              </label>
            </div>
          </div>
        </div>

        <TagListInput
          label="Actions"
          values={statement.action}
          onChange={(action) => onChange({ ...statement, action })}
          placeholder="e.g. collections:posts:read"
        />

        <TagListInput
          label="Resources"
          values={statement.resource}
          onChange={(resource) => onChange({ ...statement, resource })}
          placeholder="e.g. *"
        />
      </CardContent>
    </Card>
  );
}
```

**Step 4: Create PolicyFormBuilder**

```tsx
// ui/src/components/policy-form-builder/policy-form-builder.tsx
import { useState, useRef } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Plus } from "lucide-react";
import { StatementCard } from "./statement-card";
import { JsonEditor, type JsonEditorHandle } from "@/components/json-editor";
import type { PolicyDocument, PolicyStatement } from "@/types/policy";
import { DEFAULT_STATEMENT, DEFAULT_POLICY_DOCUMENT } from "@/types/policy";

interface PolicyFormBuilderProps {
  value: PolicyDocument;
  onChange: (doc: PolicyDocument) => void;
}

export function PolicyFormBuilder({ value, onChange }: PolicyFormBuilderProps) {
  const [mode, setMode] = useState<"form" | "json">("form");
  const editorRef = useRef<JsonEditorHandle>(null);
  const [jsonResetKey, setJsonResetKey] = useState(0);

  const updateStatement = (index: number, stmt: PolicyStatement) => {
    const next = [...value.statement];
    next[index] = stmt;
    onChange({ ...value, statement: next });
  };

  const removeStatement = (index: number) => {
    onChange({
      ...value,
      statement: value.statement.filter((_, i) => i !== index),
    });
  };

  const addStatement = () => {
    onChange({
      ...value,
      statement: [...value.statement, { ...DEFAULT_STATEMENT }],
    });
  };

  const handleModeChange = (newMode: string) => {
    if (newMode === "json") {
      // Sync form state to JSON editor
      setJsonResetKey((k) => k + 1);
    } else if (newMode === "form") {
      // Sync JSON editor state back to form
      const raw = editorRef.current?.getValue();
      if (raw) {
        try {
          const parsed = JSON.parse(raw) as PolicyDocument;
          onChange(parsed);
        } catch {
          // Invalid JSON — stay on JSON tab
          return;
        }
      }
    }
    setMode(newMode as "form" | "json");
  };

  return (
    <div className="space-y-4">
      <Tabs value={mode} onValueChange={handleModeChange}>
        <TabsList>
          <TabsTrigger value="form">Form Builder</TabsTrigger>
          <TabsTrigger value="json">Raw JSON</TabsTrigger>
        </TabsList>

        <TabsContent value="form" className="space-y-4">
          <div className="space-y-1">
            <Label className="text-xs">Version</Label>
            <Input
              value={value.version}
              onChange={(e) => onChange({ ...value, version: e.target.value })}
              placeholder="2024-01-01"
              className="h-8 w-48 text-sm"
            />
          </div>

          {value.statement.map((stmt, i) => (
            <StatementCard
              key={i}
              index={i}
              statement={stmt}
              onChange={(s) => updateStatement(i, s)}
              onRemove={() => removeStatement(i)}
              canRemove={value.statement.length > 1}
            />
          ))}

          <Button type="button" variant="outline" onClick={addStatement}>
            <Plus className="mr-1 h-4 w-4" />
            Add Statement
          </Button>
        </TabsContent>

        <TabsContent value="json">
          <JsonEditor
            ref={editorRef}
            defaultValue={JSON.stringify(value, null, 2)}
            resetKey={jsonResetKey}
            onChange={(raw) => {
              try {
                const parsed = JSON.parse(raw);
                onChange(parsed);
              } catch {
                // Invalid JSON — don't update form state
              }
            }}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
```

**Step 5: Create barrel export**

```ts
// ui/src/components/policy-form-builder/index.ts
export { PolicyFormBuilder } from "./policy-form-builder";
export { StatementCard } from "./statement-card";
export { TagListInput } from "./tag-list-input";
```

**Step 6: Commit**

```bash
git add ui/src/types/ ui/src/components/policy-form-builder/
git commit -m "add PolicyFormBuilder with statement cards and JSON editor"
```

---

## Task 12: Build EffectivePermissionsTable component

Read-only table showing merged Allow/Deny per action with source attribution.

**Files:**
- Create: `ui/src/components/effective-permissions-table/effective-permissions-table.tsx`, `ui/src/components/effective-permissions-table/index.ts`
- Create: `ui/src/lib/permissions.ts`

**Step 1: Create the permissions evaluation utility**

```ts
// ui/src/lib/permissions.ts
import type { PolicyDocument } from "@/types/policy";

export interface EffectivePermission {
  action: string;
  effect: "Allow" | "Deny";
  source: string; // e.g., "Role: Editor", "Direct: PolicyName", "Group: Devs"
  policyName: string;
  sid: string;
}

interface PolicyWithSource {
  policyName: string;
  document: PolicyDocument;
  source: string;
}

/**
 * Merges policies and computes effective permissions.
 * Deny overrides Allow (AWS IAM model).
 */
export function computeEffectivePermissions(
  policies: PolicyWithSource[]
): EffectivePermission[] {
  const permissions = new Map<string, EffectivePermission>();

  for (const { policyName, document: doc, source } of policies) {
    for (const stmt of doc.statement) {
      for (const action of stmt.action) {
        const existing = permissions.get(action);
        // Deny always wins
        if (existing && existing.effect === "Deny") continue;
        // New Deny overrides existing Allow
        if (stmt.effect === "Deny" || !existing) {
          permissions.set(action, {
            action,
            effect: stmt.effect as "Allow" | "Deny",
            source,
            policyName,
            sid: stmt.sid,
          });
        }
      }
    }
  }

  return Array.from(permissions.values()).sort((a, b) =>
    a.action.localeCompare(b.action)
  );
}
```

**Step 2: Create the component**

```tsx
// ui/src/components/effective-permissions-table/effective-permissions-table.tsx
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import type { EffectivePermission } from "@/lib/permissions";

interface EffectivePermissionsTableProps {
  permissions: EffectivePermission[];
  showSource?: boolean;
}

export function EffectivePermissionsTable({
  permissions,
  showSource = false,
}: EffectivePermissionsTableProps) {
  if (permissions.length === 0) {
    return (
      <p className="py-4 text-sm text-muted-foreground">
        No permissions to display.
      </p>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Action</TableHead>
          <TableHead>Effect</TableHead>
          {showSource && <TableHead>Source</TableHead>}
        </TableRow>
      </TableHeader>
      <TableBody>
        {permissions.map((p) => (
          <TableRow key={p.action}>
            <TableCell className="font-mono text-sm">{p.action}</TableCell>
            <TableCell>
              <Badge
                variant={p.effect === "Allow" ? "default" : "destructive"}
              >
                {p.effect}
              </Badge>
            </TableCell>
            {showSource && (
              <TableCell className="text-sm text-muted-foreground">
                {p.source}
              </TableCell>
            )}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
```

**Step 3: Create barrel export**

```ts
// ui/src/components/effective-permissions-table/index.ts
export { EffectivePermissionsTable } from "./effective-permissions-table";
```

**Step 4: Commit**

```bash
git add ui/src/lib/permissions.ts ui/src/components/effective-permissions-table/
git commit -m "add EffectivePermissionsTable component"
```

---

## Task 13: Build Policies list page

**Files:**
- Create: `ui/src/pages/policies/list.tsx`, `ui/src/pages/policies/index.ts`
- Modify: `ui/src/App.tsx` — replace placeholder route

**Step 1: Create the policies list page**

```tsx
// ui/src/pages/policies/list.tsx
import { useTable } from "@refinedev/react-table";
import { useNavigation, useDelete } from "@refinedev/core";
import {
  type ColumnDef,
  flexRender,
  getCoreRowModel,
} from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Pencil, Trash2 } from "lucide-react";

interface IPolicy {
  id: string;
  name: string;
  description: string;
  document: { statement: unknown[] };
}

export function PolicyList() {
  const { create, edit } = useNavigation();
  const { mutate: deleteRecord } = useDelete();

  const columns: ColumnDef<IPolicy>[] = [
    { accessorKey: "name", header: "Name" },
    {
      accessorKey: "description",
      header: "Description",
      cell: ({ getValue }) => {
        const v = getValue<string>();
        return v && v.length > 60 ? v.slice(0, 60) + "..." : v || "—";
      },
    },
    {
      id: "statements",
      header: "Statements",
      cell: ({ row }) => {
        const doc = row.original.document;
        return doc?.statement?.length ?? 0;
      },
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => (
        <div className="flex gap-1">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => edit("iam_policies", row.original.id)}
          >
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() =>
              deleteRecord({
                resource: "iam_policies",
                id: row.original.id,
              })
            }
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      ),
    },
  ];

  const table = useTable<IPolicy>({
    columns,
    refineCoreProps: { resource: "iam_policies" },
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-2xl font-bold">Policies</h1>
        <Button onClick={() => create("iam_policies")}>Create Policy</Button>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((hg) => (
              <TableRow key={hg.id}>
                {hg.headers.map((h) => (
                  <TableHead key={h.id}>
                    {h.isPlaceholder
                      ? null
                      : flexRender(h.column.columnDef.header, h.getContext())}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={columns.length} className="h-24 text-center">
                  No policies found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-end gap-2 py-4">
        <Button
          variant="outline"
          size="sm"
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          Next
        </Button>
      </div>
    </div>
  );
}
```

**Step 2: Create barrel export**

```ts
// ui/src/pages/policies/index.ts
export { PolicyList } from "./list";
```

**Step 3: Update App.tsx** — replace the policies route placeholder:

```tsx
// In the /policies route section, replace:
<Route index element={<div>Policies list (TODO)</div>} />
// With:
<Route index element={<PolicyList />} />
```

Add the import at top of App.tsx:
```tsx
import { PolicyList } from "@/pages/policies";
```

**Step 4: Verify**

With PB running, navigate to `/policies`. Should show the table (empty if no policies exist yet).

**Step 5: Commit**

```bash
git add ui/src/pages/policies/ ui/src/App.tsx
git commit -m "add Policies list page"
```

---

## Task 14: Build Policies create/edit pages

**Files:**
- Create: `ui/src/pages/policies/create.tsx`, `ui/src/pages/policies/edit.tsx`
- Modify: `ui/src/pages/policies/index.ts`, `ui/src/App.tsx`

**Step 1: Create the policy create page**

```tsx
// ui/src/pages/policies/create.tsx
import { useState } from "react";
import { useForm } from "@refinedev/react-hook-form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PolicyFormBuilder } from "@/components/policy-form-builder";
import type { PolicyDocument } from "@/types/policy";
import { DEFAULT_POLICY_DOCUMENT } from "@/types/policy";

export function PolicyCreate() {
  const [policyDoc, setPolicyDoc] = useState<PolicyDocument>({
    ...DEFAULT_POLICY_DOCUMENT,
  });

  const {
    refineCore: { onFinish, formLoading },
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    refineCoreProps: {
      resource: "iam_policies",
      action: "create",
    },
  });

  const onSubmit = (data: Record<string, any>) => {
    onFinish({ ...data, document: policyDoc });
  };

  return (
    <div className="max-w-3xl">
      <h1 className="mb-6 text-2xl font-bold">Create Policy</h1>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="name">Name</Label>
          <Input
            id="name"
            {...register("name", { required: "Name is required" })}
            placeholder="e.g. ReadOnlyPosts"
          />
          {errors.name && (
            <p className="text-sm text-destructive">
              {errors.name.message as string}
            </p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="description">Description</Label>
          <Input
            id="description"
            {...register("description")}
            placeholder="Optional description"
          />
        </div>

        <div className="space-y-2">
          <Label>Policy Document</Label>
          <PolicyFormBuilder value={policyDoc} onChange={setPolicyDoc} />
        </div>

        <Button type="submit" disabled={formLoading}>
          {formLoading ? "Saving..." : "Create Policy"}
        </Button>
      </form>
    </div>
  );
}
```

**Step 2: Create the policy edit page**

```tsx
// ui/src/pages/policies/edit.tsx
import { useState, useEffect } from "react";
import { useForm } from "@refinedev/react-hook-form";
import { useList, useCreate, useDelete } from "@refinedev/core";
import { useParams } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { PolicyFormBuilder } from "@/components/policy-form-builder";
import { EntityChipList } from "@/components/entity-chip-list";
import type { PolicyDocument } from "@/types/policy";
import { DEFAULT_POLICY_DOCUMENT } from "@/types/policy";

export function PolicyEdit() {
  const { id } = useParams();
  const [policyDoc, setPolicyDoc] = useState<PolicyDocument>(DEFAULT_POLICY_DOCUMENT);

  const {
    refineCore: { onFinish, formLoading, query },
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    refineCoreProps: {
      resource: "iam_policies",
      action: "edit",
      id,
    },
  });

  // Initialize policy document from fetched data
  useEffect(() => {
    if (query?.data?.data?.document) {
      setPolicyDoc(query.data.data.document as PolicyDocument);
    }
  }, [query?.data?.data?.document]);

  const onSubmit = (data: Record<string, any>) => {
    onFinish({ ...data, document: policyDoc });
  };

  // --- Attachments ---
  const { mutate: createJoin } = useCreate();
  const { mutate: deleteJoin } = useDelete();

  // Users attached via iam_user_policies
  const { data: userPolicies } = useList({
    resource: "iam_user_policies",
    filters: [{ field: "policy", operator: "eq", value: id }],
    meta: { expand: ["user"] },
  });
  const attachedUsers = (userPolicies?.data ?? []).map((r: any) => ({
    id: r.id,
    entityId: r.user,
    label: r.expand?.user?.email ?? r.user,
  }));

  // Roles attached via iam_role_policies
  const { data: rolePolicies } = useList({
    resource: "iam_role_policies",
    filters: [{ field: "policy", operator: "eq", value: id }],
    meta: { expand: ["role"] },
  });
  const attachedRoles = (rolePolicies?.data ?? []).map((r: any) => ({
    id: r.id,
    entityId: r.role,
    label: r.expand?.role?.name ?? r.role,
  }));

  // Groups attached via iam_group_policies
  const { data: groupPolicies } = useList({
    resource: "iam_group_policies",
    filters: [{ field: "policy", operator: "eq", value: id }],
    meta: { expand: ["group"] },
  });
  const attachedGroups = (groupPolicies?.data ?? []).map((r: any) => ({
    id: r.id,
    entityId: r.group,
    label: r.expand?.group?.name ?? r.group,
  }));

  return (
    <div className="max-w-3xl">
      <h1 className="mb-6 text-2xl font-bold">Edit Policy</h1>
      {query?.isLoading ? (
        <p>Loading...</p>
      ) : (
        <>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                {...register("name", { required: "Name is required" })}
              />
              {errors.name && (
                <p className="text-sm text-destructive">
                  {errors.name.message as string}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Input id="description" {...register("description")} />
            </div>

            <div className="space-y-2">
              <Label>Policy Document</Label>
              <PolicyFormBuilder value={policyDoc} onChange={setPolicyDoc} />
            </div>

            <Button type="submit" disabled={formLoading}>
              {formLoading ? "Saving..." : "Save Policy"}
            </Button>
          </form>

          <Separator className="my-8" />

          <div className="space-y-6">
            <h2 className="text-lg font-semibold">Attachments</h2>

            <EntityChipList
              title="Direct Users"
              items={attachedUsers.map((u) => ({ id: u.id, label: u.label }))}
              resource="users"
              labelField="email"
              addLabel="Attach User"
              onAdd={(userId) =>
                createJoin({
                  resource: "iam_user_policies",
                  values: { user: userId, policy: id },
                })
              }
              onRemove={(joinId) =>
                deleteJoin({ resource: "iam_user_policies", id: joinId })
              }
            />

            <EntityChipList
              title="Roles"
              items={attachedRoles.map((r) => ({ id: r.id, label: r.label }))}
              resource="iam_roles"
              labelField="name"
              addLabel="Attach to Role"
              onAdd={(roleId) =>
                createJoin({
                  resource: "iam_role_policies",
                  values: { role: roleId, policy: id },
                })
              }
              onRemove={(joinId) =>
                deleteJoin({ resource: "iam_role_policies", id: joinId })
              }
            />

            <EntityChipList
              title="Groups"
              items={attachedGroups.map((g) => ({ id: g.id, label: g.label }))}
              resource="iam_groups"
              labelField="name"
              addLabel="Attach to Group"
              onAdd={(groupId) =>
                createJoin({
                  resource: "iam_group_policies",
                  values: { group: groupId, policy: id },
                })
              }
              onRemove={(joinId) =>
                deleteJoin({ resource: "iam_group_policies", id: joinId })
              }
            />
          </div>
        </>
      )}
    </div>
  );
}
```

**Step 3: Update barrel and App.tsx**

Add to `ui/src/pages/policies/index.ts`:
```ts
export { PolicyCreate } from "./create";
export { PolicyEdit } from "./edit";
```

Update App.tsx routes and imports for policies create/edit.

**Step 4: Commit**

```bash
git add ui/src/pages/policies/ ui/src/App.tsx
git commit -m "add Policy create and edit pages with attachments"
```

---

## Task 15: Build Roles list and detail pages

**Files:**
- Create: `ui/src/pages/roles/list.tsx`, `ui/src/pages/roles/create.tsx`, `ui/src/pages/roles/show.tsx`, `ui/src/pages/roles/index.ts`
- Modify: `ui/src/App.tsx`

**Step 1: Create roles list page**

Same DataTable pattern as PolicyList, columns: Name, Description, Policies (count), Users (count), Actions (Edit, Delete).

**Step 2: Create roles create page**

Simple form with Name + Description fields using `useForm`.

**Step 3: Create roles show/detail page**

This is the key page — it shows:
- Name + Description (editable)
- EntityChipList for attached policies (via `iam_role_policies`)
- EntityChipList for assigned users (via `iam_user_roles`)
- EffectivePermissionsTable computed from attached policies

```tsx
// ui/src/pages/roles/show.tsx
import { useShow, useList, useCreate, useDelete } from "@refinedev/core";
import { useParams, Link } from "react-router";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { Pencil } from "lucide-react";
import { EntityChipList } from "@/components/entity-chip-list";
import { EffectivePermissionsTable } from "@/components/effective-permissions-table";
import { computeEffectivePermissions } from "@/lib/permissions";
import type { PolicyDocument } from "@/types/policy";

export function RoleShow() {
  const { id } = useParams();
  const { query } = useShow({ resource: "iam_roles", id });
  const record = query?.data?.data;

  const { mutate: createJoin } = useCreate();
  const { mutate: deleteJoin } = useDelete();

  // Attached policies
  const { data: rolePolicies } = useList({
    resource: "iam_role_policies",
    filters: [{ field: "role", operator: "eq", value: id }],
    meta: { expand: ["policy"] },
  });
  const attachedPolicies = (rolePolicies?.data ?? []).map((r: any) => ({
    id: r.id,
    entityId: r.policy,
    label: r.expand?.policy?.name ?? r.policy,
    document: r.expand?.policy?.document as PolicyDocument | undefined,
  }));

  // Assigned users
  const { data: userRoles } = useList({
    resource: "iam_user_roles",
    filters: [{ field: "role", operator: "eq", value: id }],
    meta: { expand: ["user"] },
  });
  const assignedUsers = (userRoles?.data ?? []).map((r: any) => ({
    id: r.id,
    entityId: r.user,
    label: r.expand?.user?.email ?? r.user,
  }));

  // Effective permissions
  const effectivePerms = computeEffectivePermissions(
    attachedPolicies
      .filter((p) => p.document)
      .map((p) => ({
        policyName: p.label,
        document: p.document!,
        source: `Role: ${record?.name ?? ""}`,
      }))
  );

  if (query?.isLoading) return <p>Loading...</p>;

  return (
    <div className="max-w-3xl space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{record?.name}</h1>
        <Button variant="outline" asChild>
          <Link to={`/roles/edit/${id}`}>
            <Pencil className="mr-1 h-4 w-4" /> Edit
          </Link>
        </Button>
      </div>

      {record?.description && (
        <p className="text-muted-foreground">{record.description}</p>
      )}

      <Separator />

      <EntityChipList
        title="Attached Policies"
        items={attachedPolicies.map((p) => ({ id: p.id, label: p.label }))}
        resource="iam_policies"
        labelField="name"
        addLabel="Attach Policy"
        onAdd={(policyId) =>
          createJoin({
            resource: "iam_role_policies",
            values: { role: id, policy: policyId },
          })
        }
        onRemove={(joinId) =>
          deleteJoin({ resource: "iam_role_policies", id: joinId })
        }
      />

      <EntityChipList
        title="Assigned Users"
        items={assignedUsers.map((u) => ({ id: u.id, label: u.label }))}
        resource="users"
        labelField="email"
        addLabel="Assign User"
        onAdd={(userId) =>
          createJoin({
            resource: "iam_user_roles",
            values: { user: userId, role: id },
          })
        }
        onRemove={(joinId) =>
          deleteJoin({ resource: "iam_user_roles", id: joinId })
        }
      />

      <Separator />

      <div>
        <h2 className="mb-3 text-lg font-semibold">Effective Permissions</h2>
        <EffectivePermissionsTable permissions={effectivePerms} />
      </div>
    </div>
  );
}
```

**Step 4: Update barrel and App.tsx**

Wire routes: list, create, edit (reuse simple form), show.

**Step 5: Commit**

```bash
git add ui/src/pages/roles/ ui/src/App.tsx
git commit -m "add Roles list and detail pages"
```

---

## Task 16: Build Groups list and detail pages

**Files:**
- Create: `ui/src/pages/groups/list.tsx`, `ui/src/pages/groups/create.tsx`, `ui/src/pages/groups/show.tsx`, `ui/src/pages/groups/index.ts`
- Modify: `ui/src/App.tsx`

Same pattern as Roles (Task 15) with two differences:
- **Members** section (EntityChipList for `iam_group_users`) instead of "Assigned Users"
- **Attached Policies** section (EntityChipList for `iam_group_policies`)
- EffectivePermissionsTable included

**Step 1-4:** Follow the same structure as Task 15, substituting group-specific collections.

**Step 5: Commit**

```bash
git add ui/src/pages/groups/ ui/src/App.tsx
git commit -m "add Groups list and detail pages"
```

---

## Task 17: Build Users list and IAM summary pages

**Files:**
- Create: `ui/src/pages/users/list.tsx`, `ui/src/pages/users/show.tsx`, `ui/src/pages/users/index.ts`
- Modify: `ui/src/App.tsx`

**Step 1: Create users list page**

DataTable with columns: Email, Roles (count), Groups (count), Direct Policies (count), Actions (View).

Users come from PB's `users` collection — read-only, no create/edit/delete.

**Step 2: Create user IAM summary page**

This is the most important page. It shows:
- User info (email)
- EntityChipList: Roles (via `iam_user_roles`)
- EntityChipList: Groups (via `iam_group_users`)
- EntityChipList: Direct Policies (via `iam_user_policies`)
- EffectivePermissionsTable with `showSource={true}` — collects policies from all 3 sources, merges them, shows source attribution
- "Test in Simulator" link

```tsx
// ui/src/pages/users/show.tsx
import { useShow, useList, useCreate, useDelete } from "@refinedev/core";
import { useParams, Link } from "react-router";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { FlaskConical } from "lucide-react";
import { EntityChipList } from "@/components/entity-chip-list";
import { EffectivePermissionsTable } from "@/components/effective-permissions-table";
import { computeEffectivePermissions } from "@/lib/permissions";
import type { PolicyDocument } from "@/types/policy";

export function UserShow() {
  const { id } = useParams();
  const { query } = useShow({ resource: "users", id });
  const user = query?.data?.data;

  const { mutate: createJoin } = useCreate();
  const { mutate: deleteJoin } = useDelete();

  // User's roles
  const { data: userRoles } = useList({
    resource: "iam_user_roles",
    filters: [{ field: "user", operator: "eq", value: id }],
    meta: { expand: ["role"] },
  });
  const roles = (userRoles?.data ?? []).map((r: any) => ({
    id: r.id,
    entityId: r.role,
    label: r.expand?.role?.name ?? r.role,
  }));

  // User's groups
  const { data: groupUsers } = useList({
    resource: "iam_group_users",
    filters: [{ field: "user", operator: "eq", value: id }],
    meta: { expand: ["group"] },
  });
  const groups = (groupUsers?.data ?? []).map((r: any) => ({
    id: r.id,
    entityId: r.group,
    label: r.expand?.group?.name ?? r.group,
  }));

  // User's direct policies
  const { data: userPolicies } = useList({
    resource: "iam_user_policies",
    filters: [{ field: "user", operator: "eq", value: id }],
    meta: { expand: ["policy"] },
  });
  const directPolicies = (userPolicies?.data ?? []).map((r: any) => ({
    id: r.id,
    entityId: r.policy,
    label: r.expand?.policy?.name ?? r.policy,
    document: r.expand?.policy?.document as PolicyDocument | undefined,
  }));

  // Role policies (need to fetch for each role)
  const roleIds = roles.map((r) => r.entityId);
  const { data: rolePoliciesData } = useList({
    resource: "iam_role_policies",
    filters: roleIds.length
      ? [{ field: "role", operator: "in", value: roleIds }]
      : [],
    meta: { expand: ["policy", "role"] },
    queryOptions: { enabled: roleIds.length > 0 },
    pagination: { pageSize: 100 },
  });

  // Group policies (need to fetch for each group)
  const groupIds = groups.map((g) => g.entityId);
  const { data: groupPoliciesData } = useList({
    resource: "iam_group_policies",
    filters: groupIds.length
      ? [{ field: "group", operator: "in", value: groupIds }]
      : [],
    meta: { expand: ["policy", "group"] },
    queryOptions: { enabled: groupIds.length > 0 },
    pagination: { pageSize: 100 },
  });

  // Compute effective permissions from all sources
  const allPolicySources = [
    ...directPolicies
      .filter((p) => p.document)
      .map((p) => ({
        policyName: p.label,
        document: p.document!,
        source: `Direct: ${p.label}`,
      })),
    ...(rolePoliciesData?.data ?? [])
      .filter((r: any) => r.expand?.policy?.document)
      .map((r: any) => ({
        policyName: r.expand.policy.name,
        document: r.expand.policy.document as PolicyDocument,
        source: `Role: ${r.expand?.role?.name ?? r.role}`,
      })),
    ...(groupPoliciesData?.data ?? [])
      .filter((r: any) => r.expand?.policy?.document)
      .map((r: any) => ({
        policyName: r.expand.policy.name,
        document: r.expand.policy.document as PolicyDocument,
        source: `Group: ${r.expand?.group?.name ?? r.group}`,
      })),
  ];

  const effectivePerms = computeEffectivePermissions(allPolicySources);

  if (query?.isLoading) return <p>Loading...</p>;

  return (
    <div className="max-w-3xl space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{user?.email}</h1>
        <Button variant="outline" asChild>
          <Link to={`/simulator?user=${id}`}>
            <FlaskConical className="mr-1 h-4 w-4" /> Test in Simulator
          </Link>
        </Button>
      </div>

      <EntityChipList
        title="Roles"
        items={roles.map((r) => ({ id: r.id, label: r.label }))}
        resource="iam_roles"
        labelField="name"
        addLabel="Assign Role"
        onAdd={(roleId) =>
          createJoin({
            resource: "iam_user_roles",
            values: { user: id, role: roleId },
          })
        }
        onRemove={(joinId) =>
          deleteJoin({ resource: "iam_user_roles", id: joinId })
        }
      />

      <EntityChipList
        title="Groups"
        items={groups.map((g) => ({ id: g.id, label: g.label }))}
        resource="iam_groups"
        labelField="name"
        addLabel="Add to Group"
        onAdd={(groupId) =>
          createJoin({
            resource: "iam_group_users",
            values: { user: id, group: groupId },
          })
        }
        onRemove={(joinId) =>
          deleteJoin({ resource: "iam_group_users", id: joinId })
        }
      />

      <EntityChipList
        title="Direct Policies"
        items={directPolicies.map((p) => ({ id: p.id, label: p.label }))}
        resource="iam_policies"
        labelField="name"
        addLabel="Attach Policy"
        onAdd={(policyId) =>
          createJoin({
            resource: "iam_user_policies",
            values: { user: id, policy: policyId },
          })
        }
        onRemove={(joinId) =>
          deleteJoin({ resource: "iam_user_policies", id: joinId })
        }
      />

      <Separator />

      <div>
        <h2 className="mb-3 text-lg font-semibold">All Effective Permissions</h2>
        <EffectivePermissionsTable permissions={effectivePerms} showSource />
      </div>
    </div>
  );
}
```

**Step 3: Wire routes and commit**

```bash
git add ui/src/pages/users/ ui/src/App.tsx
git commit -m "add Users list and IAM summary pages"
```

---

## Task 18: Build Managed Collections page

**Files:**
- Create: `ui/src/pages/managed-collections/list.tsx`, `ui/src/pages/managed-collections/index.ts`
- Modify: `ui/src/App.tsx`

**Step 1: Create the page**

Two-section layout:
- **Currently Managed**: fetches `iam_managed_collections`, shows table with Unregister button
- **Available Collections**: fetches all PB collections via `useCustom` to `/api/collections`, filters out already-managed ones, shows Register button

Each action (register/unregister) uses a confirmation dialog (`AlertDialog` from shadcn).

Register = `useCreate` on `iam_managed_collections` with `{ collection_name: "..." }`.
Unregister = `useDelete` on `iam_managed_collections`.

**Step 2: Wire route and commit**

```bash
git add ui/src/pages/managed-collections/ ui/src/App.tsx
git commit -m "add Managed Collections page"
```

---

## Task 19: Add /api/iam/simulate backend endpoint

The Policy Simulator needs a verbose endpoint that returns evaluation trace.

**Files:**
- Modify: `iam/routes.go`
- Modify: `iam/engine.go` — add `EvaluateVerbose` function that returns trace details

**Step 1: Add EvaluateVerbose to engine.go**

Returns structured trace info: policies collected (count + sources), statements checked, deny/allow matches with policy name + source.

**Step 2: Add the route to routes.go**

```go
// POST /api/iam/simulate (superuser-only)
se.Router.POST("/api/iam/simulate", func(e *core.RequestEvent) error {
    var body struct {
        UserID   string `json:"user_id"`
        Action   string `json:"action"`
        Resource string `json:"resource"`
    }
    // ... bind + validate
    // ... call EvaluateVerbose
    // ... return full trace JSON
}).Bind(apis.RequireSuperuserAuth())
```

**Step 3: Commit**

```bash
git add iam/routes.go iam/engine.go
git commit -m "add /api/iam/simulate endpoint for policy simulator"
```

---

## Task 20: Build Policy Simulator page

**Files:**
- Create: `ui/src/pages/simulator/simulator.tsx`, `ui/src/pages/simulator/index.ts`
- Modify: `ui/src/App.tsx`

**Step 1: Create the simulator page**

```tsx
// ui/src/pages/simulator/simulator.tsx
import { useState } from "react";
import { useList, useCustom } from "@refinedev/core";
import { useSearchParams } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { CheckCircle2, XCircle } from "lucide-react";

export function Simulator() {
  const [searchParams] = useSearchParams();
  const [userId, setUserId] = useState(searchParams.get("user") ?? "");
  const [action, setAction] = useState("");
  const [resource, setResource] = useState("*");
  const [result, setResult] = useState<any>(null);
  const [submitted, setSubmitted] = useState(false);

  // Fetch users for dropdown
  const { data: usersData } = useList({
    resource: "users",
    pagination: { pageSize: 100 },
  });

  // Fetch all policies to extract action suggestions
  const { data: policiesData } = useList({
    resource: "iam_policies",
    pagination: { pageSize: 100 },
  });
  const actionSuggestions = Array.from(
    new Set(
      (policiesData?.data ?? []).flatMap((p: any) =>
        (p.document?.statement ?? []).flatMap((s: any) => s.action ?? [])
      )
    )
  ).sort();

  // Simulate call
  const { refetch, isLoading } = useCustom({
    url: "/api/iam/simulate",
    method: "post",
    config: {
      payload: { user_id: userId, action, resource },
    },
    queryOptions: {
      enabled: false,
      onSuccess: (data: any) => setResult(data?.data),
    },
  });

  const handleSimulate = () => {
    setSubmitted(true);
    refetch();
  };

  return (
    <div className="max-w-2xl space-y-6">
      <h1 className="text-2xl font-bold">Policy Simulator</h1>
      <p className="text-muted-foreground">
        Test whether a user is allowed to perform an action.
      </p>

      <div className="space-y-4">
        <div className="space-y-2">
          <Label>User</Label>
          <Select value={userId} onValueChange={setUserId}>
            <SelectTrigger>
              <SelectValue placeholder="Select user..." />
            </SelectTrigger>
            <SelectContent>
              {(usersData?.data ?? []).map((u: any) => (
                <SelectItem key={u.id} value={u.id}>
                  {u.email}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-2">
          <Label>Action</Label>
          <Input
            value={action}
            onChange={(e) => setAction(e.target.value)}
            placeholder="e.g. collections:posts:read"
            list="action-suggestions"
          />
          <datalist id="action-suggestions">
            {actionSuggestions.map((a) => (
              <option key={a} value={a} />
            ))}
          </datalist>
        </div>

        <div className="space-y-2">
          <Label>Resource</Label>
          <Input
            value={resource}
            onChange={(e) => setResource(e.target.value)}
            placeholder="* (default)"
          />
        </div>

        <Button
          onClick={handleSimulate}
          disabled={!userId || !action || isLoading}
        >
          {isLoading ? "Simulating..." : "Simulate"}
        </Button>
      </div>

      {submitted && result && (
        <Card>
          <CardContent className="pt-6 space-y-4">
            <div className="flex items-center gap-2">
              {result.allowed ? (
                <>
                  <CheckCircle2 className="h-6 w-6 text-green-600" />
                  <Badge variant="default" className="text-lg">
                    ALLOWED
                  </Badge>
                </>
              ) : (
                <>
                  <XCircle className="h-6 w-6 text-red-600" />
                  <Badge variant="destructive" className="text-lg">
                    DENIED
                  </Badge>
                </>
              )}
            </div>

            {result.matched_statement && (
              <div className="space-y-1 text-sm">
                <p>
                  <span className="text-muted-foreground">Matched: </span>
                  {result.matched_statement.sid} ({result.matched_statement.effect})
                </p>
                <p>
                  <span className="text-muted-foreground">Policy: </span>
                  {result.matched_statement.policy_name}
                </p>
                {result.attached_via && (
                  <p>
                    <span className="text-muted-foreground">Via: </span>
                    {result.attached_via.type}: {result.attached_via.name}
                  </p>
                )}
              </div>
            )}

            {result.trace && (
              <div>
                <h3 className="mb-2 text-sm font-medium">Evaluation Trace</h3>
                <ol className="list-decimal space-y-1 pl-5 text-sm text-muted-foreground">
                  {result.trace.map((step: string, i: number) => (
                    <li key={i}>{step}</li>
                  ))}
                </ol>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
```

**Step 2: Wire route and commit**

```bash
git add ui/src/pages/simulator/ ui/src/App.tsx
git commit -m "add Policy Simulator page"
```

---

## Task 21: Final wiring and build verification

**Files:**
- Modify: `ui/src/App.tsx` — ensure all route placeholders replaced
- Modify: `.gitignore`

**Step 1: Verify all routes are wired**

Check App.tsx — every `<Route>` should point to a real component, no `(TODO)` placeholders.

**Step 2: Build for production**

```bash
cd ui && npm run build
```

Expected: Output appears in `../pb_public/` — `index.html` + `assets/` folder.

**Step 3: Test embedded serving**

```bash
cd .. && go run . serve
```

Navigate to `http://127.0.0.1:8090` — should serve the dashboard.

**Step 4: Commit**

```bash
git add ui/ .gitignore
git commit -m "finalize dashboard wiring and build config"
```

---

## Summary

| Task | What | Key Files |
|------|------|-----------|
| 1 | Scaffold Vite project | `ui/` |
| 2 | Tailwind + shadcn/ui | `ui/vite.config.ts`, `ui/components.json` |
| 3 | Install Refine + PB provider | `ui/package.json` |
| 4 | PB client + providers | `ui/src/providers/pocketbase.ts` |
| 5 | Sidebar layout | `ui/src/components/layout/` |
| 6 | App.tsx + routing + login | `ui/src/App.tsx`, `ui/src/pages/login.tsx` |
| 7 | .gitignore | `.gitignore` |
| 8 | EntityPicker | `ui/src/components/entity-picker/` |
| 9 | EntityChipList | `ui/src/components/entity-chip-list/` |
| 10 | JSON editor (CodeMirror) | `ui/src/components/json-editor/` |
| 11 | PolicyFormBuilder | `ui/src/components/policy-form-builder/` |
| 12 | EffectivePermissionsTable | `ui/src/components/effective-permissions-table/` |
| 13 | Policies list | `ui/src/pages/policies/list.tsx` |
| 14 | Policies create/edit | `ui/src/pages/policies/create.tsx`, `edit.tsx` |
| 15 | Roles pages | `ui/src/pages/roles/` |
| 16 | Groups pages | `ui/src/pages/groups/` |
| 17 | Users pages | `ui/src/pages/users/` |
| 18 | Managed Collections page | `ui/src/pages/managed-collections/` |
| 19 | Backend: /api/iam/simulate | `iam/routes.go`, `iam/engine.go` |
| 20 | Policy Simulator page | `ui/src/pages/simulator/` |
| 21 | Final wiring + build | All |
