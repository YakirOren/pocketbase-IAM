import { useShow, useList, useCreate, useDelete } from "@refinedev/core";
import { useParams, Link, useNavigate } from "react-router";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Pencil } from "lucide-react";
import { EntityChipList } from "@/components/entity-chip-list";
import { EffectivePermissionsTable } from "@/components/effective-permissions-table";
import { computeEffectivePermissions } from "@/lib/permissions";
import type { PolicyDocument } from "@/types/policy";

interface AttachedPolicy {
  id: string;
  entityId: string;
  label: string;
  document: PolicyDocument | undefined;
}

interface AttachedEntity {
  id: string;
  entityId: string;
  label: string;
}

export function RoleShow() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { query } = useShow({ resource: "iam_roles", id });
  const record = query?.data?.data;

  const { mutate: createJoin } = useCreate();
  const { mutate: deleteJoin } = useDelete();

  // Attached policies
  const { result: rolePoliciesResult } = useList({
    resource: "iam_role_policies",
    filters: [{ field: "role", operator: "eq", value: id }],
    meta: { expand: ["policy"] },
  });
  const attachedPolicies: AttachedPolicy[] = (rolePoliciesResult.data ?? []).map((r: Record<string, unknown>) => {
    const expand = r.expand as { policy?: { name?: string; document?: PolicyDocument } } | undefined;
    return {
      id: r.id as string,
      entityId: r.policy as string,
      label: (expand?.policy?.name ?? r.policy) as string,
      document: expand?.policy?.document,
    };
  });

  // Assigned users
  const { result: userRolesResult } = useList({
    resource: "iam_user_roles",
    filters: [{ field: "role", operator: "eq", value: id }],
    meta: { expand: ["user"] },
  });
  const assignedUsers: AttachedEntity[] = (userRolesResult.data ?? []).map((r: Record<string, unknown>) => {
    const expand = r.expand as { user?: { email?: string } } | undefined;
    return {
      id: r.id as string,
      entityId: r.user as string,
      label: (expand?.user?.email ?? r.user) as string,
    };
  });

  // Effective permissions
  const effectivePerms = computeEffectivePermissions(
    attachedPolicies
      .filter((p) => p.document)
      .map((p) => ({
        policyName: p.label,
        document: p.document!,
        source: `Role: ${(record?.name as string) ?? ""}`,
      }))
  );

  if (query?.isLoading) return <p>Loading...</p>;

  return (
    <div className="max-w-3xl space-y-6">
      <Button variant="ghost" size="sm" onClick={() => navigate(-1)}>
        <ArrowLeft className="mr-1 h-4 w-4" /> Back
      </Button>
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{record?.name as string}</h1>
        <Button variant="outline" asChild>
          <Link to={`/roles/edit/${id}`}>
            <Pencil className="mr-1 h-4 w-4" /> Edit
          </Link>
        </Button>
      </div>

      {record?.description && (
        <p className="text-muted-foreground">{record.description as string}</p>
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
