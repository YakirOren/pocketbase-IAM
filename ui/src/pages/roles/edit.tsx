import { useForm } from "@refinedev/react-hook-form";
import { useParams, useNavigate } from "react-router";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export function RoleEdit() {
  const { id } = useParams();
  const navigate = useNavigate();
  const {
    refineCore: { onFinish, formLoading, query },
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    refineCoreProps: {
      resource: "iam_roles",
      action: "edit",
      id,
    },
  });

  if (query?.isLoading) return <p>Loading...</p>;

  return (
    <div className="max-w-xl">
      <Button variant="ghost" size="sm" className="mb-4" onClick={() => navigate(-1)}>
        <ArrowLeft className="mr-1 h-4 w-4" /> Back
      </Button>
      <h1 className="mb-6 text-2xl font-bold">Edit Role</h1>
      <form onSubmit={handleSubmit(onFinish)} className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="name">Name</Label>
          <Input id="name" {...register("name", { required: "Name is required" })} />
          {errors.name && <p className="text-sm text-destructive">{errors.name.message as string}</p>}
        </div>
        <div className="space-y-2">
          <Label htmlFor="description">Description</Label>
          <Input id="description" {...register("description")} />
        </div>
        <Button type="submit" disabled={formLoading}>
          {formLoading ? "Saving..." : "Save Role"}
        </Button>
      </form>
    </div>
  );
}
