import { Refine, Authenticated } from "@refinedev/core";
import routerProvider, {
  NavigateToResource,
  CatchAllNavigate,
  UnsavedChangesNotifier,
  DocumentTitleHandler,
} from "@refinedev/react-router";
import { BrowserRouter, Routes, Route, Outlet } from "react-router";

import { pbDataProvider, pbAuthProvider, pbLiveProvider } from "@/providers/pocketbase";
import { accessControlProvider } from "@/providers/access-control";
import { ThemeProvider } from "@/components/refine-ui/theme/theme-provider";
import { Layout } from "@/components/layout";
import { LoginPage } from "@/pages/login";
import { PolicyList, PolicyCreate, PolicyEdit } from "@/pages/policies";
import { RoleList, RoleCreate, RoleEdit, RoleShow } from "@/pages/roles";
import { GroupList, GroupCreate, GroupEdit, GroupShow } from "@/pages/groups";
import { UserList, UserShow } from "@/pages/users";
import { ManagedCollectionList } from "@/pages/managed-collections";
import { Simulator } from "@/pages/simulator";

function App() {
  return (
    <ThemeProvider>
    <BrowserRouter basename="/_/iam">
      <Refine
        routerProvider={routerProvider}
        dataProvider={pbDataProvider}
        authProvider={pbAuthProvider}
        liveProvider={pbLiveProvider}
        accessControlProvider={accessControlProvider}
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
              <Route index element={<PolicyList />} />
              <Route path="create" element={<PolicyCreate />} />
              <Route path="edit/:id" element={<PolicyEdit />} />
              <Route path="show/:id" element={<PolicyEdit />} />
            </Route>

            <Route path="/roles">
              <Route index element={<RoleList />} />
              <Route path="create" element={<RoleCreate />} />
              <Route path="edit/:id" element={<RoleEdit />} />
              <Route path="show/:id" element={<RoleShow />} />
            </Route>

            <Route path="/groups">
              <Route index element={<GroupList />} />
              <Route path="create" element={<GroupCreate />} />
              <Route path="edit/:id" element={<GroupEdit />} />
              <Route path="show/:id" element={<GroupShow />} />
            </Route>

            <Route path="/users">
              <Route index element={<UserList />} />
              <Route path="show/:id" element={<UserShow />} />
            </Route>

            <Route path="/managed-collections">
              <Route index element={<ManagedCollectionList />} />
            </Route>

            <Route path="/simulator" element={<Simulator />} />

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
    </ThemeProvider>
  );
}

export default App;
