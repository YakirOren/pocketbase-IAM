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
  const roles = (userRoles?.data ?? []).map((r: Record<string, unknown>) => {
    const expand = r.expand as { role?: { name?: string } } | undefined;
    return {
      id: r.id as string,
      entityId: r.role as string,
      label: (expand?.role?.name ?? r.role) as string,
    };
  });

  // User's groups
  const { data: groupUsers } = useList({
    resource: "iam_group_users",
    filters: [{ field: "user", operator: "eq", value: id }],
    meta: { expand: ["group"] },
  });
  const groups = (groupUsers?.data ?? []).map((r: Record<string, unknown>) => {
    const expand = r.expand as { group?: { name?: string } } | undefined;
    return {
      id: r.id as string,
      entityId: r.group as string,
      label: (expand?.group?.name ?? r.group) as string,
    };
  });

  // User's direct policies
  const { data: userPolicies } = useList({
    resource: "iam_user_policies",
    filters: [{ field: "user", operator: "eq", value: id }],
    meta: { expand: ["policy"] },
  });
  const directPolicies = (userPolicies?.data ?? []).map((r: Record<string, unknown>) => {
    const expand = r.expand as { policy?: { name?: string; document?: PolicyDocument } } | undefined;
    return {
      id: r.id as string,
      entityId: r.policy as string,
      label: (expand?.policy?.name ?? r.policy) as string,
      document: expand?.policy?.document,
    };
  });

  // Role policies (for effective permissions)
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

  // Group policies (for effective permissions)
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
      .filter((r: Record<string, unknown>) => {
        const expand = r.expand as { policy?: { document?: unknown } } | undefined;
        return expand?.policy?.document;
      })
      .map((r: Record<string, unknown>) => {
        const expand = r.expand as { policy?: { name?: string; document?: PolicyDocument }; role?: { name?: string } } | undefined;
        return {
          policyName: expand?.policy?.name ?? (r.policy as string),
          document: expand!.policy!.document as PolicyDocument,
          source: `Role: ${expand?.role?.name ?? r.role}`,
        };
      }),
    ...(groupPoliciesData?.data ?? [])
      .filter((r: Record<string, unknown>) => {
        const expand = r.expand as { policy?: { document?: unknown } } | undefined;
        return expand?.policy?.document;
      })
      .map((r: Record<string, unknown>) => {
        const expand = r.expand as { policy?: { name?: string; document?: PolicyDocument }; group?: { name?: string } } | undefined;
        return {
          policyName: expand?.policy?.name ?? (r.policy as string),
          document: expand!.policy!.document as PolicyDocument,
          source: `Group: ${expand?.group?.name ?? r.group}`,
        };
      }),
  ];

  const effectivePerms = computeEffectivePermissions(allPolicySources);

  if (query?.isLoading) return <p>Loading...</p>;

  return (
    <div className="max-w-3xl space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{user?.email as string}</h1>
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
          createJoin({ resource: "iam_user_roles", values: { user: id, role: roleId } })
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
          createJoin({ resource: "iam_group_users", values: { user: id, group: groupId } })
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
          createJoin({ resource: "iam_user_policies", values: { user: id, policy: policyId } })
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
