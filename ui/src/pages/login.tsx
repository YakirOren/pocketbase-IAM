import { AuthPage } from "@refinedev/core";
import { Shield } from "lucide-react";

export function LoginPage() {
  return (
    <div className="flex h-screen items-center justify-center bg-muted">
      <div className="w-full max-w-sm space-y-6 rounded-lg border bg-card p-8 shadow-sm">
        <div className="flex flex-col items-center gap-2">
          <Shield className="h-8 w-8" />
          <h1 className="text-xl font-semibold">IAM Dashboard</h1>
        </div>
        <AuthPage type="login" />
      </div>
    </div>
  );
}
