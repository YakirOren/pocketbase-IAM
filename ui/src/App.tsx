import { Refine, Authenticated, AuthPage } from "@refinedev/core";
import routerProvider, {
  NavigateToResource,
  CatchAllNavigate,
  UnsavedChangesNotifier,
  DocumentTitleHandler,
} from "@refinedev/react-router";
import { BrowserRouter, Routes, Route, Outlet } from "react-router";

import { pbDataProvider, pbAuthProvider, pbLiveProvider } from "@/providers/pocketbase";
import { Layout } from "@/components/layout";

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
            meta: { canDelete: true, label: "Policies" },
          },
          {
            name: "iam_roles",
            list: "/roles",
            create: "/roles/create",
            edit: "/roles/edit/:id",
            show: "/roles/show/:id",
            meta: { canDelete: true, label: "Roles" },
          },
          {
            name: "iam_groups",
            list: "/groups",
            create: "/groups/create",
            edit: "/groups/edit/:id",
            show: "/groups/show/:id",
            meta: { canDelete: true, label: "Groups" },
          },
          {
            name: "users",
            list: "/users",
            show: "/users/show/:id",
            meta: { label: "Users" },
          },
          {
            name: "iam_managed_collections",
            list: "/managed-collections",
            meta: { label: "Managed Collections" },
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
            <Route path="/login" element={<AuthPage type="login" />} />
          </Route>
        </Routes>

        <UnsavedChangesNotifier />
        <DocumentTitleHandler />
      </Refine>
    </BrowserRouter>
  );
}

export default App;
