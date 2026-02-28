import { AuthPage } from "@refinedev/core";
import { Shield } from "lucide-react";

export function LoginPage() {
  return (
    <div className="flex h-screen items-center justify-center bg-muted">
      <div className="w-full max-w-sm space-y-6 rounded-lg border bg-card p-8 shadow-sm [&_form]:space-y-4 [&_label]:text-sm [&_label]:font-medium [&_input]:flex [&_input]:h-9 [&_input]:w-full [&_input]:rounded-md [&_input]:border [&_input]:border-black [&_input]:rounded [&_input]:bg-transparent [&_input]:px-3 [&_input]:py-1 [&_input]:text-sm [&_input]:transition-colors [&_input]:placeholder:text-muted-foreground [&_input]:focus-visible:outline-none [&_input]:focus-visible:ring-1 [&_input]:focus-visible:ring-black [&_button[type=submit]]:inline-flex [&_button[type=submit]]:w-full [&_button[type=submit]]:items-center [&_button[type=submit]]:justify-center [&_button[type=submit]]:gap-2 [&_button[type=submit]]:whitespace-nowrap [&_button[type=submit]]:rounded-md [&_button[type=submit]]:bg-primary [&_button[type=submit]]:px-4 [&_button[type=submit]]:py-2 [&_button[type=submit]]:text-sm [&_button[type=submit]]:font-medium [&_button[type=submit]]:text-primary-foreground [&_button[type=submit]]:shadow [&_button[type=submit]]:transition-colors [&_button[type=submit]]:hover:bg-primary/90 [&_button[type=submit]]:cursor-pointer">
        <div className="flex flex-col items-center gap-2">
          <Shield className="h-8 w-8" />
          <h1 className="text-xl font-semibold">IAM Dashboard</h1>
        </div>
        <AuthPage
          type="login"
          registerLink={false}
          forgotPasswordLink={false}
          rememberMe={false}
        />
      </div>
    </div>
  );
}
