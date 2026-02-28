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

  const onSubmit = (data: Record<string, unknown>) => {
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
  const attachedUsers = (userPolicies?.data ?? []).map((r: Record<string, unknown>) => ({
    id: r.id as string,
    entityId: r.user as string,
    label: ((r as Record<string, unknown> & { expand?: { user?: { email?: string } } }).expand?.user?.email ?? r.user) as string,
  }));

  // Roles attached via iam_role_policies
  const { data: rolePolicies } = useList({
    resource: "iam_role_policies",
    filters: [{ field: "policy", operator: "eq", value: id }],
    meta: { expand: ["role"] },
  });
  const attachedRoles = (rolePolicies?.data ?? []).map((r: Record<string, unknown>) => ({
    id: r.id as string,
    entityId: r.role as string,
    label: ((r as Record<string, unknown> & { expand?: { role?: { name?: string } } }).expand?.role?.name ?? r.role) as string,
  }));

  // Groups attached via iam_group_policies
  const { data: groupPolicies } = useList({
    resource: "iam_group_policies",
    filters: [{ field: "policy", operator: "eq", value: id }],
    meta: { expand: ["group"] },
  });
  const attachedGroups = (groupPolicies?.data ?? []).map((r: Record<string, unknown>) => ({
    id: r.id as string,
    entityId: r.group as string,
    label: ((r as Record<string, unknown> & { expand?: { group?: { name?: string } } }).expand?.group?.name ?? r.group) as string,
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
