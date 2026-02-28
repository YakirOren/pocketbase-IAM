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
