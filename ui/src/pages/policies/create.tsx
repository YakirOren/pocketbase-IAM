import { useState } from "react";
import { useNavigate } from "react-router";
import { useForm } from "@refinedev/react-hook-form";
import { ArrowLeft } from "lucide-react";
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

  const onSubmit = (data: Record<string, unknown>) => {
    onFinish({ ...data, document: policyDoc });
  };

  const navigate = useNavigate();

  return (
    <div className="max-w-3xl">
      <Button variant="ghost" size="sm" className="mb-4" onClick={() => navigate(-1)}>
        <ArrowLeft className="mr-1 h-4 w-4" /> Back
      </Button>
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
